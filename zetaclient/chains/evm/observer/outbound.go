package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchOutbound watches evm chain for outgoing txs status
// TODO(revamp): move ticker function to ticker file
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchOutbound(ctx context.Context) error {
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("EVM_WatchOutbound_%d", ob.Chain().ChainId),
		ob.GetChainParams().OutboundTicker,
	)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Msg("error creating ticker")
		return err
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ob.Logger().Outbound.Info().Msgf("WatchOutbound started for chain %d", ob.Chain().ChainId)
	sampledLogger := ob.Logger().Outbound.Sample(&zerolog.BasicSampler{N: 10})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !app.IsOutboundObservationEnabled() {
				sampledLogger.Info().
					Msgf("WatchOutbound: outbound observation is disabled for chain %d", ob.Chain().ChainId)
				continue
			}
			trackers, err := ob.ZetacoreClient().
				GetAllOutboundTrackerByChain(ctx, ob.Chain().ChainId, interfaces.Ascending)
			if err != nil {
				continue
			}
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
				if ob.IsTxConfirmed(nonceInt) { // Go to next tracker if this one already has a confirmed tx
					continue
				}
				txCount := 0
				var outboundReceipt *ethtypes.Receipt
				var outbound *ethtypes.Transaction
				for _, txHash := range tracker.HashList {
					if receipt, tx, ok := ob.checkConfirmedTx(ctx, txHash.TxHash, nonceInt); ok {
						txCount++
						outboundReceipt = receipt
						outbound = tx
						ob.Logger().Outbound.Info().
							Msgf("WatchOutbound: confirmed outbound %s for chain %d nonce %d", txHash.TxHash, ob.Chain().ChainId, nonceInt)
						if txCount > 1 {
							ob.Logger().Outbound.Error().Msgf(
								"WatchOutbound: checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v transaction %v", txCount, ob.Chain().ChainId, nonceInt, outboundReceipt, outbound)
						}
					}
				}
				if txCount == 1 { // should be only one txHash confirmed for each nonce.
					ob.SetTxNReceipt(nonceInt, outboundReceipt, outbound)
				} else if txCount > 1 { // should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.Logger().Outbound.Error().Msgf("WatchOutbound: confirmed multiple (%d) outbound for chain %d nonce %d", txCount, ob.Chain().ChainId, nonceInt)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.Logger().Outbound)
		case <-ob.StopChannel():
			ob.Logger().Outbound.Info().Msg("WatchOutbound: stopped")
			return nil
		}
	}
}

// PostVoteOutbound posts vote to zetacore for the confirmed outbound
func (ob *Observer) PostVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	receiveValue *big.Int,
	receiveStatus chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
	logger zerolog.Logger,
) {
	chainID := ob.Chain().ChainId

	signerAddress := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		signerAddress.String(),
		cctxIndex,
		receipt.TxHash.Hex(),
		receipt.BlockNumber.Uint64(),
		receipt.GasUsed,
		math.NewIntFromBigInt(transaction.GasPrice()),
		transaction.Gas(),
		math.NewUintFromBigInt(receiveValue),
		receiveStatus,
		chainID,
		nonce,
		coinType,
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit

	var retryGasLimit uint64
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	logFields := map[string]any{
		"chain":    chainID,
		"nonce":    nonce,
		"outbound": receipt.TxHash.String(),
	}
	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Fields(logFields).
			Msgf("PostVoteOutbound: error posting vote for chain %d", chainID)
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logFields["vote"] = zetaTxHash
		logFields["ballot"] = ballot
		logger.Info().Fields(logFields).Msgf("PostVoteOutbound: posted vote for chain %d", chainID)
	}
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, error) {
	// skip if outbound is not confirmed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	if !ob.IsTxConfirmed(nonce) {
		return true, nil
	}
	receipt, transaction := ob.GetTxNReceipt(nonce)
	sendID := fmt.Sprintf("%d-%d", ob.Chain().ChainId, nonce)
	logger := ob.Logger().Outbound.With().Str("sendID", sendID).Logger()

	// get connector and erce20Custody contracts
	connectorAddr, connector, err := ob.GetConnectorContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting zeta connector for chain %d", ob.Chain().ChainId)
	}
	custodyAddr, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting erc20 custody for chain %d", ob.Chain().ChainId)
	}

	// define a few common variables
	var receiveValue *big.Int
	var receiveStatus chains.ReceiveStatus
	cointype := cctx.InboundParams.CoinType

	// compliance check, special handling the cancelled cctx
	if compliance.IsCctxRestricted(cctx) {
		// use cctx's amount to bypass the amount check in zetacore
		receiveValue = cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_failed
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receiveStatus = chains.ReceiveStatus_success
		}
		ob.PostVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return false, nil
	}

	// parse the received value from the outbound receipt
	receiveValue, receiveStatus, err = ParseOutboundReceivedValue(
		cctx,
		receipt,
		transaction,
		cointype,
		connectorAddr,
		connector,
		custodyAddr,
		custody,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("VoteOutboundIfConfirmed: error parsing outbound event for chain %d txhash %s", ob.Chain().ChainId, receipt.TxHash)
		return true, err
	}

	// post vote to zetacore
	ob.PostVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
	return false, nil
}

