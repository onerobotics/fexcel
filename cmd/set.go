package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set ./path/to/spreadsheet.xlsx ipAddress [more ipAddresses]",
	Short: "Set FANUC robots comments based on the provided Excel spreadsheet",
	Args:  validateSetArgs,
	RunE:  setMain,
}

func init() {
	rootCmd.AddCommand(setCmd)
}

func validateSetArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("requires a spreadsheet and at least one host IP address")
	}

	return nil
}

func setMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath, hosts := args[0], args[1:]

	setCmd, err := fexcel.NewSetCommand(fpath, globalCfg, hosts...)
	if err != nil {
		return err
	}

	startTime := time.Now()
	result, err := setCmd.Execute()
	// we will use err later

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)

	header := []string{""}
	sheetRow := []string{filepath.Base(fpath)}
	total := 0

	types := []fexcel.Type{fexcel.Ain, fexcel.Aout, fexcel.Din, fexcel.Dout, fexcel.Flag, fexcel.Gin, fexcel.Gout, fexcel.Numreg, fexcel.Posreg, fexcel.Rin, fexcel.Rout, fexcel.Sreg, fexcel.Ualm}
	for _, t := range types {
		header = append(header, t.String())
		defCount := len(setCmd.Definitions[t])
		total += defCount
		sheetRow = append(sheetRow, strconv.Itoa(defCount))
	}
	header = append(header, "Total")
	table.SetHeader(header)
	sheetRow = append(sheetRow, strconv.Itoa(total))
	table.Append(sheetRow)

	for _, host := range setCmd.Hosts() {
		total = 0
		row := []string{host}

		for _, t := range types {
			count := result.Counts[host][t]
			row = append(row, strconv.Itoa(count))
			total += count
		}

		row = append(row, strconv.Itoa(total))
		table.Append(row)
	}

	table.Render()

	fmt.Printf("Finished in %s.\n\n", time.Since(startTime))

	return err
}
