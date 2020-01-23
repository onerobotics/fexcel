package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unreal/fexcel/fexcel"
	"github.com/unreal/fexcel/fexcel/compile"
)

var compileCmd = &cobra.Command{
	Use:   "compile spreadsheet.xlsx filename",
	Short: "Compile a fexcel source file to a FANUC .ls file",
	Args:  validateCompileArgs,
	RunE:  compileMain,
}

var (
	o string
)

func init() {
	compileCmd.Flags().StringVarP(&o, "output", "o", "", "Output file (e.g. filename.ls)")
	rootCmd.AddCommand(compileCmd)
}

func validateCompileArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("requires a spreadsheet and a source filename")
	}

	return nil
}

func compileMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	xlspath, fpath := args[0], args[1]

	p, err := compile.NewPrinter(xlspath, globalCfg.FileConfig)
	if err != nil {
		return err
	}

	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	filename := filepath.Base(fpath)
	f, err := compile.Parse(filename, string(src))
	if err != nil {
		return err
	}

	err = p.Print(f)
	if err != nil {
		return err
	}

	if o != "" {
		err = ioutil.WriteFile(o, []byte(p.Output()), 0644)
		if err != nil {
			return err
		}

		fmt.Printf("Wrote output to %s\n", o)
	} else {
		fmt.Print(p.Output())
	}

	return nil
}
