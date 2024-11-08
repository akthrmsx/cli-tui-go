package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	var ip net.IP
	var port int

	flag.IPVarP(&ip, "ip", "i", net.IPv4(192, 168, 1, 1), "Specify an IP address to validate")
	flag.IntVarP(&port, "port", "p", 80, "Specify a port to check connectivity")

	flag.Parse()

	address := fmt.Sprintf("%s:%s", ip.String(), strconv.Itoa(port))
	timeout, _ := time.ParseDuration("10s")
	_, err := net.DialTimeout("ipv4::Tcp", address, timeout)

	if err != nil {
		fmt.Printf("Network address %s cannot be reached\n", address)
	} else {
		fmt.Printf("Network address %s is available\n", address)
	}
}
