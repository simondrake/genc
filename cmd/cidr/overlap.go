package cidr

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func newOverlapCommand() *cobra.Command {
	var cidrs string

	overlapCmd := &cobra.Command{
		Use:   "overlap",
		Short: "determine if CIDR blocks overlap with each other",
		Example: `
    $ genc cidr overlap --cidrs '["87.243.24.122/32", "87.243.24.0/24"]'
    CIDRs (87.243.24.122/32) and (87.243.24.0/24) overlap

    $ genc cidr overlap --cidrs '["87.243.24.122/32", "87.243.25.0/24"]'
    CIDRs do not overlap`,
		Run: func(cmd *cobra.Command, args []string) {
			var c []string

			if err := json.Unmarshal([]byte(cidrs), &c); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error parsing cidrs: %w", err))
				os.Exit(1)
			}

			if err := doCIDRsOverlap(c); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error determining if CIDRs overlap: %w", err))
				os.Exit(1)
			}
		},
	}

	overlapCmd.Flags().StringVar(&cidrs, "cidrs", "", "the list of CIDRs to check for overlap")

	if err := overlapCmd.MarkFlagRequired("cidrs"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'cidrs' as required: %w", err))
	}

	return overlapCmd
}

func doCIDRsOverlap(cidrs []string) error {
	var overlaps bool

	// This map is a crufty way of ensuring that we don't print out the same
	// CIDR overlap twice (as we loop through each thing twice)
	//
	// Should find a more elegant way of doing this, but it'll do for now
	processed := make(map[string]struct{})

	for i1, c1 := range cidrs {
		_, n1, err := net.ParseCIDR(c1)
		if err != nil {
			return err
		}

		for i2, c2 := range cidrs {
			if i1 == i2 {
				continue
			}

			_, n2, err := net.ParseCIDR(c2)
			if err != nil {
				return err
			}

			_, c1e := processed[c1+c2]
			_, c2e := processed[c2+c1]

			if !c1e && !c2e && (n1.Contains(n2.IP) || n2.Contains(n1.IP)) {
				fmt.Fprintf(os.Stdout, "CIDRs (%s) and (%s) overlap\n", c1, c2)

				processed[c1+c2] = struct{}{}
				processed[c2+c1] = struct{}{}

				overlaps = true
			}
		}
	}

	if !overlaps {
		fmt.Fprintln(os.Stdout, "CIDRs do not overlap")
	}

	return nil
}
