package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:     "diff spreadsheet.xlsx target(s)...",
	Short:   "Compare robot comments to spreadsheet (remote or local)",
	Example: "  fexcel diff spreadsheet.xlsx 192.168.100.101 192.168.100.102 ./backup/dir ./some/other/backup/dir",
	Args:    validateDiffArgs,
	RunE:    diffMain,
}

var (
	all bool
)

func init() {
	diffCmd.Flags().BoolVar(&all, "all", false, "show all comparisons in summary tables instead of just differences")
	rootCmd.AddCommand(diffCmd)
}

func validateDiffArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("requires a spreadsheet and at least one target (IP or backup directory)")
	}

	// TODO validate args[1:]?

	return nil
}

func diffMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath := args[0]

	d, err := fexcel.NewDiffCommand(fpath, globalCfg, args[1:]...)
	if err != nil {
		return err
	}

	for dataType, _ := range d.Locations() {
		err := d.FprintTable(os.Stdout, dataType, all)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, "")
	}

	return nil
}
