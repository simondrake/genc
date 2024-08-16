package ip

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func newInCIDRCommand() *cobra.Command {
	var (
		ip   string
		cidr string
	)

	inCIDRCmd := &cobra.Command{
		Use:   "inCIDR",
		Short: "determine if an IP is within the range of a CIDR",
		Example: `
    $ genc ip inCIDR --ip "192.168.1.68" --cidr "192.168.1.0/24"
    IP (192.168.1.68) is in the CIDR range (192.168.1.0/24)

    $ genc ip inCIDR --ip "192.168.1.68" --cidr "192.168.1.30/32"
    IP (192.168.1.68) is not in the CIDR range (192.168.1.30/32)`,
		Run: func(cmd *cobra.Command, args []string) {
			ic, err := ipInCIDR(ip, cidr)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error determining if IP is in CIDR: %w", err))
				os.Exit(1)
			}

			if ic {
				fmt.Fprintf(os.Stdout, "IP (%s) is in the CIDR range (%s)", ip, cidr)
			} else {
				fmt.Fprintf(os.Stdout, "IP (%s) is not in the CIDR range (%s)", ip, cidr)
			}
		},
	}

	inCIDRCmd.Flags().StringVar(&ip, "ip", "", "the ip address")
	inCIDRCmd.Flags().StringVar(&cidr, "cidr", "", "the cidr")

	if err := inCIDRCmd.MarkFlagRequired("ip"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'ip' as required: %w", err))
	}
	if err := inCIDRCmd.MarkFlagRequired("cidr"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'cidr' as required: %w", err))
	}

	return inCIDRCmd
}

func ipInCIDR(ip, cidr string) (bool, error) {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return n.Contains(net.ParseIP(ip)), nil
}
