package main

import (
	"log"
	"time"

	"github.com/jou66jou/go-forky-blockchain/block"
	"github.com/jou66jou/go-forky-blockchain/service"
)

func main() {
	go func() {
		t := time.Now()
		genesisBlock := block.Block{0, t.String(), "", "", 0}
		// spew.Dump(genesisBlock)
		block.BCs = append(block.BCs, genesisBlock)
	}()
	log.Fatal(service.RunHTTP())
}
