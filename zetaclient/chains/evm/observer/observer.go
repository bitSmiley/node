// Package observer implements the EVM chain observer
package observer

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.non-eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/zetacore/pkg/bg"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/db"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

var _ interfaces.ChainObserver = (*Observer)(nil)

// Observer is the observer for evm chains
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	priorityFeeConfig

	// evmClient is the EVM client for the observed chain
	evmClient interfaces.EVMRPCClient

	// evmJSONRPC is the EVM JSON RPC client for the observed chain
	evmJSONRPC interfaces.EVMJSONRPCClient

	// outboundPendingTransactions is the map to index pending transactions by hash
	outboundPendingTransactions map[string]*ethtypes.Transaction

	// outboundConfirmedReceipts is the map to index confirmed receipts by hash
	outboundConfirmedReceipts map[string]*ethtypes.Receipt

	// outboundConfirmedTransactions is the map to index confirmed transactions by hash
	outboundConfirmedTransactions map[string]*ethtypes.Transaction
}

// priorityFeeConfig is the configuration for priority fee
type priorityFeeConfig struct {
	checked   bool
	supported bool
}

// NewObserver returns a new EVM chain observer
func NewObserver(
	ctx context.Context,
	evmCfg config.EVMConfig,
	evmClient interfaces.EVMRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	database *db.DB,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	// create base observer
	baseObserver, err := base.NewObserver(
		evmCfg.Chain,
		chainParams,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		base.DefaultHeaderCacheSize,
		ts,
		database,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	// create evm observer
	ob := &Observer{
		Observer:                      *baseObserver,
		evmClient:                     evmClient,
		evmJSONRPC:                    ethrpc.NewEthRPC(evmCfg.Endpoint),
		outboundPendingTransactions:   make(map[string]*ethtypes.Transaction),
		outboundConfirmedReceipts:     make(map[string]*ethtypes.Receipt),
		outboundConfirmedTransactions: make(map[string]*ethtypes.Transaction),
		priorityFeeConfig:             priorityFeeConfig{},
	}

	// load last block scanned
	if err = ob.LoadLastBlockScanned(ctx); err != nil {
		return nil, errors.Wrap(err, "unable to load last block scanned")
	}

	return ob, nil
}

// WithEvmClient attaches a new evm client to the observer
func (ob *Observer) WithEvmClient(client interfaces.EVMRPCClient) {
	ob.evmClient = client
}

// WithEvmJSONRPC attaches a new evm json rpc client to the observer
func (ob *Observer) WithEvmJSONRPC(client interfaces.EVMJSONRPCClient) {
	ob.evmJSONRPC = client
}

// SetChainParams sets the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.WithChainParams(params)
}

// GetChainParams returns the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) GetChainParams() observertypes.ChainParams {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.ChainParams()
}

// GetConnectorContract returns the non-Eth connector address and binder
func (ob *Observer) GetConnectorContract() (ethcommon.Address, *zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := zetaconnector.NewZetaConnectorNonEth(addr, ob.evmClient)
	return addr, contract, err
}

// GetConnectorContractEth returns the Eth connector address and binder
func (ob *Observer) GetConnectorContractEth() (ethcommon.Address, *zetaconnectoreth.ZetaConnectorEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContractEth(addr, ob.evmClient)
	return addr, contract, err
}

// GetERC20CustodyContract returns ERC20Custody contract address and binder
func (ob *Observer) GetERC20CustodyContract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().Erc20CustodyContractAddress)
	contract, err := erc20custody.NewERC20Custody(addr, ob.evmClient)
	return addr, contract, err
}

// FetchConnectorContractEth returns the Eth connector address and binder
// TODO(revamp): move this to a contract package
func FetchConnectorContractEth(
	addr ethcommon.Address,
	client interfaces.EVMRPCClient,
) (*zetaconnectoreth.ZetaConnectorEth, error) {
	return zetaconnectoreth.NewZetaConnectorEth(addr, client)
}

