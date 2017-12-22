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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/eparis/access-daemon/client/util"
	"github.com/eparis/access-daemon/operations/journalctl"
)

var (
	journalArgs = journalctl.ClientArgs{}
)

func journalDoIt(cmd *cobra.Command, cmdArgs []string) error {
	if len(cmdArgs) != 0 {
		return fmt.Errorf("Unsupported arguments: %v", cmdArgs)
	}

	reader, err := util.WriteJSONGetStream(serverAddr, role, "journalctl", journalArgs)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	return nil
}

// journalctlCmd represents the journalctl command
var journalctlCmd = &cobra.Command{
	Use:   "journalctl",
	Short: "Run journalctl on a remote machine",
	Long: `Run journalctl on a remote machine

client --role=all journalctl --since=-1h -u atomic-openshift-master-api -u atomic-openshift-master-controllers -u atomic-openshift-node -u etcd`,
	RunE: journalDoIt,
}

func init() {
	rootCmd.AddCommand(journalctlCmd)
	journalctlCmd.Flags().StringVar(&journalArgs.Since, "since", journalArgs.Since, "Start showing entries on or newer than the specified date")
	journalctlCmd.Flags().StringVar(&journalArgs.Until, "until", journalArgs.Until, "Start showing entries on or older than the specified date")
	journalctlCmd.Flags().BoolVarP(&journalArgs.Follow, "follow", "f", journalArgs.Follow, "follow the logs as they are generated")
	journalctlCmd.Flags().StringSliceVarP(&journalArgs.Units, "unit", "u", journalArgs.Units, "Show messages for the specified systemd unit")
}
