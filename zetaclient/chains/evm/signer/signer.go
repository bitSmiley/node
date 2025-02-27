// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/constant"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

const (
	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1000 * time.Millisecond

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5
)

var (
	_ interfaces.ChainSigner = (*Signer)(nil)

	// zeroValue is for outbounds that carry no ETH (gas token) value
	zeroValue = big.NewInt(0)
)

// Signer deals with the signing EVM transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the EVM RPC client to interact with the EVM chain
	client interfaces.EVMRPCClient

	// ethSigner encapsulates EVM transaction signature handling
	ethSigner ethtypes.Signer

	// zetaConnectorABI is the ABI of the ZetaConnector contract
	zetaConnectorABI abi.ABI

	// erc20CustodyABI is the ABI of the ERC20Custody contract
	erc20CustodyABI abi.ABI

	// zetaConnectorAddress is the address of the ZetaConnector contract
	zetaConnectorAddress ethcommon.Address

	// er20CustodyAddress is the address of the ERC20Custody contract
	er20CustodyAddress ethcommon.Address
}

// NewSigner creates a new EVM signer
func NewSigner(
	ctx context.Context,
	chain chains.Chain,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger base.Logger,
	endpoint string,
	zetaConnectorABI string,
	erc20CustodyABI string,
	zetaConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, ts, logger)

	// create EVM client
	client, ethSigner, err := getEVMRPC(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create EVM client")
	}

	// prepare ABIs
	connectorABI, err := abi.JSON(strings.NewReader(zetaConnectorABI))
	if err != nil {
		return nil, errors.Wrap(err, "unable to build ZetaConnector ABI")
	}

	custodyABI, err := abi.JSON(strings.NewReader(erc20CustodyABI))
	if err != nil {
		return nil, errors.Wrap(err, "unable to build ERC20Custody ABI")
	}

	return &Signer{
		Signer:               baseSigner,
		client:               client,
		ethSigner:            ethSigner,
		zetaConnectorABI:     connectorABI,
		erc20CustodyABI:      custodyABI,
		zetaConnectorAddress: zetaConnectorAddress,
		er20CustodyAddress:   erc20CustodyAddress,
	}, nil
}

// SetZetaConnectorAddress sets the zeta connector address
func (signer *Signer) SetZetaConnectorAddress(addr ethcommon.Address) {
	signer.Lock()
	defer signer.Unlock()
	signer.zetaConnectorAddress = addr
}

// SetERC20CustodyAddress sets the erc20 custody address
func (signer *Signer) SetERC20CustodyAddress(addr ethcommon.Address) {
	signer.Lock()
	defer signer.Unlock()
	signer.er20CustodyAddress = addr
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(_ string) {
	// Note: do nothing for now
	// gateway address will be needed in the future contract architecture
}

// GetZetaConnectorAddress returns the zeta connector address
func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	signer.Lock()
	defer signer.Unlock()
	return signer.zetaConnectorAddress
}

// GetERC20CustodyAddress returns the erc20 custody address
func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	signer.Lock()
	defer signer.Unlock()
	return signer.er20CustodyAddress
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	// Note: return empty string for now
	// gateway address will be needed in the future contract architecture
	return ""
}

// Sign given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(
	ctx context.Context,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
	height uint64,
) (*ethtypes.Transaction, []byte, []byte, error) {
	signer.Logger().Std.Debug().
		Str("tss_pub_key", signer.TSS().EVMAddress().String()).
		Msg("Signing evm transaction")

	chainID := big.NewInt(signer.Chain().ChainId)
	tx, err := newTx(chainID, data, to, amount, gas, nonce)
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes := signer.ethSigner.Hash(tx).Bytes()

	sig, err := signer.TSS().Sign(ctx, hashBytes, height, nonce, signer.Chain().ChainId, "")
	if err != nil {
		return nil, nil, nil, err
	}

	log.Debug().Msgf("Sign: Signature: %s", hex.EncodeToString(sig[:]))
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msgf("SigToPub error")
	}

	addr := crypto.PubkeyToAddress(*pubk)
	signer.Logger().Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.ethSigner, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}

	return signedTX, sig[:], hashBytes[:], nil
}

