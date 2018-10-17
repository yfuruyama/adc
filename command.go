package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

type CurrentCommand struct{}

func (c *CurrentCommand) Run(args []string) int {
	credential, err := GetDefaultCredential()
	if err != nil {
		log.Println(err)
		return -1
	}
	// if user account, get user name and print it
	log.Println(credential.ClientId)
	return 0
}

func (c *CurrentCommand) Synopsis() string {
	return "Show current credential"
}

func (c *CurrentCommand) Help() string {
	return "TODO"
}

type ListCommand struct{}

func (c *ListCommand) Run(args []string) int {
	credentials, err := GetAllCredentials()
	if err != nil {
		log.Println(err)
		return -1
	}

	fmt.Println("No  Credential                      Type")
	fmt.Println("-------------------------------------------------------")
	for i, credential := range credentials {
		fmt.Printf("%d   ", i+1)
		fmt.Printf("%s   ", credential.Name())
		fmt.Printf("%s", credential.Type.Name())
		fmt.Println()
	}

	// No  Active  Credential                        Type
	// -------------------------------------------------------------
	// 1           addsict@gmail.com                 User Account
	// 2           yfuruyama@gmail.com               User Account
	// 3   *       yfuruyama-sandbox-98246f6f9623    Service Account
	// 4           yfuruyama-sandbox-123456789012    Service Account

	return 0
}

func (c *ListCommand) Synopsis() string {
	return "Show available credentials"
}

func (c *ListCommand) Help() string {
	return "TODO"
}

type LoginCommand struct{}

func (c *LoginCommand) Run(args []string) int {
	// remember current active credential
	// call `gcloud auth application-default login`
	// then call `gcloud auth application-default print-access-token`
	// then request to https://www.googleapis.com/oauth2/v3/tokeninfo?access_token={TOKEN}
	// then response.email to filename and copy
	// then copy ~/.config/gcloud/application_default_credentials.json to ~/.config/adc/{filename}.json
	// in the last, recover active credential
	return 0
}

func (c *LoginCommand) Synopsis() string {
	return "Alias for `gcloud auth application-default login`"
}

func (c *LoginCommand) Help() string {
	return "TODO"
}

type UseCommand struct{}

func (c *UseCommand) Run(args []string) int {
	if len(args) == 0 {
		log.Println("invalid usage")
		return -1
	}
	credentialName := args[0]

	credential, err := GetCredentialByName(credentialName)
	if err != nil {
		return -1
	}
	if credential == nil {
		fmt.Printf("Credential `%s` not found\n", credentialName)
		return -1
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	adcpath := path.Join(currentUser.HomeDir, ".config", "gcloud", "application_default_credentials.json")

	// delete old one
	if err := os.Remove(adcpath); err != nil {
		fmt.Println(err)
	}

	credpath := path.Join(currentUser.HomeDir, ".config", "adc", credential.fileName)

	if err := os.Symlink(credpath, adcpath); err != nil {
		fmt.Printf("Activate failed: %s\n", err)
		return -1
	}

	fmt.Printf("Credential `%s` activated\n", credentialName)

	return 0
}

func (c *UseCommand) Synopsis() string {
	return "Set a credential to the default credential"
}

func (c *UseCommand) Help() string {
	return "TODO"
}
