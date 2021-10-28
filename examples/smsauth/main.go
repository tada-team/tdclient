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

	prompt := promptui.Prompt{Label: "Enter phone"}
	phone, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	codeResp, err := client.AuthBySmsSendCode(phone)
	if err != nil {
		panic(err)
	}

	fmt.Println("SMS sent to:", codeResp.Phone)
	prompt = promptui.Prompt{Label: "Enter code from SMS"}
	code, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	tokenResp, err := client.AuthBySmsGetToken(phone, code)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your token:", tokenResp.Token)
}