func newTx(
	chainID *big.Int,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
) (*ethtypes.Transaction, error) {
	if err := gas.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid gas parameters")
	}

	if gas.isLegacy() {
		return ethtypes.NewTx(&ethtypes.LegacyTx{
			To:       &to,
			Value:    amount,
			Data:     data,
			GasPrice: gas.Price,
			Gas:      gas.Limit,
			Nonce:    nonce,
		}), nil
	}

	return ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		To:        &to,
		Value:     amount,
		Data:      data,
		GasFeeCap: gas.Price,
		GasTipCap: gas.PriorityFee,
		Gas:       gas.Limit,
		Nonce:     nonce,
	}), nil
}

func (signer *Signer) broadcast(ctx context.Context, tx *ethtypes.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	return signer.client.SendTransaction(ctx, tx)
}

// SignOutbound
// function onReceive(
//
//	bytes calldata originSenderAddress,
//	uint256 originChainId,
//	address destinationAddress,
//	uint zetaAmount,
//	bytes calldata message,
//	bytes32 internalSendHash
//
// ) external virtual {}
func (signer *Signer) SignOutbound(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	data, err = signer.zetaConnectorABI.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	if err != nil {
		return nil, fmt.Errorf("onReceive pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign onReceive error: %w", err)
	}

	return tx, nil
}

// SignRevertTx
// function onRevert(
// address originSenderAddress,
// uint256 originChainId,
// bytes calldata destinationAddress,
// uint256 destinationChainId,
// uint256 zetaAmount,
// bytes calldata message,
// bytes32 internalSendHash
// ) external override whenNotPaused onlyTssAddress
func (signer *Signer) SignRevertTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := signer.zetaConnectorABI.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.cctxIndex,
	)
	if err != nil {
		return nil, fmt.Errorf("onRevert pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign onRevert error: %w", err)
	}

	return tx, nil
}

// SignCancelTx signs a transaction from TSS address to itself with a zero amount in order to increment the nonce
func (signer *Signer) SignCancelTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	txData.gas.Limit = evm.EthTransferGasLimit
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		signer.TSS().EVMAddress(),
		zeroValue, // zero out the amount to cancel the tx
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignCancelTx error: %w", err)
	}

	return tx, nil
}

// SignWithdrawTx signs a withdrawal transaction sent from the TSS address to the destination
func (signer *Signer) SignWithdrawTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	txData.gas.Limit = evm.EthTransferGasLimit
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignWithdrawTx error: %w", err)
	}

	return tx, nil
}

