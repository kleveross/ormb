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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kleveross/ormb/pkg/oras"
	"github.com/kleveross/ormb/pkg/ormb"
)

func preRunE(cmd *cobra.Command, args []string) error {
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
	return err
}
