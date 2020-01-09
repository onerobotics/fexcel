package cmd

import (
	"errors"
	"os"

	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var (
	globalFlags = struct {
		Numregs string // e.g. A2 or Sheet1:A2
		Posregs string
		Ualms   string
		Rins    string
		Routs   string
		Dins    string
		Douts   string
		Gins    string
		Gouts   string
		Ains    string
		Aouts   string
		Sregs   string
		Flags   string

		Sheet    string
		Offset   int
		NoUpdate bool
	}{}
)

var rootCmd = &cobra.Command{
	Use:  "fexcel",
	Long: fexcel.Logo() + "fexcel lets you use Excel to manage your FANUC robot data",
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if !globalFlags.NoUpdate {
			err := fexcel.CheckForUpdates(os.Stdout)
			if err != nil {
				return errors.New("failed to get latest version id from GitHub.")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Numregs, "numregs", "", "", "start cell of numeric register ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Posregs, "posregs", "", "", "start cell of position register ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Ualms, "ualms", "", "", "start cell of user alarm ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Rins, "rins", "", "", "start cell of robot input ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Routs, "routs", "", "", "start cell of robot output ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Dins, "dins", "", "", "start cell of digital input ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Douts, "douts", "", "", "start cell of digital output ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Ains, "ains", "", "", "start cell of analog input ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Aouts, "aouts", "", "", "start cell of analog output ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Sregs, "sregs", "", "", "start cell of string register ids")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Flags, "flags", "", "", "start cell of flag ids")

	rootCmd.PersistentFlags().StringVarP(&globalFlags.Sheet, "sheet", "", "Sheet1", "default sheet to look at when unspecified in the start cell")
	rootCmd.PersistentFlags().IntVarP(&globalFlags.Offset, "offset", "", 1, "column offset between ids and comments")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.NoUpdate, "noupdate", "", false, "don't check for fexcel updates")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
