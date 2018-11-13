adc - Application Default Credentials manager for GCP [![CircleCI](https://circleci.com/gh/yfuruyama/adc.svg?style=svg)](https://circleci.com/gh/yfuruyama/adc)
===

adc is a tool for managing GCP credentials such as service account keys or user credentials which are used as [Application Default Credentials](https://cloud.google.com/docs/authentication/production) (ADC) from your application.

With this tool, you will be free from the typical credential management problem: There are a lot of service account keys in my Downloads folder.

## Usage

![gif](https://github.com/yfuruyama/adc/blob/master/screencast.gif)

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
After adding the credential, you can delete the original one.

```sh
$ adc add ~/Downloads/my-service-account-key-e50710fb4883.json
Added credential `e50710fb4883`

# remove the original one (optional)
$ rm ~/Downloads/my-service-account-key-e50710fb4883.json
```

### adc exec

`adc exec <CREDENTIAL>` executes an arbitrary command with specified credential.  
You can specify `<CREDENTIAL>` with first several characters of the credential.

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

`adc env` displays commands to active the credential for current shell.

```sh
# set environment variable
$ eval "$(adc env 27f)"

# following commands are executed with the credential `27f98a8b3b11`
$ terraform execute
```

### adc token

`adc token <CREDENTIAL>` prints access token for the credential.

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
