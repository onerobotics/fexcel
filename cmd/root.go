package cmd

import (
	"errors"
	"os"

	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var (
	globalCfg fexcel.Config
)

var rootCmd = &cobra.Command{
	Use:  "fexcel",
	Long: fexcel.Logo() + "fexcel lets you use Excel to manage your FANUC robot data",
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if !globalCfg.NoUpdate {
			err := fexcel.CheckForUpdates(os.Stdout)
			if err != nil {
				return errors.New("failed to get latest version id from GitHub.")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Numregs, "numregs", "", "", "start cell of numeric register ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Posregs, "posregs", "", "", "start cell of position register ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Ualms, "ualms", "", "", "start cell of user alarm ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Rins, "rins", "", "", "start cell of robot input ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Routs, "routs", "", "", "start cell of robot output ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Dins, "dins", "", "", "start cell of digital input ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Douts, "douts", "", "", "start cell of digital output ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Ains, "ains", "", "", "start cell of analog input ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Aouts, "aouts", "", "", "start cell of analog output ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Sregs, "sregs", "", "", "start cell of string register ids")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Flags, "flags", "", "", "start cell of flag ids")

	rootCmd.PersistentFlags().StringVarP(&globalCfg.Sheet, "sheet", "", "Sheet1", "default sheet to look at when unspecified in the start cell")
	rootCmd.PersistentFlags().IntVarP(&globalCfg.Offset, "offset", "", 1, "column offset between ids and comments")
	rootCmd.PersistentFlags().BoolVarP(&globalCfg.NoUpdate, "noupdate", "", false, "don't check for fexcel updates")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
