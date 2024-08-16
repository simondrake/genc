package cidr

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func newOverlapCommand() *cobra.Command {
	var cidrs []string

	overlapCmd := &cobra.Command{
		Use:   "overlap",
		Short: "determine if CIDR blocks overlap with each other",
		Example: `
    $ genc cidr overlap --cidrs 87.243.24.122/32,87.243.24.0/24
    CIDRs ([87.243.24.122/32 87.243.24.0/24]) do overlap

    $ genc cidr overlap --cidrs 87.243.24.122/32,87.243.25.0/24
    CIDRs do not overlap`,
		Run: func(cmd *cobra.Command, args []string) {
			overlap, c, err := doCIDRsOverlap(cidrs)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error determining if CIDRs overlap: %w", err))
				os.Exit(1)
			}

			if overlap {
				fmt.Fprintf(os.Stdout, "CIDRs (%s) do overlap", c)
			} else {
				fmt.Fprintln(os.Stdout, "CIDRs do not overlap")
			}
		},
	}

	overlapCmd.Flags().StringSliceVar(&cidrs, "cidrs", nil, "the list of CIDRs to check for overlap")

	if err := overlapCmd.MarkFlagRequired("cidrs"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'cidrs' as required: %w", err))
	}

	return overlapCmd
}

func doCIDRsOverlap(cidrs []string) (bool, []string, error) {
	for i1, c1 := range cidrs {
		_, n1, err := net.ParseCIDR(c1)
		if err != nil {
			return false, nil, err
		}

		for i2, c2 := range cidrs {
			if i1 == i2 {
				continue
			}

			_, n2, err := net.ParseCIDR(c2)
			if err != nil {
				return false, nil, err
			}

			if n1.Contains(n2.IP) || n2.Contains(n1.IP) {
				return true, []string{c1, c2}, nil
			}
		}

	}

	return false, nil, nil
}
