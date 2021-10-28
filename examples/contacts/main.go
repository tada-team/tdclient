package main

import (
	"fmt"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
)

func main() {
	settings := examples.NewSettings()
	settings.RequireTeam()
	settings.RequireToken()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)

	contacts, err := client.Contacts(settings.TeamUid)
	if err != nil {
		panic(err)
	}

	// full fields list: https://github.com/tada-team/tdproto/blob/master/contact.go
	for _, contact := range contacts {
		fmt.Printf("%s\t%s\t%s\n", contact.TeamStatus, contact.Jid, contact.DisplayName)
	}
}
