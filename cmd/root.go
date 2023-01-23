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
	"os"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var composeFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "compose",
	Short: "Compose is a helm plugin to define and manage multiple helm releases as a single entity.",
	Long: `With the Helm Compose plugin you are able to create a single compose file
to manage a multitude of releases. Either for deploying one chart multiple times in
the same or different namespaces or to deploy multiple charts as a single entity together.

This idea is heavily inspired and influenced by the idea behind docker-compose.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&composeFile, "file", "f", "", "Compose configuration file")
}
