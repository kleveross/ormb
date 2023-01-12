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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"
)

// pushExportCmd represents the save-and-push command.
var pushExportCmd = &cobra.Command{
	Use:     "save-and-push",
	Short:   "save and push the model",
	Long:    ``,
	PreRunE: preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO(gaocegege): Check the args.
		modelURI := args[0]
		dstDir := args[1]

		// Get username and password from environment
		username := viper.GetString("ORMB_USERNAME")
		pwd := viper.GetString("ORMB_PASSWORD")
		// Get the host from the URL.
		strs := strings.Split(modelURI, "/")
		if len(strs) == 0 {
			return fmt.Errorf("Failed to get the host from %s", modelURI)
		}
		fmt.Printf("Logging to the remote registry %s\n", strs[0])
		fmt.Printf("Username: %s\n", username)
		if err := ormbClient.Login(strs[0], username, pwd, true); err != nil {
			return err
		}

		// Recreate the ORMB client to let it know the registry and config.
		rootPath, err := filepath.Abs(viper.GetString("rootPath"))
		if err != nil {
			return err
		}
		fmt.Printf("Using %s as the root path\n", rootPath)

		ormbClient, err = ormb.New(
			oras.ClientOptRootPath(rootPath),
			oras.ClientOptWriter(os.Stdout),
			oras.ClientOptPlainHTTP(plainHTTPOpt),
			oras.ClientOptInsecure(insecureOpt),
		)
		if err != nil {
			return err
		}

		// Save the model to local cache
		if err := ormbClient.Save(dstDir, modelURI); err != nil {
			return err
		}
		// push model for remote hub
		if err := ormbClient.Push(modelURI); err != nil {
			return err
		}

		if err := ormbClient.Remove(modelURI); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushExportCmd)
	pushExportCmd.Flags().BoolVarP(&plainHTTPOpt, "plain-http", "", true, "use plain http and not https")
	pushExportCmd.Flags().BoolVarP(&insecureOpt, "insecure", "", true, "allow connections to TLS registry without certs")
}
