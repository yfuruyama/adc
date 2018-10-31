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
		"ls": func() (cli.Command, error) {
			return &ListCommand{
				Format: ListFormatStandard,
				Stream: Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"cat": func() (cli.Command, error) {
			return &CatCommand{
				Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"add": func() (cli.Command, error) {
			return &AddCommand{
				Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"active": func() (cli.Command, error) {
			return &ActiveCommand{
				Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"exec": func() (cli.Command, error) {
			return &ExecCommand{
				Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"env": func() (cli.Command, error) {
			return &EnvCommand{
				IsUnset: false,
				Stream:  Stream{os.Stdout, os.Stderr},
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &TokenCommand{
				Stream{os.Stdout, os.Stderr},
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	os.Exit(exitStatus)
}