// ParseAndCheckZetaEvent parses and checks ZetaReceived/ZetaReverted event from the outbound receipt
// It either returns an ZetaReceived or an ZetaReverted event, or an error if no event found
func ParseAndCheckZetaEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	connectorAddr ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
) (*zetaconnector.ZetaConnectorNonEthZetaReceived, *zetaconnector.ZetaConnectorNonEthZetaReverted, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		// try parsing ZetaReceived event
		received, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
		if err == nil {
			err = evm.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), evm.TopicsZetaReceived)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReceived event")
			}
			if !strings.EqualFold(received.DestinationAddress.Hex(), params.Receiver) {
				return nil, nil, fmt.Errorf("receiver address mismatch in ZetaReceived event, want %s got %s",
					params.Receiver, received.DestinationAddress.Hex())
			}
			if received.ZetaValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in ZetaReceived event, want %s got %s",
					params.Amount.String(), received.ZetaValue.String())
			}
			if ethcommon.BytesToHash(received.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in ZetaReceived event, want %s got %s",
					cctx.Index, hex.EncodeToString(received.InternalSendHash[:]))
			}
			return received, nil, nil
		}
		// try parsing ZetaReverted event
		reverted, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReverted(*vLog)
		if err == nil {
			err = evm.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), evm.TopicsZetaReverted)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReverted event")
			}
			if !strings.EqualFold(
				ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex(),
				cctx.InboundParams.Sender,
			) {
				return nil, nil, fmt.Errorf("receiver address mismatch in ZetaReverted event, want %s got %s",
					cctx.InboundParams.Sender, ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex())
			}
			if reverted.RemainingZetaValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in ZetaReverted event, want %s got %s",
					params.Amount.String(), reverted.RemainingZetaValue.String())
			}
			if ethcommon.BytesToHash(reverted.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in ZetaReverted event, want %s got %s",
					cctx.Index, hex.EncodeToString(reverted.InternalSendHash[:]))
			}
			return nil, reverted, nil
		}
	}
	return nil, nil, errors.New("no ZetaReceived/ZetaReverted event found")
}

// ParseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outbound receipt
func ParseAndCheckWithdrawnEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*erc20custody.ERC20CustodyWithdrawn, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ParseWithdrawn(*vLog)
		if err == nil {
			err = evm.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), evm.TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrap(err, "error validating Withdrawn event")
			}
			if !strings.EqualFold(withdrawn.Recipient.Hex(), params.Receiver) {
				return nil, fmt.Errorf("receiver address mismatch in Withdrawn event, want %s got %s",
					params.Receiver, withdrawn.Recipient.Hex())
			}
			if !strings.EqualFold(withdrawn.Asset.Hex(), cctx.InboundParams.Asset) {
				return nil, fmt.Errorf("asset mismatch in Withdrawn event, want %s got %s",
					cctx.InboundParams.Asset, withdrawn.Asset.Hex())
			}
			if withdrawn.Amount.Cmp(params.Amount.BigInt()) != 0 {
				return nil, fmt.Errorf("amount mismatch in Withdrawn event, want %s got %s",
					params.Amount.String(), withdrawn.Amount.String())
			}
			return withdrawn, nil
		}
	}
	return nil, errors.New("no ERC20 Withdrawn event found")
}

