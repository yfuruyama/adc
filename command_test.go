package main

import (
	"bytes"
	"os"
	"testing"
)

func TestCurrentCommand(t *testing.T) {
	t.Run("env variable set", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "fixtures/service_account_credential_001.json")
		defer func() {
			os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		}()

		cmd := &CurrentCommand{Stream{outStream, errStream}}
		cmd.Run([]string{})

		got := outStream.String()
		expected := "service_account_credential_001-012345\n"
		if got != expected {
			t.Errorf("expected = %s, but got = %s", expected, got)
		}
	})

	t.Run("env variable not set, but gcloud application default credential exists", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		gcloudDefaultCredentialPath = "fixtures/user_credential_001.json"

		cmd := &CurrentCommand{Stream{outStream, errStream}}
		cmd.Run([]string{})

		got := outStream.String()
		expected := "authorized_user\n"
		if got != expected {
			t.Errorf("expected = %s, but got = %s", expected, got)
		}
	})

	t.Run("no active credential", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		gcloudDefaultCredentialPath = "fixtures/not_exists.json"

		cmd := &CurrentCommand{Stream{outStream, errStream}}
		cmd.Run([]string{})

		got := outStream.String()
		expected := ""
		if got != expected {
			t.Errorf("expected = %s, but got = %s", expected, got)
		}
	})
}
