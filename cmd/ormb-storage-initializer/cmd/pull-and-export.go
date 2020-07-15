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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kleveross/ormb/pkg/consts"
	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"
)

// pullExportCmd represents the pull-and-export command.
var pullExportCmd = &cobra.Command{
	Use:     "pull-and-export",
	Short:   "Pull and export the model",
	Long:    ``,
	PreRunE: preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO(gaocegege): Check the args.
		modelURI := args[0]
		dstDir := args[1]

		// Get username and password from environment
		// Here AWS_SECRET_ACCESS_KEY and AWS_ACCESS_KEY_ID are used
		// because Seldon Core does not support renaming the environment variable name.
		username := viper.GetString("AWS_ACCESS_KEY_ID")
		pwd := viper.GetString("AWS_SECRET_ACCESS_KEY")
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
		)
		if err != nil {
			return err
		}

		// Pull the model from the remote registry.
		if err := ormbClient.Pull(modelURI); err != nil {
			return err
		}
		// Export it to the specified directory.
		if err := ormbClient.Export(modelURI, dstDir); err != nil {
			return err
		}

		// Move the files in model directory to the upper directory.
		// e.g. Move /mnt/models/model to /mnt/models (dstDir).
		// Seldon core will run `--model_base_path=dstDir` directly.
		originalDir, err := filepath.Abs(
			filepath.Join(dstDir, consts.ORMBModelDirectory))
		if err != nil {
			return err
		}
		destinationDir, err := filepath.Abs(dstDir)
		if err != nil {
			return err
		}
		files, err := ioutil.ReadDir(originalDir)
		if err != nil {
			return err
		}
		for _, f := range files {
			oldPath := filepath.Join(originalDir, f.Name())
			newPath := filepath.Join(destinationDir, f.Name())
			fmt.Printf("Moving %s to %s\n", oldPath, newPath)
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullExportCmd)
}
