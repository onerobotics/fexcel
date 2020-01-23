package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unreal/fexcel/fexcel"
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
