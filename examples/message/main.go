package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdproto"
)

func main() {
	server := flag.String("server", "https://web.tada.team", "server address")
	team := flag.String("team", "", "team uid")
	chat := flag.String("chat", "", "chat jid")
	token := flag.String("token", "", "bot token. Type \"/newbot <NAME>\" command in @TadaBot direct chat")
	message := flag.String("message", "test message", "message text")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	if *token == "" {
		fmt.Println("-token required")
		return
	}

	if *chat == "" {
		fmt.Println("-chat required")
		return
	}

	if *team == "" {
		fmt.Println("-team required")
		return
	}

	client, err := tdclient.NewSession(*server)
	if err != nil {
		panic(err)
	}

	client.SetToken(*token)
	client.SetVerbose(*verbose)

	websocketConnection, err := client.Ws(*team, nil)
	if err != nil {
		panic(err)
	}

	recipient := *tdproto.NewJID(*chat)

	// composing like human. Full events list at https://github.com/tada-team/tdproto
	websocketConnection.Send(tdproto.NewClientChatComposing(recipient, true, nil))
	time.Sleep(3 * time.Second)

	// shortcut for simple messaging
	websocketConnection.SendPlainMessage(recipient, *message)
	if _, err := websocketConnection.WaitForConfirm(); err != nil {
		panic(err)
	}
}
