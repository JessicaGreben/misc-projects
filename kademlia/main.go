package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	kadNet "github.com/jessicagreben/kademlia/pkg/network"
	"github.com/jessicagreben/kademlia/pkg/types"
)

func main() {
	switch os.Args[1] {
	case "-s":
		port := os.Args[2]
		server(port)
	case "-c":
		boot := types.Contact{
			NodeID: types.NodeID{0},
			IP:     "boot",
			Port:   "8080",
		}

		if err := kadNet.Ping(boot); err != nil {
			os.Exit(1)
		}
	}
}

func server(port string) error {
	network := new(kadNet.Network)

	// Publishes the networks methods to the server.
	rpc.Register(network)

	// Registers an HTTP handler for RPC messages to the server.
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return err
	}
	fmt.Println("Joining network...")

	err = network.Join("localhost", port)
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
