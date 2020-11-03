package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

func main() {

	result_pull := make(map[string]string)
	all_msg := make(map[string]string)

	settings := examples.NewSettings()
	settings.RequireToken()
	settings.RequireTeam()
	settings.RequireChat()
	settings.RequireDryRun()
	settings.RequireDeep()
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

	for key := range messages.Messages {
		all_msg[messages.Messages[key].MessageId] = messages.Messages[key].Content.Text
	}

	var lastMsgId = messages.Messages[len(messages.Messages)-1].MessageId

	for i := 1; i < settings.Deep; i++ {
		messagesOld, err := client.GetMessagesOldMsg(settings.TeamUid, chatUid, lastMsgId)
		if err != nil {
			panic(err)
		}
		lastMsgId = messages.Messages[len(messagesOld.Messages)-1].MessageId
		for key := range messagesOld.Messages {
			all_msg[messagesOld.Messages[key].MessageId] = messagesOld.Messages[key].Content.Text
		}
	}

	for key := range all_msg {
		var text = all_msg[key]
		s := strings.Split(text, ":")[0]

		if s == "Удалён участник" {
			result_pull[key] = text
		}

	}

	if settings.DryRun {
		if len(result_pull) > 0 {
			fmt.Println("Список сообщений для удаления")
			for key := range result_pull {
				fmt.Println(result_pull[key])
			}
		} else {
			fmt.Println("Нет системных сообщений для удаления")
		}
	} else {
		if len(result_pull) > 0 {
			websocketConnection, err := client.Ws(settings.TeamUid, nil)
			if err != nil {
				panic(err)
			}
			for key := range result_pull {
				websocketConnection.DeleteMessage(key)
				time.Sleep(3 * time.Second)
				fmt.Printf("сообщение %s было удалено [%s]", key, result_pull[key])
				if _, err := websocketConnection.WaitForConfirm(); err != nil {
					panic(err)
				}
			}
			websocketConnection.Close()
		} else {
			fmt.Println("Нет системных сообщений для удаления")
		}
	}

}
