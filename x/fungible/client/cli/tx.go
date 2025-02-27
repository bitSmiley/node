package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdDeploySystemContracts(),
		CmdDeployFungibleCoinZRC4(),
		CmdRemoveForeignCoin(),
		CmdUpdateZRC20LiquidityCap(),
		CmdUpdateContractBytecode(),
		CmdUpdateSystemContract(),
		CmdPauseZRC20(),
		CmdUnpauseZRC20(),
		CmdUpdateZRC20WithdrawFee(),
		CmdUpdateGatewayContract(),
	)

	return cmd
}
