package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/tada-team/tdclient"
	"github.com/tada-team/tdclient/examples"
)

func main() {
	settings := examples.NewSettings()
	settings.Parse()

	client, err := tdclient.NewSession(settings.Server)
	if err != nil {
		panic(err)
	}

	client.SetVerbose(settings.Verbose)

	prompt := promptui.Prompt{Label: "Enter login"}
	login, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	prompt = promptui.Prompt{Label: "Enter password", Mask: '*'}
	password, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	tokenResp, err := client.AuthByPasswordGetToken(login, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your token:", tokenResp.Token)




}
