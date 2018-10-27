package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func initialize() {
	if err := InitCredentialsStore(); err != nil {
		fmt.Printf("failed to initialize credentials store: %s\n", err)
		os.Exit(-1)
	}
}

func main() {
	initialize()

	c := cli.NewCLI("adc", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &ListCommand{}, nil
		},
		"add": func() (cli.Command, error) {
			return &AddCommand{}, nil
		},
		"current": func() (cli.Command, error) {
			return &CurrentCommand{}, nil
		},
		"use": func() (cli.Command, error) {
			return &UseCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
