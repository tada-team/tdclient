package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

var messageMap = make(map[string]string)
var blackValue = make(map[string]string)

func getLastMessageId(messages tdproto.ChatMessages) string {
	return messages.Messages[len(messages.Messages)-1].MessageId
}

func messageMapUpdate(messages tdproto.ChatMessages) {
	for key := range messages.Messages {
		messageMap[messages.Messages[key].MessageId] = messages.Messages[key].Content.Text
	}
}

func filterMessages() {
	for key := range messageMap {
		var text = messageMap[key]
		s := strings.Split(text, ":")[0]

		if s == "Удалён участник" {
			blackValue[key] = text
		}
	}
}

func main() {
	deph := flag.Int("deph", 5, "depth degree")

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

	messages, err := client.GetMessages(settings.TeamUid, chatUid)
	if err != nil {
		panic(err)
	}

	messageMapUpdate(messages)
	var lastMsgId = getLastMessageId(messages)

	for i := 0; i < *deph; i++ {
		fmt.Println("Загружаем страницу", i, lastMsgId)
		messagesOld, err := client.GetOldMessagesFrom(settings.TeamUid, chatUid, lastMsgId)
		if err != nil {
			panic(err)
		}
		fmt.Println("На странице", len(messagesOld.Messages))

		lastMsgId = getLastMessageId(messagesOld)
		messageMapUpdate(messagesOld)
	}

	fmt.Println("Сообщений всего загружено", len(messageMap))
	filterMessages()
	fmt.Println("Кандидатов на удаление", len(blackValue))

	if len(blackValue) > 0 {
		for key := range blackValue {
			if settings.DryRun {
				fmt.Println("сообщение будет удалено (dryrun)", key, blackValue[key])
			} else {
				_, err := client.DeleteMessage(settings.TeamUid, chatUid, key)
				if err != nil {
					panic(err)
				}
				fmt.Println("сообщение удалено", key, blackValue[key])
			}
		}
	} else {
		fmt.Println("Нет системных сообщений для удаления")
	}
}
