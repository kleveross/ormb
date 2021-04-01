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
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kleveross/ormb/pkg/ormb"
)

var cfgFile string
var logLevel uint32

var ormbClient ormb.Interface

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ormb",
	Short: "Manage Machine Learning/Deep Learning Models like Docker Images.",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithField("error", err).Panicln("Failed to run the  command")
	}
}

func init() {
	viper.SetEnvPrefix("ORMB")
	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ormb/config.yaml)")
	rootCmd.PersistentFlags().Uint32Var(&logLevel, "log-level", 4, "Log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLogger() {
	logrus.SetLevel(logrus.Level(logLevel))

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetReportCaller(false)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logrus.WithField("error", err).Panicln("Failed to find the home directory")
		}

		// Search config in home directory with name "config" (without extension).
		ormbHome := filepath.Join(home, ".ormb")
		viper.AddConfigPath(ormbHome)

		viper.SetConfigName("config")
		rootPath := viper.GetString("rootPath")

		if rootPath == "" {
			viper.SetDefault("rootPath", ormbHome)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.WithFields(logrus.Fields{
			"config": viper.ConfigFileUsed(),
		}).Debugln("Found the config file")
	}
}
