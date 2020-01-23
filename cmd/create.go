package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unreal/fexcel/fexcel"
)

var createCmd = &cobra.Command{
	Use:     "create spreadsheet.xlsx target",
	Short:   "Create a spreadsheet based on a target's comments",
	Example: "  fexcel create ./doc/spreadsheet.xlsx 192.168.100.101",
	Args:    validateCreateArgs,
	RunE:    createMain,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func validateCreateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("requires a spreadsheet path and a target (IP or backup directory)")
	}

	ext := filepath.Ext(args[0])
	if ext != ".xlsx" {
		return errors.New("requires a spreadsheet path ending in .xlsx")
	}

	// TODO: validate args[1]?

	return nil
}

func createMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath, targetPath := args[0], args[1]

	c, err := fexcel.NewCreator(fpath, globalCfg, targetPath)
	if err != nil {
		return err
	}

	err = c.Create(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
