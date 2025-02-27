package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// migrationTestRoutine runs migration related e2e tests
func migrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserMigration
		// initialize runner for migration test
		migrationTestRunner, err := initTestRunner(
			"migration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, "migration"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		migrationTestRunner.Logger.Print("🏃 starting migration tests")
		startTime := time.Now()

		if len(testNames) == 0 {
			migrationTestRunner.Logger.Print("🍾 migration tests completed in %s", time.Since(startTime).String())
			return nil
		}
		// run migration test
		testsToRun, err := migrationTestRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("migration tests failed: %v", err)
		}

		if err := migrationTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("migration tests failed: %v", err)
		}
		if err := migrationTestRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		migrationTestRunner.Logger.Print("🍾 migration tests completed in %s", time.Since(startTime).String())

		return err
	}
}
