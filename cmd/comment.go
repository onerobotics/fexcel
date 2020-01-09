package cmd

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/onerobotics/fexcel/fanuc"
	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var (
	timeout int
)

var commentCmd = &cobra.Command{
	Use:   "comment ./path/to/spreadsheet.xlsx ipAddress [more ipAddresses]",
	Short: "Set FANUC robot comments",
	Long:  "Set FANUC robots comments based on the provided Excel spreadsheet",
	Args:  validateArgs,
	RunE:  main,
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.Flags().IntVarP(&timeout, "timout", "", 500, "timeout value in milliseconds")
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("requires a spreadsheet and at least one host IP address")
	}

	ext := filepath.Ext(args[0])
	if ext != ".xlsx" {
		return errors.New("requires a .xlsx file generated by Excel 2007 or later")
	}

	for _, host := range args[1:] {
		if net.ParseIP(host) == nil {
			return fmt.Errorf("%s is not a valid IP address", host)
		}
	}

	return nil
}

func main(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	hosts := args[1:]

	f, err := fexcel.PrepareFile(args[0], globalCfg)
	if err != nil {
		return err
	}

	c := fanuc.NewMultiUpdater(hosts, &fanuc.CommentToolUpdater{time.Duration(timeout) * time.Millisecond})

	var definitions []fanuc.Definition
	for d, _ := range f.Locations {
		defs, err := f.Definitions(d)
		if err != nil {
			return err
		}

		fmt.Printf("Found %d %ss.\n", len(defs), d.VerboseName())

		definitions = append(definitions, defs...)
	}

	fmt.Printf("\nUpdating %d comments on %d %s... ", len(definitions), len(hosts), fexcel.Pluralize("host", len(hosts)))

	startTime := time.Now()

	err = c.Update(definitions)
	if err != nil {
		return err
	}

	fmt.Printf("finished in %s.\n\n", time.Since(startTime))

	for _, warning := range c.Warnings {
		fmt.Printf("[warning] %s\n", warning)
	}

	for host, errs := range c.Errors {
		for _, err := range errs {
			fmt.Printf("[error] %s: %s\n", host, err)
		}
	}

	if len(c.Errors) > 0 {
		return fmt.Errorf("Finished with %d errors", len(c.Errors))
	}

	return nil
}