// SignCommandTx signs a transaction based on the given command includes:
//
//	cmd_whitelist_erc20
//	cmd_migrate_tss_funds
func (signer *Signer) SignCommandTx(
	ctx context.Context,
	txData *OutboundData,
	cmd string,
	params string,
) (*ethtypes.Transaction, error) {
	switch cmd {
	case constant.CmdWhitelistERC20:
		return signer.SignWhitelistERC20Cmd(ctx, txData, params)
	case constant.CmdMigrateTssFunds:
		return signer.SignMigrateTssFundsCmd(ctx, txData)
	}
	return nil, fmt.Errorf("SignCommandTx: unknown command %s", cmd)
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
// TODO(revamp): simplify function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundProc *outboundprocessor.Processor,
	outboundID string,
	chainObserver interfaces.ChainObserver,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	var (
		params = cctx.GetCurrentOutboundParam()
		myID   = zetacoreClient.GetKeys().GetOperatorAddress()
		logger = signer.Logger().Std.With().
			Str("method", "TryProcessOutbound").
			Int64("chain", signer.Chain().ChainId).
			Uint64("nonce", params.TssNonce).
			Str("cctx.index", cctx.Index).
			Str("cctx.receiver", params.Receiver).
			Str("cctx.amount", params.Amount.String()).
			Logger()
	)

	logger.Info().Msgf("TryProcessOutbound")

	app, err := zctx.FromContext(ctx)
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msg("error getting app context")
		return
	}

	// end outbound process on panic
	defer func() {
		if r := recover(); r != nil {
			logger.Error().Interface("panic", r).Msg("panic in TryProcessOutbound")
		}

		outboundProc.EndTryProcess(outboundID)
	}()

	evmObserver, ok := chainObserver.(*observer.Observer)
	if !ok {
		logger.Error().Msg("chain observer is not an EVM observer")
		return
	}

	// Setup Transaction input
	txData, skipTx, err := NewOutboundData(ctx, cctx, evmObserver, height, logger)
	if err != nil {
		logger.Err(err).Msg("error setting up transaction input fields")
		return
	}

	if skipTx {
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().Msgf("unable to TryProcessOutbound when toChain is zetaChain (%d)", toChain.ID())
		return
	}

	// https://github.com/zeta-chain/node/issues/2050
	var tx *ethtypes.Transaction
	// compliance check goes first
	if compliance.IsCctxRestricted(cctx) {
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			txData.to.Hex(),
			cctx.GetCurrentOutboundParam().CoinType.String(),
		)

		tx, err = signer.SignCancelTx(ctx, txData) // cancel the tx
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.InboundParams.CoinType == coin.CoinType_Cmd { // admin command
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		if to == (ethcommon.Address{}) {
			logger.Error().Msgf("invalid receiver %s", cctx.GetCurrentOutboundParam().Receiver)
			return
		}
		msg := strings.Split(cctx.RelayedMessage, ":")
		if len(msg) != 2 {
			logger.Error().Msgf("invalid message %s", msg)
			return
		}
		// cmd field is used to determine whether to execute ERC20 whitelist or migrate TSS funds given that the coin type
		// from the cctx is coin.CoinType_Cmd
		cmd := msg[0]
		// params field is used to pass input parameters for command requests, currently it is used to pass the ERC20
		// contract address when a whitelist command is requested
		params := msg[1]
		tx, err = signer.SignCommandTx(ctx, txData, cmd, params)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if IsSenderZetaChain(cctx, zetacoreClient) {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignWithdrawTx(ctx, txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf(
				"SignERC20WithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignERC20WithdrawTx(ctx, txData)
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignOutbound: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignOutbound(ctx, txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == zetacoreClient.Chain().ChainId {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Zeta:
			logger.Info().Msgf(
				"SignRevertTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			tx, err = signer.SignRevertTx(ctx, txData)
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignWithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignWithdrawTx(ctx, txData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf("SignERC20WithdrawTx: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				txData.gas.Price,
			)
			tx, err = signer.SignERC20WithdrawTx(ctx, txData)
		}
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		logger.Info().Msgf(
			"SignRevertTx: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.Price,
		)
		txData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		txData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)

		tx, err = signer.SignRevertTx(ctx, txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	} else if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
		logger.Info().Msgf(
			"SignOutbound: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			txData.gas.Price,
		)
		tx, err = signer.SignOutbound(ctx, txData)
		if err != nil {
			logger.Warn().Err(err).Msg(ErrorMsg(cctx))
			return
		}
	}

	logger.Info().Msgf(
		"Key-sign success: %d => %d, nonce %d",
		cctx.InboundParams.SenderChainId,
		toChain.ID(),
		cctx.GetCurrentOutboundParam().TssNonce,
	)

	// Broadcast Signed Tx
	signer.BroadcastOutbound(ctx, tx, cctx, logger, myID, zetacoreClient, txData)
}

// BroadcastOutbound signed transaction through evm rpc client
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *ethtypes.Transaction,
	cctx *types.CrossChainTx,
	logger zerolog.Logger,
	myID sdk.AccAddress,
	zetacoreClient interfaces.ZetacoreClient,
	txData *OutboundData,
) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Err(err).Msg("error getting app context")
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsZeta():
		// should not happen
		logger.Error().Msgf("unable to broadcast when toChain is zetaChain (%d)", toChain.ID())
		return
	case tx == nil:
		logger.Warn().Msgf("BroadcastOutbound: no tx to broadcast %s", cctx.Index)
		return
	}

	// broadcast transaction
	outboundHash := tx.Hash().Hex()

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)
		err := signer.broadcast(ctx, tx)
		if err != nil {
			log.Warn().
				Err(err).
				Msgf("BroadcastOutbound: error broadcasting tx %s on chain %d nonce %d retry %d signer %s",
					outboundHash, toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce, i, myID)
			retry, report := zetacore.HandleBroadcastError(
				err,
				strconv.FormatUint(cctx.GetCurrentOutboundParam().TssNonce, 10),
				fmt.Sprintf("%d", toChain.ID()),
				outboundHash,
			)
			if report {
				signer.reportToOutboundTracker(ctx, zetacoreClient, toChain.ID(), tx.Nonce(), outboundHash, logger)
			}
			if !retry {
				break
			}
			backOff *= 2
			continue
		}
		logger.Info().Msgf("BroadcastOutbound: broadcasted tx %s on chain %d nonce %d signer %s",
			outboundHash, toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce, myID)
		signer.reportToOutboundTracker(ctx, zetacoreClient, toChain.ID(), tx.Nonce(), outboundHash, logger)
		break // successful broadcast; no need to retry
	}
}

