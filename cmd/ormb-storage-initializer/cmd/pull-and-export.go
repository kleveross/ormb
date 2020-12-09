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
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kleveross/ormb/pkg/consts"
	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"
	"github.com/kleveross/ormb/pkg/parser"
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

		// Pull the model from the remote registry.
		if err := ormbClient.Pull(modelURI, false); err != nil {
			return err
		}
		// Export it to the specified directory.
		if err := ormbClient.Export(modelURI, dstDir); err != nil {
			return err
		}

		// For model serving, need relayout option.
		if !reLayoutOpt {
			return nil
		}

		var relayoutFunc func(string) error

		isMLflow, err := checkForMLflow(dstDir)
		if err != nil {
			return err
		}

		if isMLflow {
			relayoutFunc = relayoutForMLflowModel
		} else {
			relayoutFunc = relayoutModel
		}

		if err := relayoutFunc(dstDir); err != nil {
			return err
		}

		return nil
	},
}

func checkForMLflow(modelDir string) (bool, error) {
	path, err := filepath.Abs(modelDir)
	if err != nil {
		return false, err
	}

	dat, err := ioutil.ReadFile(filepath.Join(path, consts.ORMBfileName))
	if err != nil {
		return false, err
	}

	metadata := &model.Metadata{}
	ormbParser := parser.NewDefaultParser()
	if metadata, err = ormbParser.Parse(dat); err != nil {
		return false, err
	}

	format := model.Format(metadata.Format)
	return format == model.FormatMLflow, nil
}

func relayoutForMLflowModel(modelDir string) error {
	// Relayout for pre-packaged MLflow server to serve
	// i.e. move /mnt/models/model/* to /mnt/models/* (dstDir).
	// so that, the pre-packaged MLflow server will serving, the refenence as
	// https://github.com/SeldonIO/seldon-core/blob/v1.1.0/operator/controllers/model_initializer_injector.go#L214
	originalDir, err := filepath.Abs(filepath.Join(modelDir, consts.ORMBModelDirectory))
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(originalDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Rename(path.Join(originalDir, file.Name()), path.Join(modelDir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

func relayoutModel(modelDir string) error {
	// Rename the `model` directory as `1`.
	// e.g. move /mnt/models/model to /mnt/models/1 (dstDir).
	// so that, the trtserver will serving, the refenence as
	// https://docs.nvidia.com/deeplearning/triton-inference-server/user-guide/docs/model_repository.html#section-onnx-models.
	originalDir, err := filepath.Abs(
		filepath.Join(modelDir, consts.ORMBModelDirectory))
	if err != nil {
		return err
	}

	if err := os.Rename(originalDir, path.Join(modelDir, "1")); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pullExportCmd)

	pullExportCmd.Flags().BoolVarP(&reLayoutOpt, "relayout", "", true, "relayout data for model serving")
	pullExportCmd.Flags().BoolVarP(&plainHTTPOpt, "plain-http", "", true, "use plain http and not https")
	pullExportCmd.Flags().BoolVarP(&insecureOpt, "insecure", "", true, "allow connections to TLS registry without certs")
}
