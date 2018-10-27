package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

const (
	statusSuccess = 0
	statusError   = -1

	envKey = "GOOGLE_APPLICATION_CREDENTIALS"
)

type CurrentCommand struct{}

func (c *CurrentCommand) Run(args []string) int {
	credential, err := GetCurrentCredential()
	if err != nil {
		log.Println(err)
		return statusError
	}

	fmt.Println(credential.Name())
	return statusSuccess
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

	currentCredential, err := GetCurrentCredential()
	if err != nil {
		log.Println(err)
		return -1
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Credential", "Active", "Project", "Type"})

	for _, credential := range credentials {
		var active string
		if credential.Name() == currentCredential.Name() {
			active = "(*)"
		}
		var projectId string
		if credential.Type == CredentialTypeServiceAccount {
			projectId = credential.ProjectId
		} else {
			projectId = "-"
		}
		table.Append([]string{credential.Name(), active, projectId, credential.Type.Name()})
	}

	table.Render()

	return statusSuccess
}

func (c *ListCommand) Synopsis() string {
	return "Show available credentials"
}

func (c *ListCommand) Help() string {
	return "TODO"
}

type AddCommand struct{}

func (c *AddCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("file not specified")
		return statusError
	}

	filePath := args[0]
	credentialName := path.Base(filePath)

	// TODO: check valid credential
	src, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return statusError
	}

	storePath, err := GetCredentialStorePath()
	if err != nil {
		fmt.Println(err)
		return statusError
	}
	destPath := path.Join(storePath, credentialName)
	dest, err := os.Create(destPath)

	if _, err := io.Copy(dest, src); err != nil {
		fmt.Println(err)
		return statusError
	}

	fmt.Println("Added to credentials store")
	return statusSuccess
}

func (c *AddCommand) Synopsis() string {
	return "Add service account credential"
}

func (c *AddCommand) Help() string {
	return "TODO"
}

type ExecCommand struct{}

func (c *ExecCommand) Run(args []string) int {
	if len(args) < 2 {
		log.Println("invalid usage")
		return -1
	}
	credentialName := args[0]
	child := args[1]
	childArgs := args[2:]

	credential, err := GetCredentialByName(credentialName)
	if err != nil {
		return statusError
	}
	if credential == nil {
		fmt.Printf("Credential `%s` not found\n", credentialName)
		return -1
	}

	env := os.Environ()
	env = append(env, "GOOGLE_APPLICATION_CREDENTIALS="+credential.filePath)

	cmd := exec.Command(child, childArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			fmt.Println(err)
			return statusError
		}
	}

	return statusSuccess
}

func (c *ExecCommand) Synopsis() string {
	return "TODO"
}

func (c *ExecCommand) Help() string {
	return "TODO"
}

type EnvCommand struct{}

func (c *EnvCommand) Run(args []string) int {
	if len(args) < 1 {
		log.Println("invalid usage")
		return statusError
	}
	credentialName := args[0]

	credential, err := GetCredentialByName(credentialName)
	if err != nil {
		return statusError
	}

	fmt.Printf(`export GOOGLE_APPLICATION_CREDENTIALS="%s"
# Run this command to configure your shell:
# eval "$(adc env %s)"
`, credential.filePath, credential.filePath)

	return statusSuccess
}

func (c *EnvCommand) Synopsis() string {
	return "Display the commands to set up the credentials environment for application"
}

func (c *EnvCommand) Help() string {
	return "TODO"
}
