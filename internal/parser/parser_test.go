package parser

import (
	"testing"
)

func TestParseLLMResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid JSON",
			response: `{"command": "ls -la"}`,
			want:     "ls -la",
			wantErr:  false,
		},
		{
			name:     "valid JSON in code block",
			response: "```json\n{\"command\": \"ls -la\"}\n```",
			want:     "ls -la",
			wantErr:  false,
		},
		{
			name:     "valid JSON in inline code",
			response: "`{\"command\": \"ls -la\"}`",
			want:     "ls -la",
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			response: "not a json",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLLMResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLLMResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLLMResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextManager(t *testing.T) {
	cm := NewContextManager()

	// Test empty context
	if cm.GetContext() != "" {
		t.Errorf("Expected empty context, got %q", cm.GetContext())
	}

	// Test adding a chunk
	cm.AddChunk("test chunk")
	if cm.GetContext() != "test chunk" {
		t.Errorf("Expected context to be 'test chunk', got %q", cm.GetContext())
	}

	// Test flushing
	cm.Flush()
	if cm.GetContext() != "" {
		t.Errorf("Expected empty context after flush, got %q", cm.GetContext())
	}

	// Test adding tokens
	for _, c := range "hello" {
		cm.AddToken(c)
	}
	if cm.GetContext() != "hello" {
		t.Errorf("Expected context to be 'hello', got %q", cm.GetContext())
	}
}
