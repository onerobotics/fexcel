package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff ./path/to/spreadsheet.xlsx target(s)...",
	Short: "Diff comments between your spreadsheet and remote/local targets",
	Args:  validateDiffArgs,
	RunE:  diffMain,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().IntVarP(&timeout, "timout", "", 500, "timeout value in milliseconds")
}

func validateDiffArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("requires a spreadsheet and at least one target (IP or backup directory)")
	}

	ext := filepath.Ext(args[0])
	if ext != ".xlsx" {
		return errors.New("requires a .xlsx file generated by Excel 2007 or later")
	}

	// TODO validate args[1:]?

	return nil
}

func diffMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath := args[0]

	d, err := fexcel.NewDiffCommand(fpath, globalCfg, timeout, args[1:]...)
	if err != nil {
		return err
	}

	err = d.Execute(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}