package utils

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// RequireCCTXStatus checks if the cctx status is equal to the expected status
func RequireCCTXStatus(
	t require.TestingT,
	cctx *crosschaintypes.CrossChainTx,
	expected crosschaintypes.CctxStatus,
	msgAndArgs ...any,
) {
	msg := fmt.Sprintf("cctx status is not %q cctx index %s", expected.String(), cctx.Index)

	require.NotNil(t, cctx.CctxStatus)
	require.Equal(t, expected, cctx.CctxStatus.Status, msg+errSuffix(msgAndArgs...))
}

// RequireTxSuccessful checks if the receipt status is successful.
// Currently, it accepts eth receipt, but we can make this more generic by using type assertion.
func RequireTxSuccessful(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not successful"
	require.Equal(t, ethtypes.ReceiptStatusSuccessful, receipt.Status, msg+errSuffix(msgAndArgs...))
}

// RequiredTxFailed checks if the receipt status is failed.
// Currently, it accepts eth receipt, but we can make this more generic by using type assertion.
func RequiredTxFailed(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not successful"
	require.Equal(t, ethtypes.ReceiptStatusFailed, receipt.Status, msg+errSuffix(msgAndArgs...))
}

func errSuffix(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}

	template := "; " + msgAndArgs[0].(string)

	return fmt.Sprintf(template, msgAndArgs[1:])
}
