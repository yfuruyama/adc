adc - Application Default Credentials manager for GCP
===

```
$ adc help
Usage: adc [--version] [--help] <command> [<args>]

Available commands are:
    add        Add service account credential
    current    Show current active credential
    env        Display the commands to set up the credential environment for application
    exec       Execute the command with the specified credential
    list       Show available credentials
    show       Show credential file content
```

## Install

```
go get -u github.com/yfuruyama/adc
```

## Example Usage

Get access token for `mycredential`

```
adc exec mycredential -- gcloud auth application-default print-access-token
```
