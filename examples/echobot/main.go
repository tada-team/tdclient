package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tada-team/kozma"
	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

func main() {
	useReplyToField := flag.Bool("reply", false, "use reply_to field instead quote symbol")

	settings := examples.NewSettings()
	settings.RequireTeam()
	settings.RequireToken()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)
	client.SetVerbose(settings.Verbose)

	websocketConnection, err := client.Ws(settings.TeamUid, nil)
	if err != nil {
		panic(err)
	}

	me, err := client.Me(settings.TeamUid)
	if err != nil {
		panic(err)
	}

	log.Println("bot started:", me.DisplayName)
	for {
		message, delayed, err := websocketConnection.WaitForMessage()
		if err != nil {
			panic(err)
		}

		// skip message update
		if delayed {
			continue
		}

		// handle direct messages only
		if !message.ChatType.IsDirect() {
			continue
		}

		log.Println("got:", message.PushText)

		if *useReplyToField {
			websocketConnection.Send(tdproto.NewClientMessageUpdated(tdproto.ClientMessageUpdatedParams{
				To:      message.Chat,
				ReplyTo: message.MessageId,
				Content: tdproto.MessageContent{
					Type: tdproto.MediatypePlain,
					Text: kozma.Say(),
				},
			}))
		} else {
			reply := fmt.Sprintf("> %s\n%s", message.PushText, kozma.Say())
			websocketConnection.SendPlainMessage(message.Chat, reply)
		}
	}
}
