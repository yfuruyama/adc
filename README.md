adc - Application Default Credentials manager for GCP [![CircleCI](https://circleci.com/gh/yfuruyama/adc.svg?style=svg)](https://circleci.com/gh/yfuruyama/adc)
===

adc is a tool for managing [Application Default Credentials](https://cloud.google.com/docs/authentication/production) (ADC), service account keys and user credentials for GCP applications.

## Background

While [Application Default Credentials](https://cloud.google.com/docs/authentication/production) (ADC) are great ways for authenticating your applications with GCP services, there are no standard rules for managing the credentials themself.
You can manage service account keys in your way, but someday you will find your Download folder is filled up with multiple service account keys.

With this tool, you can manage those credentials with simple commands and will be free from credential management problems.

## Demo

![gif](https://github.com/yfuruyama/adc/blob/master/screencast.gif)

## Usage

### adc ls

`adc ls` shows all registered credentials.

```sh
$ adc ls
NAME           ACTIVE   PROJECT           SERVICE_ACCOUNT       TYPE
user           -        -                 -                     User Account
27f98a8b3b11   *        my-project        storage-admin         Service Account
e3b36c383e05   -        my-project        bigquery-user         Service Account
e50710fb4883   -        another-project   cloud-kms-encryptor   Service Account
```

### adc add

`adc add <CREDENTIAL.json>` adds a service account credential to adc.
After adding the credential, you can delete the original one safely.

```sh
$ adc add ~/Downloads/my-service-account-key-e50710fb4883.json
Added credential `e50710fb4883`

# remove the original one (optional)
$ rm ~/Downloads/my-service-account-key-e50710fb4883.json
```

### adc exec

`adc exec <CREDENTIAL>` executes an arbitrary command with the specified credential.  
You can specify `<CREDENTIAL>` with just first several characters of the credential.

```sh
# execute `terraform plan` with the credential `27f98a8b3b11`
$ adc exec 27f -- terraform plan
```

### adc active

`adc active` shows current active credential.

```sh
$ adc active
27f98a8b3b11
```

### adc env

`adc env` displays commands to activate the credential for the current shell.

```sh
# set environment variable
$ eval "$(adc env 27f)"

# following commands are executed with the credential `27f98a8b3b11`
$ terraform execute
```

### adc token

`adc token <CREDENTIAL>` prints an access token for the credential.

```sh
$ adc token 27f
ya29.c.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### adc rm

`adc rm <CREDENTIAL>` removes a service account credential.

```sh
$ adc rm 27f
Removed credential `27f98a8b3b11`
```

## Install

You can get the latest binary from [Releases](https://github.com/yfuruyama/adc/releases).

Or if you have installed `go`, just execute `go get` from your console.

```
go get -u github.com/yfuruyama/adc
```

## TODO

* Support credential from Metadata server
