package main

import (
	"flag"
	"fmt"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

func main() {
	message := flag.String("message", "test message", "message text")

	settings := examples.NewSettings()
	settings.RequireTeam()
	settings.RequireChat()
	settings.RequireToken()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)
	client.SetVerbose(settings.Verbose)

	recipient := *tdproto.NewJID(settings.Chat)

	msg, err := client.SendPlaintextMessage(settings.TeamUid, recipient, *message)
	if err != nil {
		panic(err)
	}

	fmt.Printf("message created: %s\n", msg.Created)
}
