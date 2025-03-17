package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/jwswj/shell-ai/internal/config"
)

// Client represents an LLM client
type Client struct {
	config *config.Config
	client *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewClient creates a new LLM client
func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// GenerateCompletion generates a completion from the LLM
func (c *Client) GenerateCompletion(systemPrompt, userPrompt string) (string, error) {
	var apiURL string
	var apiKey string
	var model string
	var headers map[string]string

	// Configure API based on provider
	switch c.config.APIProvider {
	case "openai":
		apiURL = "https://api.openai.com/v1/chat/completions"
		if c.config.OpenAIAPIBase != "" {
			apiURL = c.config.OpenAIAPIBase + "/v1/chat/completions"
		}
		apiKey = c.config.OpenAIAPIKey
		model = c.config.OpenAIModel
		headers = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + apiKey,
		}
		if c.config.OpenAIOrganization != "" {
			headers["OpenAI-Organization"] = c.config.OpenAIOrganization
		}
	case "groq":
		apiURL = "https://api.groq.com/openai/v1/chat/completions"
		apiKey = c.config.GroqAPIKey
		model = c.config.GroqModel
		headers = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + apiKey,
		}
	default:
		return "", fmt.Errorf("unsupported API provider: %s", c.config.APIProvider)
	}

	// Create request body
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	requestBody := ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: c.config.Temperature,
	}

	if c.config.OpenAIMaxTokens > 0 {
		requestBody.MaxTokens = c.config.OpenAIMaxTokens
	}

	// Marshal request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Create request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", err
	}

	// Check if we have choices
	if len(chatResponse.Choices) == 0 {
		return "", errors.New("no completions returned from API")
	}

	return chatResponse.Choices[0].Message.Content, nil
}

// GenerateShellCommand generates a shell command from a user prompt
func (c *Client) GenerateShellCommand(userPrompt, context string) (string, error) {
	// Create system prompt
	systemPrompt := "You are an expert at using shell commands. I need you to provide a response in the format `{\"command\": \"your_shell_command_here\"}`. Only provide a single executable line of shell code as the value for the \"command\" key. Never output any text outside the JSON structure. The command will be directly executed in a shell."

	// Add platform information
	platformInfo := getPlatformInfo()
	systemPrompt += " " + platformInfo

	// Add context if available
	if context != "" {
		systemPrompt += fmt.Sprintf(" Between [], these are the last %d tokens from the previous command's output, you can use them as context: [%s]",
			len(context), context)
	}

	// Generate completion
	userPromptWithPrefix := fmt.Sprintf("Generate a shell command that satisfies this user request: %s", userPrompt)
	return c.GenerateCompletion(systemPrompt, userPromptWithPrefix)
}

// getPlatformInfo returns information about the current platform
func getPlatformInfo() string {
	// This is a simplified version - in a real implementation, you would use
	// more detailed platform detection like in the Python version
	return fmt.Sprintf("The system the shell command will be executed on is %s.", getOSName())
}

// getOSName returns the name of the operating system
func getOSName() string {
	// Simple OS detection - could be expanded with more detailed information
	switch {
	case strings.Contains(strings.ToLower(getOSRelease()), "darwin"):
		return "macOS"
	case strings.Contains(strings.ToLower(getOSRelease()), "linux"):
		return "Linux"
	case strings.Contains(strings.ToLower(getOSRelease()), "windows"):
		return "Windows"
	default:
		return "Unknown"
	}
}

// getOSRelease returns the OS release information
func getOSRelease() string {
	return runtime.GOOS
}
