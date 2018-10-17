package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("adc", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &ListCommand{}, nil
		},
		"current": func() (cli.Command, error) {
			return &CurrentCommand{}, nil
		},
		"use": func() (cli.Command, error) {
			return &UseCommand{}, nil
		},
		"login": func() (cli.Command, error) {
			return &LoginCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
