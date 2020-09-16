package main

import (
	"flag"
	"fmt"

	"github.com/tada-team/tdclient"
)

func main() {
	server := flag.String("server", "https://web.tada.team", "server address")
	team := flag.String("team", "", "team uid")
	token := flag.String("token", "", "bot token. Type \"/newbot <NAME>\" command in @TadaBot direct chat")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	if *token == "" {
		fmt.Println("-token required")
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

	contacts, err := client.Contacts(*team)
	if err != nil {
		panic(err)
	}

	for _, contact := range contacts {
		// full fields list: https://github.com/tada-team/tdproto/blob/master/contact.go
		fmt.Printf("%s\t%s\n", contact.Jid, contact.DisplayName)
	}
}
