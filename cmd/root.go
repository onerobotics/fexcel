package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unreal/fexcel/fexcel"
)

const configFile = ".fexcel.yaml"

var (
	cfgFile   string
	save      bool
	globalCfg fexcel.Config
)

var rootCmd = &cobra.Command{
	Use:   "fexcel",
	Short: "Process a spreadsheet and report what fexcel sees",
	Args:  validateRootArgs,
	RunE:  rootMain,
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if save {
			fmt.Printf("saving flagset to config file... ")
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				_, err := os.Create(configFile)
				if err != nil {
					return err
				}
			}
			err := viper.WriteConfig()
			if err != nil {
				return err
			}
			fmt.Println("done!")
		}

		if !globalCfg.NoUpdate {
			var c fexcel.GitHubUpdateChecker
			err := c.UpdateCheck(os.Stdout)
			if err != nil {
				return errors.New("failed to get latest version id from GitHub.")
			}
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./"+configFile+")")
	rootCmd.PersistentFlags().BoolVarP(&save, "save", "", false, "save flagset to config file")
	rootCmd.PersistentFlags().BoolVar(&globalCfg.NoUpdate, "noupdate", false, "don't check for fexcel updates")

	rootCmd.PersistentFlags().IntVarP(&globalCfg.Timeout, "timeout", "", 5, "timeout value in seconds")

	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Sheet, "sheet", "Sheet1", "default sheet to look at when unspecified in the start cell")
	rootCmd.PersistentFlags().IntVar(&globalCfg.FileConfig.Offset, "offset", 1, "column offset between ids and comments")

	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Constants, "constants", "", "start cell of constant ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Numregs, "numregs", "", "start cell of numeric register ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Posregs, "posregs", "", "start cell of position register ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Sregs, "sregs", "", "start cell of string register ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Ualms, "ualms", "", "start cell of user alarm ids")

	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Ains, "ains", "", "start cell of analog input ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Aouts, "aouts", "", "start cell of analog output ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Dins, "dins", "", "start cell of digital input ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Douts, "douts", "", "start cell of digital output ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Flags, "flags", "", "start cell of flag ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Gins, "gins", "", "start cell of group input ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Gouts, "gouts", "", "start cell of group output ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Rins, "rins", "", "start cell of robot input ids")
	rootCmd.PersistentFlags().StringVar(&globalCfg.FileConfig.Routs, "routs", "", "start cell of robot output ids")

	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	viper.BindPFlag("fileconfig.sheet", rootCmd.PersistentFlags().Lookup("sheet"))
	viper.BindPFlag("fileconfig.offset", rootCmd.PersistentFlags().Lookup("offset"))

	viper.BindPFlag("fileconfig.numregs", rootCmd.PersistentFlags().Lookup("numregs"))
	viper.BindPFlag("fileconfig.posregs", rootCmd.PersistentFlags().Lookup("posregs"))
	viper.BindPFlag("fileconfig.sregs", rootCmd.PersistentFlags().Lookup("sregs"))
	viper.BindPFlag("fileconfig.ualms", rootCmd.PersistentFlags().Lookup("ualms"))

	viper.BindPFlag("fileconfig.ains", rootCmd.PersistentFlags().Lookup("ains"))
	viper.BindPFlag("fileconfig.aouts", rootCmd.PersistentFlags().Lookup("aouts"))
	viper.BindPFlag("fileconfig.dins", rootCmd.PersistentFlags().Lookup("dins"))
	viper.BindPFlag("fileconfig.douts", rootCmd.PersistentFlags().Lookup("douts"))
	viper.BindPFlag("fileconfig.flags", rootCmd.PersistentFlags().Lookup("flags"))
	viper.BindPFlag("fileconfig.gins", rootCmd.PersistentFlags().Lookup("gins"))
	viper.BindPFlag("fileconfig.gouts", rootCmd.PersistentFlags().Lookup("gouts"))
	viper.BindPFlag("fileconfig.rins", rootCmd.PersistentFlags().Lookup("rins"))
	viper.BindPFlag("fileconfig.routs", rootCmd.PersistentFlags().Lookup("routs"))
}

func initConfig() {
	ext := filepath.Ext(configFile)
	name := strings.TrimSuffix(configFile, ext)
	viper.SetConfigName(name)
	viper.SetConfigType(ext[1:]) // remove .
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		viper.Unmarshal(&globalCfg)
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config file not found, but that's ok
		} else {
			// config file found, but another error was produced.
			fmt.Println("Error reading config file: ", err)
		}
	}
}

func validateRootArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a spreadsheet")
	}

	ext := filepath.Ext(args[0])
	if ext != ".xlsx" {
		return errors.New("requires a .xlsx file generated by Excel 2007 or later")
	}

	return nil
}

func rootMain(cmd *cobra.Command, args []string) error {
	fmt.Printf(fexcel.Logo())

	fpath := args[0]

	f, err := fexcel.OpenFile(fpath, globalCfg.FileConfig)
	if err != nil {
		return err
	}

	if len(f.Locations) == 0 {
		fmt.Println("No location flags specified.")
		return nil
	}

	for d, _ := range f.Locations {
		defs, err := f.Definitions(d)
		if err != nil {
			return err
		}

		fmt.Printf("Found %d %ss.\n", len(defs), d)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
