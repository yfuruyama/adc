package main

import "log"

type CurrentCommand struct{}

func (c *CurrentCommand) Run(args []string) int {
	credential, err := FromDefaultCredentialFile()
	if err != nil {
		log.Println(err)
	}
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
	return 0
}

func (c *ListCommand) Synopsis() string {
	return "Show available credentials"
}

func (c *ListCommand) Help() string {
	return "TODO"
}