// SignERC20WithdrawTx
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *Signer) SignERC20WithdrawTx(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := signer.erc20CustodyABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdraw error: %w", err)
	}

	return tx, nil
}

// EvmClient returns the EVM RPC client
func (signer *Signer) EvmClient() interfaces.EVMRPCClient {
	return signer.client
}

// EvmSigner returns the EVM signer object for the signer
func (signer *Signer) EvmSigner() ethtypes.Signer {
	// TODO(revamp): rename field into evmSigner
	return signer.ethSigner
}

// IsSenderZetaChain checks if the sender chain is ZetaChain
// TODO(revamp): move to another package more general for cctx functions
func IsSenderZetaChain(
	cctx *types.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
) bool {
	return cctx.InboundParams.SenderChainId == zetacoreClient.Chain().ChainId &&
		cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound
}

// ErrorMsg returns a error message for SignOutbound failure with cctx data
func ErrorMsg(cctx *types.CrossChainTx) string {
	return fmt.Sprintf(
		"signer SignOutbound error: nonce %d chain %d",
		cctx.GetCurrentOutboundParam().TssNonce,
		cctx.GetCurrentOutboundParam().ReceiverChainId,
	)
}

// SignWhitelistERC20Cmd signs a whitelist command for ERC20 token
// TODO(revamp): move the cmd in a specific file
func (signer *Signer) SignWhitelistERC20Cmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
	outboundParams := txData.outboundParams
	erc20 := ethcommon.HexToAddress(params)
	if erc20 == (ethcommon.Address{}) {
		return nil, fmt.Errorf("SignCommandTx: invalid erc20 address %s", params)
	}

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	data, err := custodyAbi.Pack("whitelist", erc20)
	if err != nil {
		return nil, fmt.Errorf("whitelist pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gas,
		outboundParams.TssNonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign whitelist error: %w", err)
	}

	return tx, nil
}

// SignMigrateTssFundsCmd signs a migrate TSS funds command
// TODO(revamp): move the cmd in a specific file
func (signer *Signer) SignMigrateTssFundsCmd(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignMigrateTssFundsCmd error: %w", err)
	}
	return tx, nil
}

