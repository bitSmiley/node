package constant

const (
	// DonationMessage is the message for donation transactions
	// Transaction sent to the TSS or ERC20 Custody address containing this message are considered as a donation
	DonationMessage = "I am rich!"

	// CmdWhitelistERC20 is used for CCTX of type cmd to give the instruction to the TSS to whitelist an ERC20 on an exeternal chain
	CmdWhitelistERC20 = "cmd_whitelist_erc20"

	// CmdMigrateTssFunds is used for CCTX of type cmd to give the instruction to the TSS to transfer its funds on a new address
	CmdMigrateTssFunds = "cmd_migrate_tss_funds"

	// BTCWithdrawalDustAmount is the minimum satoshis that can be withdrawn from zEVM to avoid outbound dust output
	// The Bitcoin protocol sets a minimum output value to 546 satoshis (dust limit) but we set it to 1000 satoshis
	BTCWithdrawalDustAmount = 1000

	// SolanaWalletRentExempt is the minimum balance for a Solana wallet account to become rent exempt
	// The Solana protocol sets minimum rent exempt to 890880 lamports but we set it to 1_000_000 lamports (0.001 SOL)
	// The number 890880 comes from CLI command `solana rent 0` and has been verified on devnet gateway program
	SolanaWalletRentExempt = 1_000_000
)
