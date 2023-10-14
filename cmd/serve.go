// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//nolint:gochecknoglobals,gochecknoinits // common pattern when using cobra library
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ccamel/playground-protoactor.go/internal/system"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the protoactor platform",
	Long:  `Start the protoactor platform`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sys, err := system.Boot()
		if err != nil {
			return err
		}

		sys.Wait()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
