package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// print balanceAfter of from address
	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Info("from address %s balance of SOL before: %d", r.ZEVMAuth.From, balanceBefore)

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	// #nosec G115 e2e - always in range
	withdrawAmount := big.NewInt(int64(parseInt(r, args[0])))
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)

	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// withdraw
	r.WithdrawSOLZRC20(privkey.PublicKey(), withdrawAmount, approvedAmount)

	// print balance of from address after withdraw
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Info("from address %s balance of SOL after: %d", r.ZEVMAuth.From, balanceAfter)

	// check if the balance is reduced correctly
	amountReduced := new(big.Int).Sub(balanceBefore, balanceAfter)
	require.True(r, amountReduced.Cmp(withdrawAmount) >= 0, "balance is not reduced correctly")
}
