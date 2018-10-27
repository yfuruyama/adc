package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

const (
	statusSuccess = 0
	statusError   = 1
)

type Command struct {
	outStream io.Writer
	errStream io.Writer
}

type CurrentCommand struct {
	Command
}

func (c *CurrentCommand) Run(args []string) int {
	credential, err := GetCurrentCredential()
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get current credential: %s\n", err)
		return statusError
	}

	if credential != nil {
		fmt.Fprintf(c.outStream, credential.Name()+"\n")
	}
	return statusSuccess
}

func (c *CurrentCommand) Synopsis() string {
	return "Show current active credential"
}

func (c *CurrentCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s current`, cmd)
}

type ListCommand struct {
	Command
}

func (c *ListCommand) Run(args []string) int {
	credentials, err := GetAllCredentials()
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get credentials: %s\n", err)
		return statusError
	}

	currentCredential, err := GetCurrentCredential()
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get current active credential: %s\n", err)
		return statusError
	}

	if len(credentials) == 0 {
		fmt.Fprintf(c.outStream, "No credentials found\n")
		return statusSuccess
	}

	table := tablewriter.NewWriter(c.outStream)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Credential", "Active", "Project", "Type"})

	for _, credential := range credentials {
		var active string
		if currentCredential != nil && credential.Name() == currentCredential.Name() {
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
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s list`, cmd)
}

type ShowCommand struct {
	Command
}

func (c *ShowCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
		return statusError
	}
	credentialName := args[0]

	credential, err := GetCredentialByPrefixName(credentialName)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get credential: %s\n", err)
		return statusError
	}
	if credential == nil {
		fmt.Fprintf(c.errStream, "Credential `%s` not found\n", credentialName)
		return statusError
	}

	file, err := os.Open(credential.filePath)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to read credential: %s\n", err)
		return statusError
	}

	if _, err := io.Copy(c.outStream, file); err != nil {
		fmt.Fprintf(c.errStream, "failed to read credential: %s\n", err)
		return statusError
	}

	return statusSuccess
}

func (c *ShowCommand) Synopsis() string {
	return "Show credential file content"
}

func (c *ShowCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s show <credential>`, cmd)
}

type AddCommand struct {
	Command
}

func (c *AddCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
		return statusError
	}

	filePath := args[0]
	credentialName := path.Base(filePath)

	// TODO: check valid credential
	src, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to read credential file: %s\n", err)
		return statusError
	}

	storePath, err := GetCredentialStorePath()
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to add credential file: %s\n", err)
		return statusError
	}
	destPath := path.Join(storePath, credentialName)
	dest, err := os.Create(destPath)

	if _, err := io.Copy(dest, src); err != nil {
		fmt.Fprintf(c.errStream, "failed to add credential file: %s\n", err)
		return statusError
	}

	fmt.Fprintf(c.outStream, "Added to credentials store: %s\n", destPath)
	return statusSuccess
}

func (c *AddCommand) Synopsis() string {
	return "Add service account credential"
}

func (c *AddCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s add <SERVICE_ACCOUNT_CREDENTIAL_KEY.json>`, cmd)
}

type ExecCommand struct {
	Command
}

func (c *ExecCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
		return statusError
	}
	credentialName := args[0]
	child := args[1]
	childArgs := args[2:]

	credential, err := GetCredentialByPrefixName(credentialName)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get credential: %s\n", err)
		return statusError
	}
	if credential == nil {
		fmt.Fprintf(c.errStream, "Credential `%s` not found\n", credentialName)
		return statusError
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
			fmt.Fprintf(c.errStream, "%s", err)
			return statusError
		}
	}

	return statusSuccess
}

func (c *ExecCommand) Synopsis() string {
	return "Execute the command with the specified credential"
}

func (c *ExecCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s exec <credential> <command> <args>...`, cmd)
}

type EnvCommand struct {
	Command
}

func (c *EnvCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
		return statusError
	}
	credentialName := args[0]

	credential, err := GetCredentialByPrefixName(credentialName)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get credential: %s\n", err)
		return statusError
	}
	if credential == nil {
		fmt.Fprintf(c.errStream, "Credential `%s` not found\n", credentialName)
		return statusError
	}

	cmd := os.Args[0]
	fmt.Fprintf(c.outStream, `export GOOGLE_APPLICATION_CREDENTIALS="%s"
# Run this command to configure your shell:
# eval "$(%s env %s)"
`, credential.filePath, cmd, credential.Name())

	return statusSuccess
}

func (c *EnvCommand) Synopsis() string {
	return "Display the commands to set up the credential environment for application"
}

func (c *EnvCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s env <credential>`, cmd)
}
