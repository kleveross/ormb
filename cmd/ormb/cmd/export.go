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
	"github.com/caicloud/ormb/pkg/ormb"
	"github.com/spf13/cobra"
)

var destination string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a model stored in local registry cache",
	Long: `Export a model stored in local registry cache.

This will create a new directory with the name of
the model, in a format that developers can modify
and check into source control if desired.`,
	Run: func(cmd *cobra.Command, args []string) {
		o, err := ormb.NewDefaultOCIormb()
		if err != nil {
			panic(err)
		}

		if err := o.Export(args[0], destination); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&destination, "destination", "d", ".", "location to write the chart.")
}
