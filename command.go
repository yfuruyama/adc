package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

const (
	statusSuccess = 0
	statusError   = 1
)

type Stream struct {
	outStream io.Writer
	errStream io.Writer
}

type CurrentCommand struct {
	Stream
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
	return "Show the current active credential"
}

func (c *CurrentCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s current`, cmd)
}

const (
	ListFormatTable = "table"
	ListFormatCsv   = "csv"
	ListFormatTsv   = "tsv"
)

type ListCommand struct {
	Stream
	Format string
}

func (c *ListCommand) Run(args []string) int {
	flags := flag.NewFlagSet("list", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
	}
	flags.StringVar(&c.Format, "format", ListFormatTable, "output format")
	if err := flags.Parse(args); err != nil {
		return statusError
	}

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

	header := []string{"Credential", "Active", "Project", "Type"}

	switch c.Format {
	case ListFormatTsv:
		fallthrough
	case ListFormatCsv:
		var separator string
		if c.Format == ListFormatTsv {
			separator = "\t"
		} else if c.Format == ListFormatCsv {
			separator = ","
		}
		fmt.Fprintf(c.outStream, strings.Join(header, separator)+"\n")
		for _, credential := range credentials {
			active := "false"
			if currentCredential != nil && credential.Name() == currentCredential.Name() {
				active = "true"
			}
			var projectId string
			if credential.Type == CredentialTypeServiceAccount {
				projectId = credential.ProjectId
			}
			fmt.Fprintf(c.outStream, "%s%s%s%s%s%s%s\n", credential.Name(), separator, active, separator, projectId, separator, credential.Type.Name())
		}
	default:
		table := tablewriter.NewWriter(c.outStream)
		table.SetAutoFormatHeaders(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeader(header)

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
	}

	return statusSuccess
}

func (c *ListCommand) Synopsis() string {
	return "Show available credentials"
}

func (c *ListCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s list [--format=<format>]

Available formats: table(default), csv, tsv`, cmd)
}

type CatCommand struct {
	Stream
}

func (c *CatCommand) Run(args []string) int {
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

func (c *CatCommand) Synopsis() string {
	return "Cat the credential content"
}

func (c *CatCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s cat <credential>`, cmd)
}

type AddCommand struct {
	Stream
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
	Stream
}

func (c *ExecCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintf(c.errStream, c.Help()+"\n")
		return statusError
	}
	credentialName := args[0]

	var child string
	var childArgs []string
	if args[1] == "--" {
		if len(args) < 3 {
			fmt.Fprintf(c.errStream, c.Help()+"\n")
			return statusError
		}
		child = args[2]
		childArgs = args[3:]
	} else {
		child = args[1]
		childArgs = args[2:]
	}

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
	return fmt.Sprintf(`Usage: %s exec <credential> [--] <command> [<args>...]`, cmd)
}

type EnvCommand struct {
	Stream
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

type TokenCommand struct {
	Stream
}

func (c *TokenCommand) Run(args []string) int {
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

	token, err := credential.GetAccessToken()
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to get token: %s\n", err)
		return statusError
	}

	fmt.Fprintf(c.outStream, token+"\n")

	return statusSuccess
}

func (c *TokenCommand) Synopsis() string {
	return "Prints access token for the credential"
}

func (c *TokenCommand) Help() string {
	cmd := os.Args[0]
	return fmt.Sprintf(`Usage: %s token <credential>`, cmd)
}
