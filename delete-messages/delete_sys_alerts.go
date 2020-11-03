package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

var allMsg = make(map[string]string)
var resultPull = make(map[string]string)

func getLastMessageId(messages tdproto.ChatMessages) string {
	return messages.Messages[len(messages.Messages)-1].MessageId
}

func setMessageMap(messages tdproto.ChatMessages){
	for key := range messages.Messages {
		allMsg[messages.Messages[key].MessageId] = messages.Messages[key].Content.Text
	}
}

func filterMessages(){
	for key := range allMsg {
		var text = allMsg[key]
		s := strings.Split(text, ":")[0]

		if s == "Удалён участник" {
			resultPull[key] = text
		}
	}
}

func er(err error){
	if err != nil {
		panic(err)
	}
}

func main() {
	settings := examples.NewSettings()
	settings.RequireToken()
	settings.RequireTeam()
	settings.RequireChat()
	settings.RequireDryRun()
	settings.RequireDeep()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	er(err)

	client.SetToken(settings.Token)
	client.SetVerbose(settings.Verbose)

	chatUid := *tdproto.NewJID(settings.Chat)

	messages, err := client.GetMessages(settings.TeamUid, chatUid)
	er(err)

	setMessageMap(messages)
	var lastMsgId = getLastMessageId(messages)

	for i := 1; i < settings.Deep; i++ {
		messagesOld, err := client.GetMessagesOldMsg(settings.TeamUid, chatUid, lastMsgId)
		er(err)

		lastMsgId = getLastMessageId(messagesOld)
		setMessageMap(messagesOld)
	}

	filterMessages()

	if len(resultPull) > 0 {
		websocketConnection, err := client.Ws(settings.TeamUid, nil)
		er(err)

		for key := range resultPull {
			if settings.DryRun {
				fmt.Printf("сообщение %s будет удалено (dryrun) [%s]", key, resultPull[key])
			}else{
				websocketConnection.DeleteMessage(key)
				time.Sleep(3 * time.Second)

				fmt.Printf("сообщение %s удалено [%s]", key, resultPull[key])

				_, err := websocketConnection.WaitForConfirm()
				er(err)
			}

		}
		websocketConnection.Close()
	}else{
		fmt.Println("Нет системных сообщений для удаления")
	}
}
