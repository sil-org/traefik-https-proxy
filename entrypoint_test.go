package main

import (
	"os"
	"regexp"
	"testing"
)

func TestUpdateConfigContent(t *testing.T) {
	original := `
example TEST
another FIELD	
`
	expected := `
example val
another green	
`
	replacements := []Replacement{
		{
			Key:   "TEST",
			Value: "val",
		},
		{
			Key:   "FIELD",
			Value: "green",
		},
	}
	results := UpdateConfigContent([]byte(original), replacements)

	if string(results) != expected {
		t.Fatal("Results to not match expected. Results:", results)
	}
}

func TestBuildReplacementsFromEnv(t *testing.T) {
	// Test failure for required env var
	_, err := BuildReplacementsFromEnv()
	if err == nil {
		t.Fatal("BuildReplacementsFromEnv should have failed because no env vars have been set")
	}

	setRequiredEnvVars()

	replacements, err := BuildReplacementsFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	replacementsCount := len(replacements)
	if replacementsCount != 6 {
		t.Fatal("Replacements did not have enough entries, only found", replacementsCount, "but expected 6")
	}
}

func TestReadUpdateWrite(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Unable to get current working directory for TestReadUpdateWrite")
	}
	readFile := dir + "/traefik.toml"
	writeFile := dir + "/traefik_test.toml"

	// If writeFile already exists, delete it
	if _, err := os.Stat(writeFile); err == nil {
		os.Remove(writeFile)
	}

	configToml, err := ReadTraefikToml(readFile)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure placeholders for env vars exist
	envVars := GetEnvVarModels()
	for _, envvar := range envVars {
		search := regexp.MustCompile(envvar.Name)
		found := search.Find(configToml)
		if found == nil {
			t.Fatal("Did not find key in configToml template for env var", envvar.Name)
		}
	}

	// Update config with required env var values
	setRequiredEnvVars()
	replacements, err := BuildReplacementsFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	configToml = UpdateConfigContent(configToml, replacements)

	// Make sure placeholders for required env vars no longer exist in configToml
	for _, envvar := range envVars {
		if envvar.Required {
			search := regexp.MustCompile(envvar.Name)
			found := search.Find(configToml)
			if found != nil {
				t.Fatal("Uh oh, placeholder for required env var still present after update for env var:", envvar.Name)
			}
		}
	}

	// Write out test file for manual reivew
	err = WriteTraefikToml(writeFile, configToml)
	if err != nil {
		t.Fatal(err)
	}
}

func setRequiredEnvVars() {
	os.Setenv("LETS_ENCRYPT_EMAIL", "test@testing.com")
	os.Setenv("LETS_ENCRYPT_CA", "staging")
	os.Setenv("TLD", "testing.com")
	os.Setenv("SANS", "test.testing.com,another.testing.com")
	os.Setenv("BACKEND1_URL", "http://app:80")
	os.Setenv("FRONTEND1_DOMAIN", "test.testing.com")
}
