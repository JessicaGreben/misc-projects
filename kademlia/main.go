package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const (
	idLength    = 20                     // Length in bytes of the Node ID.
	keyLength   = 20                     // Length in bytes of the key for storing data.
	k           = 20                     // Max number of contacts in any one bucket.
	alpha       = 3                      // System wide concurrency parameter.
	bitsPerByte = 8                      // How many bits in a byte.
	bucketCount = idLength * bitsPerByte // How many buckets should be created for each route table.
)

func main() {
	switch os.Args[1] {
	case "-s":
		port := os.Args[2]
		server(port)
	case "-c":
		Ping()
	}
}

func server(port string) error {
	network := new(Network)
	rpc.Register(network)
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return err
	}
	fmt.Println("Joining network...")

	err = network.join("localhost", port)
	if err != nil {
		fmt.Println("network.join err: ", err)
		return err
	}
	fmt.Printf("Serving RPC server on port %s", port)
	err = http.Serve(ln, nil)
	if err != nil {
		fmt.Println("http.Serve err: ", err)
		return err
	}
	return nil
}
