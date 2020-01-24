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

var (
	headers  bool
	template bool
)

func init() {
	createCmd.Flags().BoolVar(&headers, "headers", false, "write column header names")
	createCmd.Flags().BoolVar(&template, "template", false, "ignore config file and cell specs flags; use fexcel default template instead")
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

func templateConfig() fexcel.FileConfig {
	return fexcel.FileConfig{
		Numregs: "Sheet1:A2",
		Posregs: "Sheet1:D2",
		Flags:   "Sheet1:G2",
		Sregs:   "Sheet1:J2",
		Dins:    "A2",
		Douts:   "D2",
		Gins:    "G2",
		Gouts:   "J2",
		Rins:    "M2",
		Routs:   "P2",
		Ains:    "S2",
		Aouts:   "V2",
		Ualms:   "Alarms:A2",
		Sheet:   "IO",
		Offset:  1,
	}
}

func createMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath, targetPath := args[0], args[1]

	if template {
		globalCfg.FileConfig = templateConfig()
	}

	c, err := fexcel.NewCreator(fpath, globalCfg, headers, targetPath)
	if err != nil {
		return err
	}

	err = c.Create(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