// ParseOutboundReceivedValue parses the received value and status from the outbound receipt
// The receivd value is the amount of Zeta/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
func ParseOutboundReceivedValue(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	cointype coin.CoinType,
	connectorAddress ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
	custodyAddress ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*big.Int, chains.ReceiveStatus, error) {
	// determine the receive status and value
	// https://docs.nethereum.com/en/latest/nethereum-receipt-status/
	receiveValue := big.NewInt(0)
	receiveStatus := chains.ReceiveStatus_failed
	if receipt.Status == ethtypes.ReceiptStatusSuccessful {
		receiveValue = transaction.Value()
		receiveStatus = chains.ReceiveStatus_success
	}

	// parse receive value from the outbound receipt for Zeta and ERC20
	switch cointype {
	case coin.CoinType_Zeta:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receivedLog, revertedLog, err := ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in ZetaReceived/ZetaReverted event for vote message
			if receivedLog != nil {
				receiveValue = receivedLog.ZetaValue
			} else if revertedLog != nil {
				receiveValue = revertedLog.RemainingZetaValue
			}
		}
	case coin.CoinType_ERC20:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			withdrawn, err := ParseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in Withdrawn event for vote message
			receiveValue = withdrawn.Amount
		}
	case coin.CoinType_Gas, coin.CoinType_Cmd:
		// nothing to do for CoinType_Gas/CoinType_Cmd, no need to parse event
	default:
		return nil, chains.ReceiveStatus_failed, fmt.Errorf("unknown coin type %s", cointype)
	}
	return receiveValue, receiveStatus, nil
}

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *Observer) checkConfirmedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", txHash).
			Int64("chainID", ob.Chain().ChainId).
			Msg("error getting transaction for outbound")
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		log.Error().
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", txHash).
			Uint64("nonce", nonce).
			Msg("transaction is nil for txHash")
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.Chain().ChainId))
	from, err := signer.Sender(transaction)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", transaction.Hash().Hex()).
			Int64("chainID", ob.Chain().ChainId).
			Msg("local recovery of sender address failed for outbound")
		return nil, nil, false
	}
	if from != ob.TSS().EVMAddress() { // must be TSS address
		// If from is not TSS address, check if it is one of the previous TSS addresses We can still try to confirm a tx which was broadcast by an old TSS
		// This is to handle situations where the outbound has already been broad-casted by an older TSS address and the zetacore is waiting for the all the required block confirmations
		// to go through before marking the cctx into a finalized state

		// TODO : improve this logic to verify that the correct TSS address is the from address.
		// https://github.com/zeta-chain/node/issues/2487
		log.Info().
			Str("function", "confirmTxByHash").
			Str("sender", from.Hex()).
			Str("outboundTxHash", transaction.Hash().Hex()).
			Int64("chainID", ob.Chain().ChainId).
			Str("currentTSSAddress", ob.TSS().EVMAddress().Hex()).
			Msg("sender is not current TSS address")
		addressList := ob.TSS().EVMAddressList()
		isOldTssAddress := false
		for _, addr := range addressList {
			if from == addr {
				isOldTssAddress = true
			}
		}
		if !isOldTssAddress {
			log.Error().
				Str("function", "confirmTxByHash").
				Str("sender", from.Hex()).
				Str("outboundTxHash", transaction.Hash().Hex()).
				Int64("chainID", ob.Chain().ChainId).
				Str("currentTSSAddress", ob.TSS().EVMAddress().Hex()).
				Msg("sender is not current or old TSS address")
			return nil, nil, false
		}
	}
	if transaction.Nonce() != nonce { // must match cctx nonce
		log.Error().
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", txHash).
			Uint64("wantedNonce", nonce).
			Uint64("gotTxNonce", transaction.Nonce()).
			Msg("outbound nonce mismatch")
		return nil, nil, false
	}

	// save pending transaction
	if isPending {
		ob.SetPendingTx(nonce, transaction)
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", txHash).
			Uint64("nonce", nonce).
			Msg("transactionReceipt error")
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		log.Error().
			Str("function", "confirmTxByHash").
			Str("outboundTxHash", txHash).
			Uint64("nonce", nonce).
			Msg("receipt is nil")
		return nil, nil, false
	}
	ob.LastBlock()
	// check confirmations
	lastHeight, err := ob.evmClient.BlockNumber(ctx)
	if err != nil {
		log.Error().
			Str("function", "confirmTxByHash").
			Err(err).
			Int64("chainID", ob.GetChainParams().ChainId).
			Msg("error getting block number for chain")
		return nil, nil, false
	}
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		log.Debug().
			Str("function", "confirmTxByHash").
			Str("txHash", txHash).
			Uint64("nonce", nonce).
			Uint64("receiptBlock", receipt.BlockNumber.Uint64()).
			Uint64("currentBlock", lastHeight).
			Msg("txHash included but not confirmed")
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.CheckTxInclusion(transaction, receipt)
	if err != nil {
		log.Error().
			Err(err).
			Str("function", "confirmTxByHash").
			Str("errorContext", "checkTxInclusion").
			Str("txHash", txHash).
			Uint64("nonce", nonce).
			Msg("checkTxInclusion error")
		return nil, nil, false
	}

	return receipt, transaction, true
}
