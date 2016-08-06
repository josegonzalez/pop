package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var cfgFile string
var env string
var debugMode bool
var version bool

var RootCmd = &cobra.Command{
	Use:   "soda",
	Short: "A tasty treat for all your database needs",
	PersistentPreRun: func(c *cobra.Command, args []string) {
		fmt.Printf("Soda v%s\n\n", Version)
		pop.Debug = debugMode
		env = defaults.String(os.Getenv("GO_ENV"), env)
		setConfigLocation()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !version {
			cmd.Help()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version information")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config/database.yml", "The configuration file you would like to use.")
	RootCmd.PersistentFlags().StringVarP(&env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	RootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Use debug/verbose mode")
}

func setConfigLocation() {
	abs, err := filepath.Abs(cfgFile)
	if err != nil {
		return
	}
	dir, file := filepath.Split(abs)
	pop.AddLookupPaths(dir)
	pop.ConfigName = file
}

func getConn() *pop.Connection {
	conn, err := pop.Connect(env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return conn
}