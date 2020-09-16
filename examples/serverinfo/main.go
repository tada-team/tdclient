package main

import (
	"flag"
	"fmt"

	"github.com/tada-team/tdclient"
)

func main() {
	server := flag.String("server", "https://web.tada.team", "server address")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	client, err := tdclient.NewSession(*server)
	if err != nil {
		panic(err)
	}

	client.SetVerbose(*verbose)
	features, err := client.Features()
	if err != nil {
		panic(err)
	}

	fmt.Println("server version:", features.Build)
}
