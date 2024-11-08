package main

import (
	"context"
	"fmt"
	"net"

	flag "github.com/spf13/pflag"
)

func main() {
	var nameServer net.IP
	var domainNames []string

	flag.IPVarP(&nameServer, "nameserver", "n", net.IP{}, "Specify a name server IP address")
	flag.StringSliceVarP(&domainNames, "domainnames", "d", []string{}, "Specify one or more domain names to resolve")

	flag.Parse()

	dns := net.JoinHostPort(nameServer.String(), "53")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial(network, dns)
		},
	}

	for _, domainName := range domainNames {
		addresses, err := resolver.LookupHost(context.Background(), domainName)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Domain name %s resolved to %s using %s as the name server\n", domainName, addresses[0], nameServer)
	}
}
