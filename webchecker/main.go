package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	hostname := flag.String("hostname", "", "Specify a hostname to check connectivity")
	timeout := flag.Duration("timeout", time.Second*10, "Specify timeout for the connectivity check")

	var port int
	flag.IntVar(&port, "port", 80, "Specify a port number to check connectivity")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <-hostname> [-port] [-timeout]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *hostname == "" {
		fmt.Println("-hostname argument must be provided with a non-empty value")
		os.Exit(1)
	}

	address := *hostname + ":" + strconv.Itoa(port)
	fmt.Printf("Trying to connect to: %s\n", address)

	conn, err := net.DialTimeout("tcp", address, *timeout)

	if err != nil {
		fmt.Println("Error connecting to", address)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Printf("Connected to %s with no errors\n", address)
}
