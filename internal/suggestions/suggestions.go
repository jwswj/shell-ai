package suggestions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jwswj/shell-ai/internal/config"
	"github.com/jwswj/shell-ai/internal/llm"
	"github.com/jwswj/shell-ai/internal/parser"
	"github.com/manifoldco/promptui"
)

// SystemOption represents a system option in the suggestions menu
type SystemOption string

// System options
const (
	OptGenSuggestions SystemOption = "Generate new suggestions"
	OptDismiss        SystemOption = "Dismiss"
	OptNewCommand     SystemOption = "Enter a new command"
)

// TextEditors is a list of common text editors
var TextEditors = []string{"vi", "vim", "emacs", "nano", "ed", "micro", "joe", "nvim"}

// ContextManager is the global context manager
var ContextManager = parser.NewContextManager()

// Run runs the suggestions engine
func Run(client *llm.Client, cfg *config.Config, promptArgs []string) error {
	// Join prompt arguments into a single string
	prompt := strings.Join(promptArgs, " ")

	// Show warning if context mode is enabled
	if cfg.ContextMode {
		fmt.Printf("WARNING Context mode: data will be sent to the LLM, be careful if any sensitive data...\n\n")
		fmt.Printf(">>> %s\n", getCurrentDir())
	}

	for {
		// Generate suggestions
		suggestions, err := generateSuggestions(client, cfg, prompt)
		if err != nil {
			return err
		}

		// Add only the Dismiss system option
		options := append(suggestions, string(OptDismiss))

		// Create a select prompt with promptui
		selectPrompt := promptui.Select{
			Label: "Select a command",
			Items: options,
			Size:  10, // Show 10 items at a time
			Templates: &promptui.SelectTemplates{
				Active:   "→ {{ if eq . \"Dismiss\" }}{{ . | red }}{{ else }}{{ . | cyan }}{{ end }}",
				Inactive: "  {{ if eq . \"Dismiss\" }}{{ . | red }}{{ else }}{{ . }}{{ end }}",
				Selected: "✓ {{ if eq . \"Dismiss\" }}{{ . | red }}{{ else }}{{ . | green }}{{ end }}",
			},
			Searcher: func(input string, index int) bool {
				option := options[index]
				return strings.Contains(strings.ToLower(option), strings.ToLower(input))
			},
		}

		_, selection, err := selectPrompt.Run()
		if err != nil {
			// Check if the error is due to Ctrl+C (interrupt)
			if err.Error() == "^C" || strings.Contains(err.Error(), "interrupt") {
				fmt.Println("\nExiting...")
				return nil
			}
			return err
		}

		// Handle selection
		switch SystemOption(selection) {
		case OptDismiss:
			return nil
		default:
			// User selected a command
			userCommand := selection

			// Confirm command if not skipping confirmation
			if !cfg.SkipConfirm {
				confirmPrompt := promptui.Prompt{
					Label:     fmt.Sprintf("Confirm [%s]", userCommand),
					Default:   userCommand,
					AllowEdit: true,
				}

				confirmedCommand, err := confirmPrompt.Run()
				if err != nil {
					// Check if the error is due to Ctrl+C (interrupt)
					if err.Error() == "^C" || strings.Contains(err.Error(), "interrupt") {
						fmt.Println("\nExiting...")
						return nil
					}
					return err
				}
				userCommand = confirmedCommand
			}

			// Write to shell history if not skipping history
			if !cfg.SkipHistory {
				err = writeToShellHistory(userCommand)
				if err != nil {
					fmt.Printf("Warning: %s\n", err)
				}
			}

			// Execute command
			if !cfg.ContextMode {
				// Default mode - execute and exit
				cmd := exec.Command("sh", "-c", userCommand)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					fmt.Printf("Error executing command: %v\n", err)
				}
				return nil
			} else {
				// Context mode - capture output and continue
				if startsWithAny(userCommand, TextEditors) {
					// For text editors, just run the command directly
					cmd := exec.Command("sh", "-c", userCommand)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						fmt.Printf("Error executing command: %v\n", err)
					}
				} else if strings.HasPrefix(userCommand, "cd") {
					// Handle cd command specially
					path := strings.TrimSpace(strings.TrimPrefix(userCommand, "cd"))
					path = os.ExpandEnv(path)
					path = filepath.Clean(path)
					err = os.Chdir(path)
					if err != nil {
						fmt.Printf("Error changing directory: %v\n", err)
					}
				} else {
					// For other commands, capture output
					cmd := exec.Command("sh", "-c", userCommand)
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Printf("Error executing command: %v\n", err)
					}
					if len(output) > 0 {
						fmt.Printf("\n%s", string(output))
					}
					ContextManager.AddChunk(string(output))
				}

				// Prompt for new command
				fmt.Printf(">>> %s\n", getCurrentDir())
				newCmdPrompt := promptui.Prompt{
					Label: "New command",
					Validate: func(input string) error {
						if strings.TrimSpace(input) == "" {
							return fmt.Errorf("Command cannot be empty")
						}
						return nil
					},
				}

				newCmd, err := newCmdPrompt.Run()
				if err != nil {
					// Check if the error is due to Ctrl+C (interrupt)
					if err.Error() == "^C" || strings.Contains(err.Error(), "interrupt") {
						fmt.Println("\nExiting...")
						return nil
					}
					return err
				}
				prompt = strings.TrimSpace(newCmd)
			}
		}
	}
}

