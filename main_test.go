package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestEnvVarClientID(t *testing.T) {
	// Save original env vars
	origClientID := os.Getenv("DROPBOX_CLIENT_ID")
	origClientSecret := os.Getenv("DROPBOX_CLIENT_SECRET")
	defer func() {
		if origClientID != "" {
			os.Setenv("DROPBOX_CLIENT_ID", origClientID)
		} else {
			os.Unsetenv("DROPBOX_CLIENT_ID")
		}
		if origClientSecret != "" {
			os.Setenv("DROPBOX_CLIENT_SECRET", origClientSecret)
		} else {
			os.Unsetenv("DROPBOX_CLIENT_SECRET")
		}
	}()

	// Test with custom environment variables
	testClientID := "custom_client_id"
	testClientSecret := "custom_client_secret"
	os.Setenv("DROPBOX_CLIENT_ID", testClientID)
	os.Setenv("DROPBOX_CLIENT_SECRET", testClientSecret)

	// Create temporary directory for settings
	tmpDir := filepath.Join(os.TempDir(), "dboxpaper_test")
	os.MkdirAll(tmpDir, 0700)
	defer os.RemoveAll(tmpDir)

	// Create a dummy settings file to prevent OAuth flow
	settingsFile := filepath.Join(tmpDir, "settings.json")
	os.WriteFile(settingsFile, []byte(`{"access_token":"dummy","token_type":"Bearer"}`), 0600)

	// Override HOME to use temp directory
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	testApp := cli.NewApp()
	testApp.Before = initialize
	if testApp.Metadata == nil {
		testApp.Metadata = make(map[string]interface{})
	}

	// Override global app for this test
	origApp := app
	app = testApp
	defer func() { app = origApp }()

	ctx := cli.NewContext(testApp, nil, nil)
	
	// Create settings directory structure
	configDir := filepath.Join(tmpDir, ".config", "dboxpaper")
	os.MkdirAll(configDir, 0700)
	os.WriteFile(filepath.Join(configDir, "settings.json"), []byte(`{"access_token":"dummy","token_type":"Bearer"}`), 0600)
	
	err := initialize(ctx)
	if err != nil {
		t.Logf("Initialize returned error (expected for test): %v", err)
	}

	dboxpaper := testApp.Metadata["dboxpaper"].(*DboxPaper)
	if dboxpaper == nil {
		t.Fatal("dboxpaper not initialized in metadata")
	}

	if dboxpaper.config.ClientID != testClientID {
		t.Errorf("Expected ClientID %s, got %s", testClientID, dboxpaper.config.ClientID)
	}

	if dboxpaper.config.ClientSecret != testClientSecret {
		t.Errorf("Expected ClientSecret %s, got %s", testClientSecret, dboxpaper.config.ClientSecret)
	}
}

func TestDefaultClientID(t *testing.T) {
	// Save original env vars
	origClientID := os.Getenv("DROPBOX_CLIENT_ID")
	origClientSecret := os.Getenv("DROPBOX_CLIENT_SECRET")
	defer func() {
		if origClientID != "" {
			os.Setenv("DROPBOX_CLIENT_ID", origClientID)
		} else {
			os.Unsetenv("DROPBOX_CLIENT_ID")
		}
		if origClientSecret != "" {
			os.Setenv("DROPBOX_CLIENT_SECRET", origClientSecret)
		} else {
			os.Unsetenv("DROPBOX_CLIENT_SECRET")
		}
	}()

	// Unset environment variables to test defaults
	os.Unsetenv("DROPBOX_CLIENT_ID")
	os.Unsetenv("DROPBOX_CLIENT_SECRET")

	// Create temporary directory for settings
	tmpDir := filepath.Join(os.TempDir(), "dboxpaper_test_default")
	os.MkdirAll(tmpDir, 0700)
	defer os.RemoveAll(tmpDir)

	// Override HOME to use temp directory
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	testApp := cli.NewApp()
	testApp.Before = initialize
	if testApp.Metadata == nil {
		testApp.Metadata = make(map[string]interface{})
	}

	// Override global app for this test
	origApp := app
	app = testApp
	defer func() { app = origApp }()

	ctx := cli.NewContext(testApp, nil, nil)
	
	// Create settings directory structure
	configDir := filepath.Join(tmpDir, ".config", "dboxpaper")
	os.MkdirAll(configDir, 0700)
	os.WriteFile(filepath.Join(configDir, "settings.json"), []byte(`{"access_token":"dummy","token_type":"Bearer"}`), 0600)
	
	err := initialize(ctx)
	if err != nil {
		t.Logf("Initialize returned error (expected for test): %v", err)
	}

	dboxpaper := testApp.Metadata["dboxpaper"].(*DboxPaper)
	if dboxpaper == nil {
		t.Fatal("dboxpaper not initialized in metadata")
	}

	defaultClientID := "nrb8y95k7yoeor6"
	defaultClientSecret := "fhme6tzwkzw5og8"

	if dboxpaper.config.ClientID != defaultClientID {
		t.Errorf("Expected default ClientID %s, got %s", defaultClientID, dboxpaper.config.ClientID)
	}

	if dboxpaper.config.ClientSecret != defaultClientSecret {
		t.Errorf("Expected default ClientSecret %s, got %s", defaultClientSecret, dboxpaper.config.ClientSecret)
	}
}

