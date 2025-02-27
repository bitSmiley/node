package app

import (
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"golang.org/x/exp/slices"

	"github.com/zeta-chain/zetacore/pkg/constant"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func SetupHandlers(app *App) {
	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case emissionstypes.ModuleName:
			keyTable = emissionstypes.ParamKeyTable()
		default:
			continue
		}
		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}
	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	needsForcedMigration := []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		emissionstypes.ModuleName,
		authoritytypes.ModuleName,
		observertypes.ModuleName,
	}
	allUpgrades := upgradeTracker{
		upgrades: []upgradeTrackerItem{
			{
				index: 1714664193,
				storeUpgrade: &storetypes.StoreUpgrades{
					Added: []string{consensustypes.ModuleName, crisistypes.ModuleName},
				},
				upgradeHandler: func(ctx sdk.Context, vm module.VersionMap) (module.VersionMap, error) {
					// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module
					// https://docs.cosmos.network/main/build/migrations/upgrading#xconsensus
					baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

					// empty version map happens when upgrading from old versions which did not correctly call
					// app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()) in InitChainer.
					// we must populate the version map if we detect this scenario
					//
					// this will only happen on the first upgrade. mainnet and testnet will not require this condition.
					if len(vm) == 0 {
						for m, mb := range app.mm.Modules {
							if module, ok := mb.(module.HasConsensusVersion); ok {
								if slices.Contains(needsForcedMigration, m) {
									vm[m] = module.ConsensusVersion() - 1
								} else {
									vm[m] = module.ConsensusVersion()
								}
							}
						}
					}
					return vm, nil
				},
			},
			// TODO: enable back IBC
			// these commented lines allow for the IBC modules to be added to the upgrade tracker
			// https://github.com/zeta-chain/node/issues/2573
			//{
			//	index: 1715624665,
			//	storeUpgrade: &storetypes.StoreUpgrades{
			//		Added: []string{
			//			capabilitytypes.ModuleName,
			//			ibcexported.ModuleName,
			//			ibctransfertypes.ModuleName,
			//		},
			//	},
			//},
			//{
			//	index: 1715707436,
			//	storeUpgrade: &storetypes.StoreUpgrades{
			//		Added: []string{ibccrosschaintypes.ModuleName},
			//	},
			//},
		},
		stateFileDir: DefaultNodeHome,
	}

	var upgradeHandlerFns []upgradeHandlerFn
	var storeUpgrades *storetypes.StoreUpgrades
	var err error
	_, useIncrementalTracker := os.LookupEnv("ZETACORED_USE_INCREMENTAL_UPGRADE_TRACKER")
	if useIncrementalTracker {
		upgradeHandlerFns, storeUpgrades, err = allUpgrades.getIncrementalUpgrades()
		if err != nil {
			panic(err)
		}
	} else {
		upgradeHandlerFns, storeUpgrades = allUpgrades.mergeAllUpgrades()
	}

	app.UpgradeKeeper.SetUpgradeHandler(
		constant.Version,
		func(ctx sdk.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + constant.Version)

			var err error
			for _, upgradeHandler := range upgradeHandlerFns {
				vm, err = upgradeHandler(ctx, vm)
				if err != nil {
					return vm, err
				}
			}

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == constant.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
