// Note: All credit for code in this command goes to https://github.com/vladimirvivien/go-networking/blob/master/ip/cidr/cidr.go

package cidr

import (
	"fmt"
	"math"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func newParseCommand() *cobra.Command {
	var cidr string

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parses the cidr and outputs relevant information",
		Example: `
$ genc cidr parse --cidr "192.168.1.0/24"

CIDR: 192.168.1.0/24
------------------------
Network:        192.168.1.0
IP Range:       192.168.1.0 - 192.168.1.255
Total Hosts:    256
Netmask:        255.255.255.0
Wildcard Mask:  0.0.0.255`,
		Run: func(cmd *cobra.Command, args []string) {
			pr, err := parseCIDR(cidr)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error parsing CIDR: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdin)
			fmt.Fprintf(os.Stdin, "CIDR: %s\n", pr.cidr)
			fmt.Fprintln(os.Stdin, "------------------------")
			fmt.Fprintf(os.Stdin, "Network:        %s\n", pr.network)
			fmt.Fprintf(os.Stdin, "IP Range:       %s\n", pr.ipRange)
			fmt.Fprintf(os.Stdin, "Total Hosts:    %s\n", pr.totalHosts)
			fmt.Fprintf(os.Stdin, "Netmask:        %s\n", pr.netmask)
			fmt.Fprintf(os.Stdin, "Wildcard Mask:  %s\n", pr.wildcardMask)
			fmt.Fprintln(os.Stdin)
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "the CIDR to parse")

	if err := cmd.MarkFlagRequired("cidr"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'cidr' as required: %w", err))
	}

	return cmd
}

type parseResponse struct {
	cidr         string
	network      string
	ipRange      string
	totalHosts   string
	netmask      string
	wildcardMask string
}

func parseCIDR(cidr string) (*parseResponse, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// Given IPv4 block 192.168.100.14/24
	// The followings uses IPNet to get:
	// - The routing address for the subnet (i.e. 192.168.100.0)
	// - one-bits of the network mask (24 out of 32 total)
	// - The subnetmask (i.e. 255.255.255.0)
	// - Total hosts in the network (2 ^(host identifer bits) or 2^8)
	// - Wildcard the inverse of subnet mask (i.e. 0.0.0.255)
	// - The maximum address of the subnet (i.e. 192.168.100.255)
	ones, totalBits := ipnet.Mask.Size()
	size := totalBits - ones                 // usable bits
	totalHosts := math.Pow(2, float64(size)) // 2^size
	wildcardIP := wildcard(net.IP(ipnet.Mask))
	last := lastIP(ip, net.IPMask(wildcardIP))

	return &parseResponse{
		cidr:         cidr,
		network:      ipnet.IP.String(),
		ipRange:      fmt.Sprintf("%s - %s", ipnet.IP, last),
		totalHosts:   fmt.Sprintf("%0.0f", totalHosts),
		netmask:      net.IP(ipnet.Mask).String(),
		wildcardMask: wildcardIP.String(),
	}, nil
}

// wildcard returns the opposite of the
// the netmask for the network.
func wildcard(mask net.IP) net.IP {
	var ipVal net.IP

	for _, octet := range mask {
		ipVal = append(ipVal, ^octet)
	}

	return ipVal
}

// lastIP calculates the highest addressable IP for given
// for a given subnet. It Loops through each octet of the
// subnet's IP address and applies bitwise OR operation
// to each corresponding octet from the mask value.
func lastIP(ip net.IP, mask net.IPMask) net.IP {
	ipIn := ip.To4() // is it an IPv4
	if ipIn == nil {
		ipIn = ip.To16() // is it IPv6
		if ipIn == nil {
			return nil
		}
	}
	var ipVal net.IP
	// apply network mask to each octet
	for i, octet := range ipIn {
		ipVal = append(ipVal, octet|mask[i])
	}
	return ipVal
}
