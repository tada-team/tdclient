package main

import (
	"flag"
	"time"

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

	websocketConnection, err := client.Ws(settings.TeamUid, nil)
	if err != nil {
		panic(err)
	}

	recipient := tdproto.JID(settings.Chat)

	// composing like human. Full events list at https://github.com/tada-team/tdproto
	websocketConnection.Send(tdproto.NewClientChatComposing(recipient, true, nil))
	time.Sleep(3 * time.Second)

	// shortcut for simple messaging
	websocketConnection.SendPlainMessage(recipient, *message)

	// stop composing
	websocketConnection.Send(tdproto.NewClientChatComposing(recipient, false, nil))
	time.Sleep(3 * time.Second)

	// stay online while message not sent
	if _, err := websocketConnection.WaitForConfirm(); err != nil {
		panic(err)
	}
}
