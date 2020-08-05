package cmd

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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

	logrus.WithFields(logrus.Fields{
		"root-path": rootPath,
	}).Debugln("Create the ormb client with the given root path")

	ormbClient, err = ormb.New(
		oras.ClientOptRootPath(rootPath),
		oras.ClientOptWriter(os.Stdout),
		oras.ClientOptPlainHTTP(plainHTTPOpt),
		oras.ClientOptInsecure(insecureOpt),
	)
	return err
}
