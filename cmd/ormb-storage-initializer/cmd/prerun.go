package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/caicloud/ormb/pkg/oci"
	"github.com/caicloud/ormb/pkg/ormb"
)

func preRunE(cmd *cobra.Command, args []string) error {
	rootPath, err := filepath.Abs(viper.GetString("rootPath"))
	if err != nil {
		return err
	}
	fmt.Printf("Using %s as the root path\n", rootPath)

	ormbClient, err = ormb.New(
		oci.ClientOptRootPath(rootPath),
		oci.ClientOptWriter(os.Stdout),
		oci.ClientOptPlainHTTP(plainHTTPOpt),
	)
	return err
}
