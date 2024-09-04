package main

import (
	"flag"
	"fmt"

	"github.com/u1f35c/grpc-test/client"
	"github.com/u1f35c/grpc-test/server"
)

var (
	port int
)

func main() {
	flag.IntVar(&port, "port", 12345, "TCP port to use")
	flag.Parse()

	args := flag.Args()

	var err error
	if len(args) != 1 {
		println("Must supply client or server as argument")
		return
	} else if args[0] == "grpcclient" {
		err = client.GRPCConnect(port)
	} else if args[0] == "grpcserver" {
		err = server.GRPCServe(port)
	} else if args[0] == "httpclient" {
		err = client.HTTP2Connect(port)
	} else if args[0] == "httpserver" {
		err = server.HTTP2Serve(port)
	} else {
		println("Unknown action:", args[0])
		return
	}

	if err != nil {
		fmt.Printf("Got error: %v\n", err)
	}
}