// generateSuggestions generates shell command suggestions
func generateSuggestions(client *llm.Client, cfg *config.Config, prompt string) ([]string, error) {
	// Generate suggestions in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	suggestions := make([]string, 0, cfg.SuggestionCount)
	errors := make([]error, 0)

	// Limit concurrency to 4
	maxWorkers := 4
	if cfg.SuggestionCount < maxWorkers {
		maxWorkers = cfg.SuggestionCount
	}

	// Create a semaphore channel to limit concurrency
	sem := make(chan struct{}, maxWorkers)

	for i := 0; i < cfg.SuggestionCount; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func() {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			// Get context if enabled
			var context string
			if cfg.ContextMode {
				context = ContextManager.GetContext()
			}

			// Generate suggestion
			response, err := client.GenerateShellCommand(prompt, context)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			// Parse response
			command, err := parser.ParseLLMResponse(response)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			// Add suggestion
			if command != "" {
				mu.Lock()
				suggestions = append(suggestions, command)
				mu.Unlock()
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check if we have any suggestions
	if len(suggestions) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("failed to generate suggestions: %v", errors[0])
	}

	// Deduplicate suggestions
	return deduplicate(suggestions), nil
}

// deduplicate removes duplicate strings from a slice
func deduplicate(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// writeToShellHistory writes a command to the shell history
func writeToShellHistory(command string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("SHELL environment variable not set")
	}

	var historyFilePath string
	var historyFormat string

	switch {
	case strings.Contains(shell, "zsh"):
		historyFilePath = filepath.Join(os.Getenv("HOME"), ".zsh_history")
		historyFormat = ": %d:0;%s\n"
	case strings.Contains(shell, "bash"):
		historyFilePath = filepath.Join(os.Getenv("HOME"), ".bash_history")
		historyFormat = "%s\n"
	case strings.Contains(shell, "csh"), strings.Contains(shell, "tcsh"):
		historyFilePath = filepath.Join(os.Getenv("HOME"), ".history")
		historyFormat = "%s\n"
	case strings.Contains(shell, "ksh"):
		historyFilePath = filepath.Join(os.Getenv("HOME"), ".sh_history")
		historyFormat = "%s\n"
	case strings.Contains(shell, "fish"):
		historyFilePath = filepath.Join(os.Getenv("HOME"), ".local/share/fish/fish_history")
		historyFormat = "- cmd: %s\n  when: %d\n"
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// Open history file
	file, err := os.OpenFile(historyFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write command to history
	timestamp := time.Now().Unix()
	var entry string
	if strings.Contains(shell, "zsh") {
		entry = fmt.Sprintf(historyFormat, timestamp, command)
	} else if strings.Contains(shell, "fish") {
		entry = fmt.Sprintf(historyFormat, command, timestamp)
	} else {
		entry = fmt.Sprintf(historyFormat, command)
	}

	_, err = file.WriteString(entry)
	return err
}

// getCurrentDir returns the current directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

// startsWithAny checks if a string starts with any of the given prefixes
func startsWithAny(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
