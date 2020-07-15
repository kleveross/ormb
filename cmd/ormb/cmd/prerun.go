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
	)
	return err
}
