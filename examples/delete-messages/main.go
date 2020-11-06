package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
	"github.com/tada-team/tdproto/tdapi"
)

func main() {
	date := flag.String("date", "2017-12-31", "last date")

	settings := examples.NewSettings()
	settings.RequireToken()
	settings.RequireTeam()
	settings.RequireChat()
	settings.RequireDryRun()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)
	client.SetVerbose(settings.Verbose)

	chatUid := *tdproto.NewJID(settings.Chat)

	var lastMsgId = ""

	for {
		filter := new(tdapi.MessageFilter)
		filter.Lang = "ru"
		filter.Limit = 200
		filter.OldFrom = lastMsgId
		filter.Type = "change"
		filter.DateTo = *date
		messages, err := client.GetMessages(settings.TeamUid, chatUid, filter)
		if err != nil {
			panic(err)
		}
		if len(messages.Messages) != 0 {
			lastMsgId = getLastMessageId(messages)

			for key := range messages.Messages {
				if strings.HasPrefix(messages.Messages[key].PushText, "Удалён участник:") {
					if settings.DryRun {
						fmt.Println("message will be deleted (dryrun)", key, messages.Messages[key].PushText)
					} else {
						_, err := client.DeleteMessage(settings.TeamUid, chatUid, messages.Messages[key].MessageId)
						if err != nil {
							panic(err)
						}
						fmt.Println("Message deleted", key, messages.Messages[key].PushText)
					}
				}
			}
			if len(messages.Messages) < 200 {
				break
			}
		} else {
			break
		}
	}
}

func getLastMessageId(messages tdproto.ChatMessages) string {
	return messages.Messages[len(messages.Messages)-1].MessageId
}
