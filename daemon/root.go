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

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDir     = "/etc/access-daemon"
	cfgFile    string
	staticPath = "/static"
	bindAddr   = ":8080"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "A REST API daemon which provides role based operational access to machines",
	Long: `A REST API deamon which provides role based operational access to machines.
	
This should allow different roles within an operations or engineering team to get
root level access while restricting what they can do.

Each role comes with 'modules' which a user may call. Modules (if they support it) are configured
via files in /etc/accress-daemon/$role/$module/*

Running moudles looks like:
   curl -v -X GET -d '{"args": "--since=-1m -u=kernel -f"}' http://127.0.0.1:8080/all/journalctl
`,
	RunE: mainFunc,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", fmt.Sprintf("config file (default %s/config.yaml)", cfgDir))
	rootCmd.PersistentFlags().StringVar(&cfgDir, "config-dir", cfgDir, "config directory")

	rootCmd.PersistentFlags().StringVar(&staticPath, "static-path", staticPath, "path to static files served at /static")

	rootCmd.PersistentFlags().StringVar(&bindAddr, "bind-addr", bindAddr, "Address to bind")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/access-daemon")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
