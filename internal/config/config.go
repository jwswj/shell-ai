package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// API configuration
	OpenAIAPIKey       string `json:"OPENAI_API_KEY"`
	OpenAIModel        string `json:"OPENAI_MODEL"`
	OpenAIMaxTokens    int    `json:"OPENAI_MAX_TOKENS"`
	OpenAIAPIBase      string `json:"OPENAI_API_BASE"`
	OpenAIOrganization string `json:"OPENAI_ORGANIZATION"`
	OpenAIProxy        string `json:"OPENAI_PROXY"`
	OpenAIAPIVersion   string `json:"OPENAI_API_VERSION"`

	// Groq configuration
	GroqAPIKey string `json:"GROQ_API_KEY"`
	GroqModel  string `json:"GROQ_MODEL"`

	// Application configuration
	APIProvider     string  `json:"SHAI_API_PROVIDER"`
	SuggestionCount int     `json:"SHAI_SUGGESTION_COUNT"`
	SkipConfirm     bool    `json:"SHAI_SKIP_CONFIRM"`
	SkipHistory     bool    `json:"SHAI_SKIP_HISTORY"`
	Temperature     float64 `json:"SHAI_TEMPERATURE"`
	Debug           bool    `json:"DEBUG"`
	ContextMode     bool    `json:"CTX"`
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Create a new config with default values
	cfg := &Config{
		OpenAIModel:      "gpt-3.5-turbo",
		SuggestionCount:  3,
		APIProvider:      "groq",
		GroqModel:        "llama-3.3-70b-versatile",
		Temperature:      0.05,
		OpenAIAPIVersion: "2023-05-15",
	}

	// Load from config file (overrides defaults)
	if err := loadFromConfigFile(cfg); err != nil {
		// Log the error but continue - config file is optional
		if os.IsNotExist(err) {
			// This is fine, config file is optional
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
		}
	}

	// Load from environment variables (overrides config file)
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %w", err)
	}

	return cfg, nil
}

// loadFromConfigFile loads configuration from a JSON file
func loadFromConfigFile(cfg *Config) error {
	configAppName := "shell-ai"
	var configPath string

	// Determine the path to the configuration file based on the platform
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("APPDATA"), configAppName, "config.json")
	} else {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", configAppName, "config.json")
	}

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Parse the JSON into a map first to handle type conversions
	var configMap map[string]string
	if err := json.Unmarshal(data, &configMap); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Apply the values from the map to the config struct
	if val, ok := configMap["OPENAI_API_KEY"]; ok {
		cfg.OpenAIAPIKey = val
	}
	if val, ok := configMap["OPENAI_MODEL"]; ok {
		cfg.OpenAIModel = val
	}
	if val, ok := configMap["OPENAI_MAX_TOKENS"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.OpenAIMaxTokens = i
		}
	}
	if val, ok := configMap["OPENAI_API_BASE"]; ok {
		cfg.OpenAIAPIBase = val
	}
	if val, ok := configMap["OPENAI_ORGANIZATION"]; ok {
		cfg.OpenAIOrganization = val
	}
	if val, ok := configMap["OPENAI_PROXY"]; ok {
		cfg.OpenAIProxy = val
	}
	if val, ok := configMap["OPENAI_API_VERSION"]; ok {
		cfg.OpenAIAPIVersion = val
	}
	if val, ok := configMap["GROQ_API_KEY"]; ok {
		cfg.GroqAPIKey = val
	}
	if val, ok := configMap["GROQ_MODEL"]; ok {
		cfg.GroqModel = val
	}
	if val, ok := configMap["SHAI_API_PROVIDER"]; ok {
		cfg.APIProvider = val
	}
	if val, ok := configMap["SHAI_SUGGESTION_COUNT"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.SuggestionCount = i
		}
	}
	if val, ok := configMap["SHAI_SKIP_CONFIRM"]; ok {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.SkipConfirm = b
		}
	}
	if val, ok := configMap["SHAI_SKIP_HISTORY"]; ok {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.SkipHistory = b
		}
	}
	if val, ok := configMap["SHAI_TEMPERATURE"]; ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.Temperature = f
		}
	}
	if val, ok := configMap["DEBUG"]; ok {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.Debug = b
		}
	}
	if val, ok := configMap["CTX"]; ok {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.ContextMode = b
		}
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) error {
	// Check each environment variable and override if set
	if val := os.Getenv("OPENAI_API_KEY"); val != "" {
		cfg.OpenAIAPIKey = val
	}
	if val := os.Getenv("OPENAI_MODEL"); val != "" {
		cfg.OpenAIModel = val
	}
	if val := os.Getenv("OPENAI_MAX_TOKENS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.OpenAIMaxTokens = i
		}
	}
	if val := os.Getenv("OPENAI_API_BASE"); val != "" {
		cfg.OpenAIAPIBase = val
	}
	if val := os.Getenv("OPENAI_ORGANIZATION"); val != "" {
		cfg.OpenAIOrganization = val
	}
	if val := os.Getenv("OPENAI_PROXY"); val != "" {
		cfg.OpenAIProxy = val
	}
	if val := os.Getenv("OPENAI_API_VERSION"); val != "" {
		cfg.OpenAIAPIVersion = val
	}
	if val := os.Getenv("GROQ_API_KEY"); val != "" {
		cfg.GroqAPIKey = val
	}
	if val := os.Getenv("GROQ_MODEL"); val != "" {
		cfg.GroqModel = val
	}
	if val := os.Getenv("SHAI_API_PROVIDER"); val != "" {
		cfg.APIProvider = val
	}
	if val := os.Getenv("SHAI_SUGGESTION_COUNT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			cfg.SuggestionCount = i
		}
	}
	if val := os.Getenv("SHAI_SKIP_CONFIRM"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.SkipConfirm = b
		}
	}
	if val := os.Getenv("SHAI_SKIP_HISTORY"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.SkipHistory = b
		}
	}
	if val := os.Getenv("SHAI_TEMPERATURE"); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.Temperature = f
		}
	}
	if val := os.Getenv("DEBUG"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.Debug = b
		}
	}
	if val := os.Getenv("CTX"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.ContextMode = b
		}
	}

	return nil
}

// DebugPrint prints debug information if debug mode is enabled
func (c *Config) DebugPrint(format string, args ...interface{}) {
	if c.Debug {
		fmt.Printf(format, args...)
	}
}
