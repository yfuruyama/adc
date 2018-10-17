package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type CredentialType string

const (
	CredentialTypeUserAccount    CredentialType = "authorized_user"
	CredentialTypeServiceAccount CredentialType = "service_account"
)

func (t CredentialType) Name() string {
	switch t {
	case CredentialTypeUserAccount:
		return "User Account"
	case CredentialTypeServiceAccount:
		return "Service Account"
	}
	return ""
}

type Credential struct {
	// for both keys
	Type     CredentialType `json:"type"`
	ClientId string         `json:"client_id"`

	// for only user account key
	RefreshToken string `json:"refresh_token"`
	ClientSecret string `json:"client_secret"`

	// for only service account key
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`

	// additional information by this tool
	fileName string
}

func GetDefaultCredential() (*Credential, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	filepath := path.Join(currentUser.HomeDir, ".config", "gcloud", "application_default_credentials.json")
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var credential Credential
	if err := json.Unmarshal(b, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

func GetCredentialByName(name string) (*Credential, error) {
	credentials, err := GetAllCredentials()
	if err != nil {
		return nil, err
	}
	for _, credential := range credentials {
		if credential.Name() == name {
			return credential, nil
		}
	}
	return nil, nil
}

func GetAllCredentials() ([]*Credential, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	dirpath := path.Join(currentUser.HomeDir, ".config", "adc")
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}

	fileinfoList, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	credentials := make([]*Credential, 0)
	for _, fileinfo := range fileinfoList {
		fileName := fileinfo.Name()
		filepath := path.Join(currentUser.HomeDir, ".config", "adc", fileName)
		credential, err := readCredentialFile(filepath)
		if err != nil {
			return nil, err
		}
		credential.fileName = fileName
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func (c *Credential) Name() string {
	switch c.Type {
	case CredentialTypeUserAccount:
		return c.fileName
	case CredentialTypeServiceAccount:
		return fmt.Sprintf("%s-%s", c.ProjectId, c.PrivateKeyId[0:12])
	}
	return ""
}

func readCredentialFile(filename string) (*Credential, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var credential Credential
	if err := json.Unmarshal(b, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}
