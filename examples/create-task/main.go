package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tada-team/tdproto/tdapi"

	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
	"github.com/tada-team/tdproto"
)

func main() {
	assignee := flag.String("assignee", "", "assignee jid")
	description := flag.String("description", "", "task text")
	public := flag.Bool("public", false, "public")

	settings := examples.NewSettings()
	settings.RequireTeam()
	settings.RequireToken()
	settings.Parse()

	if *assignee == "" {
		fmt.Println("-assignee required")
		os.Exit(0)
	}

	if *description == "" {
		fmt.Println("-description required")
		os.Exit(0)
	}

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetToken(settings.Token)
	client.SetVerbose(settings.Verbose)

	recipient := *tdproto.NewJID(*assignee)
	if !recipient.Valid() {
		panic("invalid assignee jid")
	}

	chat, err := client.CreateTask(settings.TeamUid, tdapi.Task{
		Description: *description,
		Assignee:    recipient,
		Public:      *public,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("task created: %s\n", chat.Jid)
}
