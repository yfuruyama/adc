adc - Application Default Credentials manager for GCP [![CircleCI](https://circleci.com/gh/yfuruyama/adc.svg?style=svg)](https://circleci.com/gh/yfuruyama/adc)
===

```
Usage: adc [--version] [--help] <command> [<args>]

Available commands are:
    add        Add service account credential
    cat        Cat the credential content
    current    Show the current active credential
    env        Display the commands to set up the credential environment for application
    exec       Execute the command with the specified credential
    list       Show available credentials
    token      Prints access token for the credential
```

## Install

```
go get -u github.com/yfuruyama/adc
```

## TODO

* Support credential from Metadata server
