/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download a model from a remote registry",
	Long: `Download a model from a remote registry.

This will store the model in the local registry cache to be used later.`,
	PreRunE: preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO(gaocegege): Validate.

		return ormbClient.Pull(args[0], plainHTTPOpt)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().BoolVarP(&plainHTTPOpt, "plain-http", "", false, "use plain http and not https")
	pullCmd.Flags().BoolVarP(&insecureOpt, "insecure", "", true, "allow connections to TLS registry without certs")
}
