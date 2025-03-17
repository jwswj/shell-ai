package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// MaxContextTokens is the maximum number of tokens to keep in context
const MaxContextTokens = 1500

// CommandResponse represents the parsed command from the LLM response
type CommandResponse struct {
	Command string `json:"command"`
}

// ContextManager manages the context for the LLM
type ContextManager struct {
	tokenBuffer []rune
	maxTokens   int
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		tokenBuffer: make([]rune, 0, MaxContextTokens),
		maxTokens:   MaxContextTokens,
	}
}

// AddToken adds a token to the context
func (cm *ContextManager) AddToken(token rune) {
	if len(cm.tokenBuffer) >= cm.maxTokens {
		// Remove the first token
		cm.tokenBuffer = cm.tokenBuffer[1:]
	}
	cm.tokenBuffer = append(cm.tokenBuffer, token)
}

// Flush clears the context
func (cm *ContextManager) Flush() {
	cm.tokenBuffer = make([]rune, 0, cm.maxTokens)
}

// AddChunk adds a chunk of text to the context
func (cm *ContextManager) AddChunk(chunk string) {
	cm.Flush()
	for _, c := range chunk {
		cm.AddToken(c)
	}
}

// GetContext returns the current context
func (cm *ContextManager) GetContext() string {
	if len(cm.tokenBuffer) == 0 {
		return ""
	}
	return string(cm.tokenBuffer)
}

// ParseLLMResponse parses the LLM response to extract the command
func ParseLLMResponse(response string) (string, error) {
	// Try to extract JSON from markdown code blocks
	jsonContent := extractJSONFromMarkdown(response)
	if jsonContent == "" {
		// If no code blocks found, try to parse the whole response as JSON
		jsonContent = response
	}

	// Parse the JSON
	var commandResp CommandResponse
	err := json.Unmarshal([]byte(jsonContent), &commandResp)
	if err != nil {
		return "", err
	}

	return commandResp.Command, nil
}

// extractJSONFromMarkdown extracts JSON content from markdown code blocks
func extractJSONFromMarkdown(markdown string) string {
	// Try to find JSON in code blocks
	codeBlockRegex := regexp.MustCompile("```(?:json)?\\s*\\n([\\s\\S]*?)\\n```")
	matches := codeBlockRegex.FindAllStringSubmatch(markdown, -1)

	if len(matches) > 0 {
		// Return the content of the first code block
		return strings.TrimSpace(matches[0][1])
	}

	// If no code blocks found, try to find inline code
	inlineCodeRegex := regexp.MustCompile("`([^`]*)`")
	inlineMatches := inlineCodeRegex.FindAllStringSubmatch(markdown, -1)

	if len(inlineMatches) > 0 {
		// Return the content of the first inline code
		return strings.TrimSpace(inlineMatches[0][1])
	}

	return ""
}
