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
	"github.com/eparis/access-daemon/operations/cat"
)

func catDoIt(cmd *cobra.Command, cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("Must specify at least 1 file to cat")
	}

	catCA := cat.ClientArgs{
		Files: cmdArgs,
	}

	reader, err := util.WriteJSONGetStream(serverAddr, role, "cat", catCA)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	return nil
}

func init() {
	catCmd := &cobra.Command{
		Use:   "cat",
		Short: "Run cat on a remote machine",
		Long: `Run cat on a remote machine

client --role=all cat /etc/resolv.conf`,
		RunE: catDoIt,
	}

	rootCmd.AddCommand(catCmd)
}