// FetchZetaTokenContract returns the non-Eth ZETA token binder
// TODO(revamp): move this to a contract package
func FetchZetaTokenContract(
	addr ethcommon.Address,
	client interfaces.EVMRPCClient,
) (*zeta.ZetaNonEth, error) {
	return zeta.NewZetaNonEth(addr, client)
}

// Start all observation routines for the evm chain
func (ob *Observer) Start(ctx context.Context) {
	if noop := ob.Observer.Start(); noop {
		ob.Logger().Chain.Info().Msgf("observer is already started for chain %d", ob.Chain().ChainId)
		return
	}

	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))
	bg.Work(ctx, ob.WatchOutbound, bg.WithName("WatchOutbound"), bg.WithLogger(ob.Logger().Outbound))
	bg.Work(ctx, ob.WatchGasPrice, bg.WithName("WatchGasPrice"), bg.WithLogger(ob.Logger().GasPrice))
	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))
	bg.Work(ctx, ob.WatchRPCStatus, bg.WithName("WatchRPCStatus"), bg.WithLogger(ob.Logger().Chain))
}

// WatchRPCStatus watches the RPC status of the evm chain
// TODO(revamp): move ticker to ticker file
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchRPCStatus(ctx context.Context) error {
	ob.Logger().Chain.Info().Msgf("Starting RPC status check for chain %d", ob.Chain().ChainId)
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}
			bn, err := ob.evmClient.BlockNumber(ctx)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			gasPrice, err := ob.evmClient.SuggestGasPrice(ctx)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			header, err := ob.evmClient.HeaderByNumber(ctx, new(big.Int).SetUint64(bn))
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			// #nosec G115 always in range
			blockTime := time.Unix(int64(header.Time), 0).UTC()
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 100 {
				ob.Logger().Chain.Warn().
					Msgf("RPC Status Check warning: RPC stale or chain stuck (check explorer)? Latest block %d timestamp is %.0fs ago", bn, elapsedSeconds)
				continue
			}
			ob.Logger().Chain.Info().
				Msgf("[OK] RPC status: latest block num %d, timestamp %s ( %.0fs ago), suggested gas price %d", header.Number, blockTime.String(), elapsedSeconds, gasPrice.Uint64())
		case <-ob.StopChannel():
			return nil
		}
	}
}

// SetPendingTx sets the pending transaction in memory
func (ob *Observer) SetPendingTx(nonce uint64, transaction *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.outboundPendingTransactions[ob.OutboundID(nonce)] = transaction
}

// GetPendingTx gets the pending transaction from memory
func (ob *Observer) GetPendingTx(nonce uint64) *ethtypes.Transaction {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.outboundPendingTransactions[ob.OutboundID(nonce)]
}

// SetTxNReceipt sets the receipt and transaction in memory
func (ob *Observer) SetTxNReceipt(nonce uint64, receipt *ethtypes.Receipt, transaction *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	delete(ob.outboundPendingTransactions, ob.OutboundID(nonce)) // remove pending transaction, if any
	ob.outboundConfirmedReceipts[ob.OutboundID(nonce)] = receipt
	ob.outboundConfirmedTransactions[ob.OutboundID(nonce)] = transaction
}

// GetTxNReceipt gets the receipt and transaction from memory
func (ob *Observer) GetTxNReceipt(nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	receipt := ob.outboundConfirmedReceipts[ob.OutboundID(nonce)]
	transaction := ob.outboundConfirmedTransactions[ob.OutboundID(nonce)]
	return receipt, transaction
}

// IsTxConfirmed returns true if there is a confirmed tx for 'nonce'
func (ob *Observer) IsTxConfirmed(nonce uint64) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.outboundConfirmedReceipts[ob.OutboundID(nonce)] != nil &&
		ob.outboundConfirmedTransactions[ob.OutboundID(nonce)] != nil
}

