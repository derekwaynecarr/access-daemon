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
	"github.com/eparis/access-daemon/operations/command"
)

func commandDoIt(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Command operations require at least 1 argument.")
	}

	commandCA := command.ClientArgs{
		CmdName: args[0],
		Args:    args[1:],
	}

	reader, err := util.WriteJSONGetStream(serverAddr, role, "command", commandCA)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	return nil
}

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Run remote commands",
	Long: `Run remote commands. For example:

client --role=all command rpm -q kernel
`,
	RunE: commandDoIt,
}

func init() {
	rootCmd.AddCommand(commandCmd)
	commandCmd.Flags().SetInterspersed(false)
}
