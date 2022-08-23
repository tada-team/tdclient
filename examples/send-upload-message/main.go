package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

func main() {
	filePath := flag.String("filepath", "", "path to file")

	settings := examples.NewSettings()
	settings.RequireTeam()
	settings.RequireChat()
	settings.RequireToken()
	settings.Parse()

	if *filePath == "" {
		fmt.Println("-filepath required")
		os.Exit(0)
	}

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)
	msg, err := client.SendUploadMessage(settings.TeamUid, tdproto.JID(settings.Chat), filepath.Base(*filePath), file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("message upload created: %s\n", msg.Created)
}
