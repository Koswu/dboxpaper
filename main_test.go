package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestEnvVarCredentials(t *testing.T) {
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
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
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
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "settings.json"), []byte(`{"access_token":"dummy","token_type":"Bearer"}`), 0600); err != nil {
		t.Fatalf("Failed to write settings file: %v", err)
	}
	
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

func TestMissingEnvVars(t *testing.T) {
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

	testCases := []struct {
		name          string
		setClientID   bool
		setSecret     bool
		expectedError string
	}{
		{
			name:          "Missing both credentials",
			setClientID:   false,
			setSecret:     false,
			expectedError: "DROPBOX_CLIENT_ID environment variable is required",
		},
		{
			name:          "Missing client secret",
			setClientID:   true,
			setSecret:     false,
			expectedError: "DROPBOX_CLIENT_SECRET environment variable is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unset environment variables
			os.Unsetenv("DROPBOX_CLIENT_ID")
			os.Unsetenv("DROPBOX_CLIENT_SECRET")

			if tc.setClientID {
				os.Setenv("DROPBOX_CLIENT_ID", "test_id")
			}
			if tc.setSecret {
				os.Setenv("DROPBOX_CLIENT_SECRET", "test_secret")
			}

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
			
			err := initialize(ctx)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error containing %q, got %q", tc.expectedError, err.Error())
			}
		})
	}
}


