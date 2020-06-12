/*
Copyright © 2020 Caicloud Authors

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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/caicloud/ormb/pkg/consts"
	"github.com/caicloud/ormb/pkg/oci"
	"github.com/caicloud/ormb/pkg/ormb"
)

var ormbClient ormb.ORMB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ormb-storage-initializer",
	Short:   "Download model from remote registry in Seldon Core and KFSerivng",
	Long:    ``,
	PreRunE: preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		ormbClient, err = ormb.NewOCIORMB(
			oci.ClientOptRootPath(rootPath),
			oci.ClientOptWriter(os.Stdout),
			oci.ClientOptPlainHTTP(plainHTTPOpt),
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.AutomaticEnv()

	logrus.SetLevel(logrus.DebugLevel)
}
