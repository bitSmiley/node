package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UnpauseZRC20(t *testing.T) {
	t.Run("can unpause status of zrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		assertUnpaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.False(t, fc.Paused)
		}
		assertPaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.True(t, fc.Paused)
		}

		// setup zrc20
		zrc20A, zrc20B, zrc20C := sample.EthAddress().
			String(),
			sample.EthAddress().
				String(),
			sample.EthAddress().
				String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		fcB := sample.ForeignCoins(t, zrc20B)
		fcB.Paused = true
		k.SetForeignCoins(ctx, fcB)
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		// can unpause zrc20
		msg := types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UnpauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		// can unpause already unpaused zrc20
		msg = types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20C,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UnpauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		// can unpause all zrc20
		msg = types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UnpauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUnpauseZRC20(
			admin,
			[]string{sample.EthAddress().String()},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UnpauseZRC20(ctx, msg)

		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		zrc20A, zrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))

		msg := types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
				sample.EthAddress().String(),
				zrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UnpauseZRC20(ctx, msg)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
