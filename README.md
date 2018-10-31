adc - Application Default Credentials manager for GCP [![CircleCI](https://circleci.com/gh/yfuruyama/adc.svg?style=svg)](https://circleci.com/gh/yfuruyama/adc)
===

```
Usage: adc [--version] [--help] <command> [<args>]

Available commands are:
    active    Print which credential is active
    add       Add service account credential
    cat       Cat credential content
    env       Display commands to set up the credential environment for application
    exec      Execute command with the specified credential
    ls        Show available credentials
    token     Prints access token for the credential
```

## Install

```
go get -u github.com/yfuruyama/adc
```

## TODO

* Support credential from Metadata server