// reportToOutboundTracker reports outboundHash to tracker only when tx receipt is available
// TODO(revamp): move outbound tracker function to a outbound tracker file
func (signer *Signer) reportToOutboundTracker(
	ctx context.Context,
	zetacoreClient interfaces.ZetacoreClient,
	chainID int64,
	nonce uint64,
	outboundHash string,
	logger zerolog.Logger,
) {
	// set being reported flag to avoid duplicate reporting
	alreadySet := signer.Signer.SetBeingReportedFlag(outboundHash)
	if alreadySet {
		logger.Info().
			Msgf("reportToOutboundTracker: outboundHash %s for chain %d nonce %d is being reported", outboundHash, chainID, nonce)
		return
	}

	// report to outbound tracker with goroutine
	go func() {
		defer func() {
			signer.Signer.ClearBeingReportedFlag(outboundHash)
		}()

		// try monitoring tx inclusion status for 10 minutes
		var err error
		report := false
		isPending := false
		blockNumber := uint64(0)
		tStart := time.Now()
		for {
			// give up after 10 minutes of monitoring
			time.Sleep(10 * time.Second)

			if time.Since(tStart) > evm.OutboundInclusionTimeout {
				// if tx is still pending after timeout, report to outboundTracker anyway as we cannot monitor forever
				if isPending {
					report = true // probably will be included later
				}
				logger.Info().
					Msgf("reportToOutboundTracker: timeout waiting tx inclusion for chain %d nonce %d outboundHash %s report %v", chainID, nonce, outboundHash, report)
				break
			}
			// try getting the tx
			_, isPending, err = signer.client.TransactionByHash(ctx, ethcommon.HexToHash(outboundHash))
			if err != nil {
				logger.Info().
					Err(err).
					Msgf("reportToOutboundTracker: error getting tx for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				continue
			}
			// if tx is include in a block, try getting receipt
			if !isPending {
				report = true // included
				receipt, err := signer.client.TransactionReceipt(ctx, ethcommon.HexToHash(outboundHash))
				if err != nil {
					logger.Info().
						Err(err).
						Msgf("reportToOutboundTracker: error getting receipt for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				}
				if receipt != nil {
					blockNumber = receipt.BlockNumber.Uint64()
				}
				break
			}
			// keep monitoring pending tx
			logger.Info().
				Msgf("reportToOutboundTracker: tx has not been included yet for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
		}

		// try adding to outbound tracker for 10 minutes
		if report {
			tStart := time.Now()
			for {
				// give up after 10 minutes of retrying
				if time.Since(tStart) > evm.OutboundTrackerReportTimeout {
					logger.Info().
						Msgf("reportToOutboundTracker: timeout adding outbound tracker for chain %d nonce %d outboundHash %s, please add manually", chainID, nonce, outboundHash)
					break
				}
				// stop if the cctx is already finalized
				cctx, err := zetacoreClient.GetCctxByNonce(ctx, chainID, nonce)
				if err != nil {
					logger.Err(err).
						Msgf("reportToOutboundTracker: error getting cctx for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				} else if !crosschainkeeper.IsPending(cctx) {
					logger.Info().Msgf("reportToOutboundTracker: cctx already finalized for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
					break
				}
				// report to outbound tracker
				zetaHash, err := zetacoreClient.AddOutboundTracker(ctx, chainID, nonce, outboundHash, nil, "", -1)
				if err != nil {
					logger.Err(err).
						Msgf("reportToOutboundTracker: error adding to outbound tracker for chain %d nonce %d outboundHash %s", chainID, nonce, outboundHash)
				} else if zetaHash != "" {
					logger.Info().Msgf("reportToOutboundTracker: added outboundHash to core successful %s, chain %d nonce %d outboundHash %s block %d",
						zetaHash, chainID, nonce, outboundHash, blockNumber)
				} else {
					// stop if the tracker contains the outboundHash
					logger.Info().Msgf("reportToOutboundTracker: outbound tracker contains outboundHash %s for chain %d nonce %d", outboundHash, chainID, nonce)
					break
				}
				// retry otherwise
				time.Sleep(evm.ZetaBlockTime * 3)
			}
		}
	}()
}

// getEVMRPC is a helper function to set up the client and signer, also initializes a mock client for unit tests
func getEVMRPC(ctx context.Context, endpoint string) (interfaces.EVMRPCClient, ethtypes.Signer, error) {
	if endpoint == mocks.EVMRPCEnabled {
		chainID := big.NewInt(chains.BscMainnet.ChainId)
		ethSigner := ethtypes.NewLondonSigner(chainID)
		client := &mocks.MockEvmClient{}
		return client, ethSigner, nil
	}
	httpClient, err := metrics.GetInstrumentedHTTPClient(endpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get instrumented HTTP client")
	}

	rpcClient, err := ethrpc.DialHTTPWithClient(endpoint, httpClient)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to dial EVM client (endpoint %q)", endpoint)
	}
	client := ethclient.NewClient(rpcClient)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get chain ID")
	}

	ethSigner := ethtypes.LatestSignerForChainID(chainID)

	return client, ethSigner, nil
}
