package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jou66jou/go-forky-blockchain/p2p"
	"github.com/jou66jou/go-forky-blockchain/service"
)

var (
	port string
	seed string
)

func main() {
	initFlag()
	fmt.Println(port)
	p2p.MyPort = port
	if ("127.0.0.1:" + port) != seed { // 連上p2p節點
		p2p.ConnectionToAddr(seed, false)

	}
	log.Fatal(service.RunHTTP(port))
}

func initFlag() {
	flag.StringVar(&port, "p", "8080", "listen port")               // 8080
	flag.StringVar(&seed, "seed", "127.0.0.1:8080", "seed ip:port") // 127.0.0.1:8080
	flag.Parse()
}