// CheckTxInclusion returns nil only if tx is included at the position indicated by the receipt ([block, index])
func (ob *Observer) CheckTxInclusion(tx *ethtypes.Transaction, receipt *ethtypes.Receipt) error {
	block, err := ob.GetBlockByNumberCached(receipt.BlockNumber.Uint64())
	if err != nil {
		return errors.Wrapf(err, "GetBlockByNumberCached error for block %d txHash %s nonce %d",
			receipt.BlockNumber.Uint64(), tx.Hash(), tx.Nonce())
	}

	// #nosec G115 non negative value
	if receipt.TransactionIndex >= uint(len(block.Transactions)) {
		return fmt.Errorf("transaction index %d out of range [0, %d), txHash %s nonce %d block %d",
			receipt.TransactionIndex, len(block.Transactions), tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	txAtIndex := block.Transactions[receipt.TransactionIndex]
	if !strings.EqualFold(txAtIndex.Hash, tx.Hash().Hex()) {
		ob.RemoveCachedBlock(receipt.BlockNumber.Uint64()) // clean stale block from cache
		return fmt.Errorf("transaction at index %d has different hash %s, txHash %s nonce %d block %d",
			receipt.TransactionIndex, txAtIndex.Hash, tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	return nil
}

// TransactionByHash query transaction by hash via JSON-RPC
// TODO(revamp): update this method as a pure RPC method that takes two parameters (jsonRPC, and txHash) and move to upper package to file rpc.go
func (ob *Observer) TransactionByHash(txHash string) (*ethrpc.Transaction, bool, error) {
	tx, err := ob.evmJSONRPC.EthGetTransactionByHash(txHash)
	if err != nil {
		return nil, false, err
	}
	err = evm.ValidateEvmTransaction(tx)
	if err != nil {
		return nil, false, err
	}
	return tx, tx.BlockNumber == nil, nil
}

// GetBlockHeaderCached get block header by number from cache
func (ob *Observer) GetBlockHeaderCached(ctx context.Context, blockNumber uint64) (*ethtypes.Header, error) {
	if result, ok := ob.HeaderCache().Get(blockNumber); ok {
		if header, ok := result.(*ethtypes.Header); ok {
			return header, nil
		}
		return nil, errors.New("cached value is not of type *ethtypes.Header")
	}
	header, err := ob.evmClient.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.HeaderCache().Add(blockNumber, header)
	return header, nil
}

// GetBlockByNumberCached get block by number from cache
// returns block, ethrpc.Block, isFallback, isSkip, error
func (ob *Observer) GetBlockByNumberCached(blockNumber uint64) (*ethrpc.Block, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*ethrpc.Block); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *ethrpc.Block")
	}
	if blockNumber > math.MaxInt32 {
		return nil, fmt.Errorf("block number %d is too large", blockNumber)
	}
	// #nosec G115 always in range, checked above
	block, err := ob.BlockByNumber(int(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.BlockCache().Add(blockNumber, block)
	return block, nil
}

// RemoveCachedBlock remove block from cache
func (ob *Observer) RemoveCachedBlock(blockNumber uint64) {
	ob.BlockCache().Remove(blockNumber)
}

// BlockByNumber query block by number via JSON-RPC
func (ob *Observer) BlockByNumber(blockNumber int) (*ethrpc.Block, error) {
	block, err := ob.evmJSONRPC.EthGetBlockByNumber(blockNumber, true)
	if err != nil {
		return nil, err
	}
	for i := range block.Transactions {
		err := evm.ValidateEvmTransaction(&block.Transactions[i])
		if err != nil {
			return nil, err
		}
	}
	return block, nil
}

// LoadLastBlockScanned loads the last scanned block from the database
// TODO(revamp): move to a db file
func (ob *Observer) LoadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned(ob.Logger().Chain)
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.evmClient.BlockNumber(ctx)
		if err != nil {
			return errors.Wrapf(err, "error BlockNumber for chain %d", ob.Chain().ChainId)
		}
		ob.WithLastBlockScanned(blockNumber)
	}
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from block %d", ob.Chain().ChainId, ob.LastBlockScanned())

	return nil
}
