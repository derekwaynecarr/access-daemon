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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/eparis/access-daemon/client/util"
	"github.com/eparis/access-daemon/operations/ip"
)

func ipDoIt(cmd *cobra.Command, ipCmd string, cmdArgs []string) error {
	ipCA := ip.ClientArgs{
		Command: ipCmd,
		Args:    cmdArgs,
	}

	reader, err := util.WriteJSONGetStream(serverAddr, role, "ip", ipCA)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	return nil
}

func init() {
	ipNeighCmd := &cobra.Command{
		Use:     "neigh",
		Aliases: []string{"neighbour", "n"},
		RunE: func(cmd *cobra.Command, cmdArgs []string) error {
			return ipDoIt(cmd, "neigh", cmdArgs)
		},
	}

	ipAddrCmd := &cobra.Command{
		Use:     "address",
		Aliases: []string{"addr", "a"},
		RunE: func(cmd *cobra.Command, cmdArgs []string) error {
			return ipDoIt(cmd, "addr", cmdArgs)
		},
	}

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "Run ip on a remote machine",
		Long: `Run ip on a remote machine

client ip addr show
client ip neighbour`,
	}

	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipNeighCmd)
	ipCmd.AddCommand(ipAddrCmd)
}
