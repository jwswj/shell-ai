package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the application configuration
type Config struct {
	// API configuration
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAIMaxTokens   int
	OpenAIAPIBase     string
	OpenAIOrganization string
	OpenAIProxy       string
	OpenAIAPIVersion  string

	// Groq configuration
	GroqAPIKey string
	GroqModel  string

	// Application configuration
	APIProvider        string
	SuggestionCount    int
	SkipConfirm        bool
	SkipHistory        bool
	Temperature        float64
	Debug              bool
	ContextMode        bool
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Default configuration
	cfg := &Config{
		OpenAIModel:      "gpt-3.5-turbo",
		SuggestionCount:  3,
		APIProvider:      "groq",
		GroqModel:        "llama-3.3-70b-versatile",
		Temperature:      0.05,
		OpenAIAPIVersion: "2023-05-15",
	}

	// Load from config file
	configFile, err := loadConfigFile()
	if err == nil {
		// Merge config file values
		if configFile["OPENAI_API_KEY"] != "" {
			cfg.OpenAIAPIKey = configFile["OPENAI_API_KEY"]
		}
		if configFile["OPENAI_MODEL"] != "" {
			cfg.OpenAIModel = configFile["OPENAI_MODEL"]
		}
		if configFile["SHAI_SUGGESTION_COUNT"] != "" {
			fmt.Sscanf(configFile["SHAI_SUGGESTION_COUNT"], "%d", &cfg.SuggestionCount)
		}
		if configFile["SHAI_API_PROVIDER"] != "" {
			cfg.APIProvider = configFile["SHAI_API_PROVIDER"]
		}
		if configFile["GROQ_MODEL"] != "" {
			cfg.GroqModel = configFile["GROQ_MODEL"]
		}
		if configFile["GROQ_API_KEY"] != "" {
			cfg.GroqAPIKey = configFile["GROQ_API_KEY"]
		}
		if configFile["SHAI_TEMPERATURE"] != "" {
			fmt.Sscanf(configFile["SHAI_TEMPERATURE"], "%f", &cfg.Temperature)
		}
	}

	// Load from environment variables (overrides config file)
	if os.Getenv("OPENAI_API_KEY") != "" {
		cfg.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	}
	if os.Getenv("OPENAI_MODEL") != "" {
		cfg.OpenAIModel = os.Getenv("OPENAI_MODEL")
	}
	if os.Getenv("OPENAI_MAX_TOKENS") != "" {
		fmt.Sscanf(os.Getenv("OPENAI_MAX_TOKENS"), "%d", &cfg.OpenAIMaxTokens)
	}
	if os.Getenv("OPENAI_API_BASE") != "" {
		cfg.OpenAIAPIBase = os.Getenv("OPENAI_API_BASE")
	}
	if os.Getenv("OPENAI_ORGANIZATION") != "" {
		cfg.OpenAIOrganization = os.Getenv("OPENAI_ORGANIZATION")
	}
	if os.Getenv("OPENAI_PROXY") != "" {
		cfg.OpenAIProxy = os.Getenv("OPENAI_PROXY")
	}
	if os.Getenv("OPENAI_API_VERSION") != "" {
		cfg.OpenAIAPIVersion = os.Getenv("OPENAI_API_VERSION")
	}
	if os.Getenv("GROQ_API_KEY") != "" {
		cfg.GroqAPIKey = os.Getenv("GROQ_API_KEY")
	}
	if os.Getenv("GROQ_MODEL") != "" {
		cfg.GroqModel = os.Getenv("GROQ_MODEL")
	}
	if os.Getenv("SHAI_API_PROVIDER") != "" {
		cfg.APIProvider = os.Getenv("SHAI_API_PROVIDER")
	}
	if os.Getenv("SHAI_SUGGESTION_COUNT") != "" {
		fmt.Sscanf(os.Getenv("SHAI_SUGGESTION_COUNT"), "%d", &cfg.SuggestionCount)
	}
	if os.Getenv("SHAI_SKIP_CONFIRM") == "true" {
		cfg.SkipConfirm = true
	}
	if os.Getenv("SHAI_SKIP_HISTORY") == "true" {
		cfg.SkipHistory = true
	}
	if os.Getenv("SHAI_TEMPERATURE") != "" {
		fmt.Sscanf(os.Getenv("SHAI_TEMPERATURE"), "%f", &cfg.Temperature)
	}
	if os.Getenv("CTX") == "true" {
		cfg.ContextMode = true
	}
	if os.Getenv("DEBUG") == "true" {
		cfg.Debug = true
	}

	return cfg, nil
}

// loadConfigFile loads the configuration from a JSON file
func loadConfigFile() (map[string]string, error) {
	configAppName := "shell-ai"
	var configPath string

	// Determine the path to the configuration file based on the platform
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("APPDATA"), configAppName, "config.json")
	} else {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", configAppName, "config.json")
	}

	// Debug: Print the config path
	fmt.Printf("DEBUG: Looking for config file at: %s\n", configPath)

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("DEBUG: Error reading config file: %v\n", err)
		return nil, err
	}

	// Debug: Print the config file contents
	fmt.Printf("DEBUG: Config file contents: %s\n", string(data))

	// Parse the JSON
	var config map[string]string
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("DEBUG: Error parsing config file: %v\n", err)
		return nil, err
	}

	// Debug: Print the parsed config
	fmt.Printf("DEBUG: Parsed config: %v\n", config)

	return config, nil
}

// DebugPrint prints debug information if debug mode is enabled
func (c *Config) DebugPrint(format string, args ...interface{}) {
	if c.Debug {
		fmt.Printf(format, args...)
	}
} 