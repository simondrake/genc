package version

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "prints the version of the CLI",
		Run: func(cmd *cobra.Command, args []string) {
			bi, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Fprintln(os.Stderr, "Unable to determine version")
				os.Exit(1)
			}

			if bi.Main.Version != "" {
				fmt.Fprintf(os.Stdout, "Version: %s\n", bi.Main.Version)
			} else {
				fmt.Fprintln(os.Stdout, "Version: unknown")
			}
		},
	}

	return cmd
}
