package main

import (
	"fmt"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
)

func main() {
	settings := examples.NewSettings()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetVerbose(settings.Verbose)
	features, err := client.Features()
	if err != nil {
		panic(err)
	}

	fmt.Println("server version:", features.Build)
	fmt.Println("max message length:", features.MaxMessageLength)
}
