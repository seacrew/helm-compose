/*
Copyright © 2023 The Helm Compose Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/seacrew/helm-compose/internal/compose"
	"github.com/seacrew/helm-compose/internal/config"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [REVISION]",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if err := compose.CompatibleHelmVersion(); err != nil {
			return err
		}

		config, err := config.ParseConfigFile(composeFile)
		if err != nil {
			return err
		}

		revision, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("REVISION must be a number")
		}

		if revision < 1 {
			return fmt.Errorf("REVISION must be a positiv number")
		}

		return compose.GetRevision(revision, config)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
