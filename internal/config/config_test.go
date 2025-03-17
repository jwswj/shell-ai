package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "shell-ai-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, ".config", "shell-ai")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a test config file
	configFile := filepath.Join(configDir, "config.json")
	testConfig := map[string]string{
		"OPENAI_API_KEY":        "test-openai-key",
		"OPENAI_MODEL":          "test-openai-model",
		"GROQ_API_KEY":          "test-groq-key",
		"GROQ_MODEL":            "test-groq-model",
		"SHAI_API_PROVIDER":     "test-provider",
		"SHAI_TEMPERATURE":      "0.7",
		"SHAI_SUGGESTION_COUNT": "5",
	}

	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Save original HOME env var and set it to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Clear any existing environment variables that might interfere with the test
	clearEnvVars := []string{
		"OPENAI_API_KEY", "OPENAI_MODEL", "GROQ_API_KEY", "GROQ_MODEL",
		"SHAI_API_PROVIDER", "SHAI_TEMPERATURE", "SHAI_SUGGESTION_COUNT",
	}

	origEnvValues := make(map[string]string)
	for _, name := range clearEnvVars {
		origEnvValues[name] = os.Getenv(name)
		os.Unsetenv(name)
	}

	// Restore environment variables after the test
	defer func() {
		for name, value := range origEnvValues {
			if value != "" {
				os.Setenv(name, value)
			}
		}
	}()

	// Load the config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Print the loaded config for debugging
	t.Logf("Loaded config: %+v", cfg)

	// Test that all values were loaded correctly
	if cfg.OpenAIAPIKey != "test-openai-key" {
		t.Errorf("OpenAIAPIKey not loaded correctly, got: %s, want: %s", cfg.OpenAIAPIKey, "test-openai-key")
	}

	if cfg.OpenAIModel != "test-openai-model" {
		t.Errorf("OpenAIModel not loaded correctly, got: %s, want: %s", cfg.OpenAIModel, "test-openai-model")
	}

	if cfg.GroqAPIKey != "test-groq-key" {
		t.Errorf("GroqAPIKey not loaded correctly, got: %s, want: %s", cfg.GroqAPIKey, "test-groq-key")
	}

	if cfg.GroqModel != "test-groq-model" {
		t.Errorf("GroqModel not loaded correctly, got: %s, want: %s", cfg.GroqModel, "test-groq-model")
	}

	if cfg.APIProvider != "test-provider" {
		t.Errorf("APIProvider not loaded correctly, got: %s, want: %s", cfg.APIProvider, "test-provider")
	}

	if cfg.Temperature != 0.7 {
		t.Errorf("Temperature not loaded correctly, got: %f, want: %f", cfg.Temperature, 0.7)
	}

	if cfg.SuggestionCount != 5 {
		t.Errorf("SuggestionCount not loaded correctly, got: %d, want: %d", cfg.SuggestionCount, 5)
	}
}

func TestEnvironmentOverridesConfigFile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "shell-ai-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, ".config", "shell-ai")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a test config file
	configFile := filepath.Join(configDir, "config.json")
	testConfig := map[string]string{
		"OPENAI_API_KEY": "file-openai-key",
		"GROQ_API_KEY":   "file-groq-key",
	}

	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Save original env vars
	originalHome := os.Getenv("HOME")
	originalOpenAIKey := os.Getenv("OPENAI_API_KEY")
	originalGroqKey := os.Getenv("GROQ_API_KEY")

	// Set env vars for test
	os.Setenv("HOME", tempDir)
	os.Setenv("OPENAI_API_KEY", "env-openai-key")
	os.Setenv("GROQ_API_KEY", "env-groq-key")

	// Restore env vars after test
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("OPENAI_API_KEY", originalOpenAIKey)
		os.Setenv("GROQ_API_KEY", originalGroqKey)
	}()

	// Load the config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test that environment variables override config file
	if cfg.OpenAIAPIKey != "env-openai-key" {
		t.Errorf("OpenAIAPIKey not overridden correctly, got: %s, want: %s", cfg.OpenAIAPIKey, "env-openai-key")
	}

	if cfg.GroqAPIKey != "env-groq-key" {
		t.Errorf("GroqAPIKey not overridden correctly, got: %s, want: %s", cfg.GroqAPIKey, "env-groq-key")
	}
}

func TestDefaultValues(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "shell-ai-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp dir to avoid loading any existing config
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Clear any existing environment variables that might interfere with the test
	clearEnvVars := []string{
		"OPENAI_API_KEY", "OPENAI_MODEL", "GROQ_API_KEY", "GROQ_MODEL",
		"SHAI_API_PROVIDER", "SHAI_TEMPERATURE", "SHAI_SUGGESTION_COUNT",
	}

	origEnvValues := make(map[string]string)
	for _, name := range clearEnvVars {
		origEnvValues[name] = os.Getenv(name)
		os.Unsetenv(name)
	}

	// Restore environment variables after the test
	defer func() {
		for name, value := range origEnvValues {
			if value != "" {
				os.Setenv(name, value)
			}
		}
	}()

	// Load the config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test that default values are set correctly
	if cfg.OpenAIModel != "gpt-3.5-turbo" {
		t.Errorf("Default OpenAIModel not set correctly, got: %s, want: %s", cfg.OpenAIModel, "gpt-3.5-turbo")
	}

	if cfg.GroqModel != "llama-3.3-70b-versatile" {
		t.Errorf("Default GroqModel not set correctly, got: %s, want: %s", cfg.GroqModel, "llama-3.3-70b-versatile")
	}

	if cfg.APIProvider != "groq" {
		t.Errorf("Default APIProvider not set correctly, got: %s, want: %s", cfg.APIProvider, "groq")
	}

	if cfg.SuggestionCount != 3 {
		t.Errorf("Default SuggestionCount not set correctly, got: %d, want: %d", cfg.SuggestionCount, 3)
	}

	if cfg.Temperature != 0.05 {
		t.Errorf("Default Temperature not set correctly, got: %f, want: %f", cfg.Temperature, 0.05)
	}

	if cfg.OpenAIAPIVersion != "2023-05-15" {
		t.Errorf("Default OpenAIAPIVersion not set correctly, got: %s, want: %s", cfg.OpenAIAPIVersion, "2023-05-15")
	}
}
