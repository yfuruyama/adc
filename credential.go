package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"sort"
	"strings"

	"golang.org/x/oauth2/google"
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

	// for internal use
	filePath string
	json     []byte
}

func GetDefaultCredential() (*Credential, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(currentUser.HomeDir, ".config", "gcloud", "application_default_credentials.json")
	if _, err := os.Stat(filePath); err != nil {
		// application_default_credentials.json not found
		return nil, nil
	}

	return readCredentialFile(filePath)
}

func GetCredentialByPrefixName(name string) (*Credential, error) {
	credentials, err := GetAllCredentials()
	if err != nil {
		return nil, err
	}

	candidates := make([]*Credential, 0)
	for _, credential := range credentials {
		if strings.HasPrefix(credential.Name(), name) {
			candidates = append(candidates, credential)
		}
	}

	if len(candidates) == 1 {
		return candidates[0], nil
	} else if len(candidates) >= 2 {
		return nil, fmt.Errorf("multiple credentials found, `%s` is ambiguous.", name)
	}

	return nil, nil
}

func GetCredentialByPath(path string) (*Credential, error) {
	return readCredentialFile(path)
}

func GetAllCredentials() ([]*Credential, error) {
	storePath, err := GetCredentialStorePath()
	if err != nil {
		return nil, err
	}

	dir, err := os.Open(storePath)
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
		filePath := path.Join(storePath, fileName)
		credential, err := readCredentialFile(filePath)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}

	// add user account credential
	defaultCredential, err := GetDefaultCredential()
	if err != nil {
		return nil, err
	}
	if defaultCredential != nil {
		credentials = append(credentials, defaultCredential)
	}

	// sort alphabetically
	sort.Slice(credentials, func(i, j int) bool {
		return credentials[i].Name() < credentials[j].Name()
	})

	return credentials, nil
}

func GetCurrentCredential() (*Credential, error) {
	if envVar := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); envVar != "" {
		if _, err := os.Stat(envVar); err == nil {
			return GetCredentialByPath(envVar)
		}
	}
	return GetDefaultCredential()
}

func (c *Credential) Name() string {
	switch c.Type {
	case CredentialTypeUserAccount:
		return "authorized_user"
	case CredentialTypeServiceAccount:
		parts := strings.Split(c.ClientEmail, "@")
		serviceAccountId := parts[0]
		return fmt.Sprintf("%s-%s", serviceAccountId, c.PrivateKeyId[0:6])
	}
	return ""
}

func (c *Credential) GetAccessToken() (string, error) {
	ctx := context.Background()
	const scope = "https://www.googleapis.com/auth/cloud-platform"
	googleCred, nil := google.CredentialsFromJSON(ctx, c.json, scope)

	token, err := googleCred.TokenSource.Token()
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func InitCredentialsStore() error {
	path, err := GetCredentialStorePath()
	if err != nil {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

func GetCredentialStorePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// TODO: customize by user
	return path.Join(currentUser.HomeDir, ".config", "adc", "credentials"), nil
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

	credential.filePath = filename
	credential.json = b

	return &credential, nil
}
