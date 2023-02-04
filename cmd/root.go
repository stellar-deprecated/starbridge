package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/stellar/go/support/config"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/starbridge/app"
)

var RootCmd = &cobra.Command{
	Use:           "starbridge",
	Short:         "starbridge validator software",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			cfg     app.Config
			cfgPath = cmd.PersistentFlags().Lookup("conf").Value.String()
		)

		err := config.Read(cfgPath, &cfg)
		if err != nil {
			switch cause := errors.Cause(err).(type) {
			case *config.InvalidConfigError:
				return errors.Wrap(cause, "config file")
			default:
				return err
			}
		}

		cfg.WithdrawalWindow = time.Hour * 24
		app := app.NewApp(cfg)
		app.Run()
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().String("conf", "./starbridge.cfg", "config file path")
}
