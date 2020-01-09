package cmd

import (
	"fmt"

	"github.com/onerobotics/fexcel/fexcel"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of fexcel",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("fexcel-v%s", fexcel.Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
