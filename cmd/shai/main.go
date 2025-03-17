package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jwswj/shell-ai/internal/config"
	"github.com/jwswj/shell-ai/internal/llm"
	"github.com/jwswj/shell-ai/internal/suggestions"
)

var CLI struct {
	Debug bool     `help:"Enable debug mode" env:"DEBUG"`
	Ctx   bool     `help:"Set context mode to True" env:"CTX"`
	Prompt []string `arg:"" optional:"" help:"The prompt to generate shell commands for"`
}

func main() {
	ctx := kong.Parse(&CLI)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Set debug mode from CLI flag
	if CLI.Debug {
		cfg.Debug = true
	}

	// Set context mode from CLI flag
	if CLI.Ctx {
		cfg.ContextMode = true
	}

	// Check if API keys are set
	if cfg.OpenAIAPIKey == "" && cfg.GroqAPIKey == "" {
		fmt.Println("DEBUG: OpenAI API Key:", cfg.OpenAIAPIKey)
		fmt.Println("DEBUG: Groq API Key:", cfg.GroqAPIKey)
		fmt.Println("DEBUG: API Provider:", cfg.APIProvider)
		fmt.Println("Please set either the OPENAI_API_KEY or GROQ_API_KEY environment variable.")
		fmt.Println("You can also create `config.json` under `~/.config/shell-ai/` to set the API key, see README.md for more information.")
		os.Exit(1)
	}

	// Create LLM client based on configuration
	client, err := llm.NewClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating LLM client: %v\n", err)
		os.Exit(1)
	}

	// Run the command
	switch ctx.Command() {
	default:
		if len(CLI.Prompt) == 0 {
			fmt.Println("Describe what you want to do as a single sentence. `shai <sentence>`")
			os.Exit(0)
		}

		// Run the suggestions engine
		err = suggestions.Run(client, cfg, CLI.Prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running suggestions: %v\n", err)
			os.Exit(1)
		}
	}
} 