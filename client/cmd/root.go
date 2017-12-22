// Copyright © 2017 Eric Paris <eparis@redhat.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDir     string
	cfgFile    string
	serverAddr = "http://127.0.0.1:8080"
	role       = "all"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "A REST API client which provides role based operational access to machines",
	Long: `A REST API client which provides role based operational access to machines.
	
This should allow different roles within an operations or engineering team to get
root level access while restricting what they can do.

Each role comes with 'modules' which a user may call. Modules (if they support it) are configured
via files in /etc/accress-daemon/$role/$module/*

Running modules looks like:
   client get-roles // will return the list of all roles
   client --role=all get-modules // will return all modules for the 'all' role
   client bash rpm -q kernel // run the bash module in the default role for rpm -q kernel
   client --role=all journalctl --since=-10m -u kernel -u atomic-openshift-master-api -f // run the journalctl module to get logs for kernel and atomic-openshift-master-api
`,
	//RunE: func(cmd *cobra.Command, args []string) error {
	//return doIt(cmd, args)
	//return nil
	//},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", fmt.Sprintf("config file (default $HOME/.ops-client.yaml)"))
	rootCmd.PersistentFlags().StringVar(&cfgDir, "config-dir", cfgDir, "config directory (default $HOME)")
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", serverAddr, "URL of server")
	rootCmd.PersistentFlags().StringVar(&role, "role", role, "role of module we with to run")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if cfgDir == "" {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfgDir = home
		}

		// Search config in home directory with name ".crap" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName(".ops-client")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
