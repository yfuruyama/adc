package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func initialize() {
	if err := InitCredentialsStore(); err != nil {
		fmt.Printf("failed to initialize credentials store: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	initialize()

	c := cli.NewCLI("adc", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &ListCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
		"show": func() (cli.Command, error) {
			return &ShowCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
		"add": func() (cli.Command, error) {
			return &AddCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
		"current": func() (cli.Command, error) {
			return &CurrentCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
		"exec": func() (cli.Command, error) {
			return &ExecCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
		"env": func() (cli.Command, error) {
			return &EnvCommand{
				Command{os.Stdout, os.Stderr},
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	os.Exit(exitStatus)
}
