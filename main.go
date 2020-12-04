package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

func main() {
	localPortPtr := flag.Int("local-port", 9090, "the local port that is opened")
	remoteEnabledPtr := flag.Bool("remote-enabled", false, "allow remotes to connect to your local port")
	var target string
	flag.StringVar(&target, "target", "157.230.121.173:9090", "target where the traff should be forwarded to")

	flag.Parse()

	var addressFormat = "127.0.0.1:%d"
	if *remoteEnabledPtr {
		addressFormat = ":%d"
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(addressFormat, *localPortPtr))
	if err != nil {
		panic(err)
	}

	fmt.Printf("started for port %d and target %s\n", *localPortPtr, target)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn, target)
	}
}

func handleRequest(conn net.Conn, target string) {
	fmt.Println("new client")

	proxy, err := net.Dial("tcp", target)
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
