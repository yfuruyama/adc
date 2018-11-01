adc - Application Default Credentials manager for GCP [![CircleCI](https://circleci.com/gh/yfuruyama/adc.svg?style=svg)](https://circleci.com/gh/yfuruyama/adc)
===

adc is a tool for managing GCP credentials such as service account keys or user credentials which are used as [Application Default Credentials](https://cloud.google.com/docs/authentication/production) (ADC) from your application.

With this tool, you will be free from the typical credential management problem: *There are always unknown service account keys in my Downloads folder.*

## Usage

![gif](https://github.com/yfuruyama/adc/blob/master/screencast.gif)

```
Usage: adc [--version] [--help] <command> [<args>]

Available commands are:
    active    Print which credential is active
    add       Add service account credential
    cat       Cat credential content
    env       Display commands to set up the credential environment for application
    exec      Execute command with the specified credential
    ls        Show all credentials
    rm        Remove service account credential
    token     Prints access token for the credential
```

## Install

```
go get -u github.com/yfuruyama/adc
```

## TODO

* Support credential from Metadata server
