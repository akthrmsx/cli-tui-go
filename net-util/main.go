package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	port       int
	nameServer net.IP

	rootCmd = &cobra.Command{
		Use:   "net-util COMMAND ARGUMENTS [OPTIONS]",
		Short: "Perform network and connection checks",
		Long:  "A CLI tool for performing network connection and DNS checks",
		Args:  cobra.NoArgs,
	}

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Validate if an IP socket can be connected",
		Long:  "Validate if a provided IP address and port can be connected from the local machine",
		Args:  cobra.ExactArgs(1),
		RunE:  checkSocket,
	}

	resolveCmd = &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a domain name to an IPv4 address",
		Long:  "Resolve a domain name to an IPv4 address and, optionally, using a custom DNS server IP address",
		Args:  cobra.ExactArgs(1),
		RunE:  resolveName,
	}
)

func checkSocket(cmd *cobra.Command, args []string) error {
	address := fmt.Sprintf("%s:%s", args[0], strconv.Itoa(port))
	timeout, _ := time.ParseDuration("10s")
	_, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		fmt.Printf("Network address %s cannot be reached\n", address)
		return err
	} else {
		fmt.Printf("Network address %s is available\n", address)
		return nil
	}
}

func resolveName(cmd *cobra.Command, args []string) error {
	dns := net.JoinHostPort(nameServer.String(), "53")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial(network, dns)
		},
	}
	addresses, err := resolver.LookupHost(context.Background(), args[0])

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Domain name %s resolved to %s using %s as the name server\n", args[0], addresses[0], nameServer)
		return nil
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	checkCmd.Flags().IntVarP(&port, "port", "p", 80, "Port number to check")
	rootCmd.AddCommand(checkCmd)

	resolveCmd.Flags().IPVarP(&nameServer, "nameserver", "n", net.IPv4(8, 8, 8, 8), "Custom DNS server IP address")
	rootCmd.AddCommand(resolveCmd)
}
