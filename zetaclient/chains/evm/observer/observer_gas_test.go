package observer_test

import (
	"context"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

func TestPostGasPrice(t *testing.T) {
	const (
		gwei        = 10e9
		blockNumber = 1000
		anything    = mock.Anything
	)

	ctx := context.Background()

	t.Run("Pre EIP-1559 doesn't support priorityFee", func(t *testing.T) {
		// ARRANGE
		// Given ETH rpc mock
		ethRPC := mocks.NewMockEvmClient().WithBlockNumber(blockNumber)

		// Given zetacore client mock
		zetacoreClient := mocks.NewZetacoreClient(t).WithZetaChain()

		// Given an observer
		chain := chains.Ethereum
		confirmation := uint64(10)
		chainParam := mocks.MockChainParams(chain.ChainId, confirmation)

		observer, _ := MockEVMObserver(t, chain, ethRPC, nil, zetacoreClient, nil, blockNumber, chainParam)

		// Given empty baseFee from RPC
		ethRPC.WithHeader(&ethtypes.Header{BaseFee: nil})

		// Given gas price from RPC
		ethRPC.WithSuggestGasPrice(big.NewInt(3 * gwei))

		// Given mock collector for zetacore call
		// PostVoteGasPrice(ctx, chain, gasPrice, priorityFee, blockNum)
		var gasPrice, priorityFee uint64
		collector := func(args mock.Arguments) {
			gasPrice = args.Get(2).(uint64)
			priorityFee = args.Get(3).(uint64)
		}

		zetacoreClient.
			On("PostVoteGasPrice", anything, anything, anything, anything, anything).
			Run(collector).
			Return("0xABC123...", nil)

		// ACT
		err := observer.PostGasPrice(ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that gas price is posted with proper gasPrice and priorityFee
		assert.Equal(t, uint64(3*gwei), gasPrice)
		assert.Equal(t, uint64(0), priorityFee)
	})

	t.Run("Post EIP-1559 supports priorityFee", func(t *testing.T) {
		// ARRANGE
		// Given ETH rpc mock
		ethRPC := mocks.NewMockEvmClient().WithBlockNumber(blockNumber)

		// Given zetacore client mock
		zetacoreClient := mocks.NewZetacoreClient(t).WithZetaChain()

		// Given an observer
		chain := chains.Ethereum
		confirmation := uint64(10)
		chainParam := mocks.MockChainParams(chain.ChainId, confirmation)

		observer, _ := MockEVMObserver(t, chain, ethRPC, nil, zetacoreClient, nil, blockNumber, chainParam)

		// Given 1 gwei baseFee from RPC
		ethRPC.WithHeader(&ethtypes.Header{BaseFee: big.NewInt(gwei)})

		// Given gasPrice and priorityFee from RPC
		ethRPC.
			WithSuggestGasPrice(big.NewInt(3 * gwei)).
			WithSuggestGasTipCap(big.NewInt(2 * gwei))

		// Given mock collector for zetacore call
		// PostVoteGasPrice(ctx, chain, gasPrice, priorityFee, blockNum)
		var gasPrice, priorityFee uint64
		collector := func(args mock.Arguments) {
			gasPrice = args.Get(2).(uint64)
			priorityFee = args.Get(3).(uint64)
		}

		zetacoreClient.
			On("PostVoteGasPrice", anything, anything, anything, anything, anything).
			Run(collector).
			Return("0xABC123...", nil)

		// ACT
		err := observer.PostGasPrice(ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that gas price is posted with proper gasPrice and priorityFee
		assert.Equal(t, uint64(3*gwei), gasPrice)
		assert.Equal(t, uint64(2*gwei), priorityFee)
	})
}
