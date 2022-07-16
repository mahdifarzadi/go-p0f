package main

import (
	"flag"
	"fmt"
	gop0f "github.com/mahdifarzadi/go-p0f"
	"log"
)

func main() {
	// parse flags
	ip := ""
	socketAddress := ""
	flag.StringVar(&ip, "i", "127.0.0.1", "ip to query")
	flag.StringVar(&socketAddress, "s", "p0f.socket", "path to p0f socket")
	flag.Parse()

	// create a client and connect to p0f
	p0f, err := gop0f.New(socketAddress)
	if err != nil {
		panic(err)
	}

	// send a query (device's ip)
	res, err := p0f.Query(ip)
	if err != nil {
		log.Println(err)
	}

	// print device's os name and flavor
	fmt.Println("OS:", res.OsName, res.OsFlavor)
}
