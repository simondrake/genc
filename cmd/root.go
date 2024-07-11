package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/simondrake/genc/cmd/aesgcm"
	"github.com/simondrake/genc/cmd/jwt"
	"github.com/simondrake/genc/cmd/pkcs7"
	"github.com/simondrake/genc/cmd/rc4"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "genc",
		Short: "genc is a utility tool for encryption and decryption",
		Long:  "Perform common encryption and decryption operations",
	}

	rootCmd.AddCommand(pkcs7.NewCommand())
	rootCmd.AddCommand(aesgcm.NewCommand())
	rootCmd.AddCommand(rc4.NewCommand())
	rootCmd.AddCommand(jwt.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
