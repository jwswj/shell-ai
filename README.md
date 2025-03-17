# Shell-AI: let AI write your shell commands

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Shell-AI (`shai`) is a CLI utility that brings the power of natural language understanding to your command line. Simply input what you want to do in natural language, and `shai` will suggest single-line commands that achieve your intent.

This is a Go fork of [@ricklamers's original project](https://github.com/ricklamers/shell-ai).

## Features

- Generate shell commands from natural language descriptions
- Multiple command suggestions to choose from
- Support for OpenAI, and Groq LLM providers
- Context mode to maintain command history and output for better suggestions
- Shell history integration
- Configurable via environment variables or config file

## Installation

### Using Go Install (Recommended)

The easiest way to install Shell-AI is using Go's built-in install command:

```bash
go install github.com/jwswj/shell-ai/cmd/shai@latest
```

This will install the `shai` binary to your `$GOPATH/bin` directory, which should be in your PATH.

You can then run it using:

```bash
shai find all files modified in the last 24 hours
```

### From Source

Alternatively, you can build from source:

```bash
git clone https://github.com/jwswj/shell-ai.git
cd shell-ai
make build
```

Then run the binary directly:

```bash
./bin/shai find all files modified in the last 24 hours
```

Or install it using:

```bash
make install
```

## Configuration

Shell-AI can be configured using environment variables or a config file.

### Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key
- `GROQ_API_KEY`: Your Groq API key
- `OPENAI_MODEL`: The OpenAI model to use (default: `gpt-3.5-turbo`)
- `GROQ_MODEL`: The Groq model to use (default: `llama-3.3-70b-versatile`)
- `SHAI_API_PROVIDER`: The API provider to use (`openai`, or `groq`, default: `groq`)
- `SHAI_SUGGESTION_COUNT`: The number of suggestions to generate (default: `3`)
- `SHAI_SKIP_CONFIRM`: Skip confirmation of the command to execute (default: `false`)
- `SHAI_SKIP_HISTORY`: Skip writing selected command to shell history (default: `false`)
- `SHAI_TEMPERATURE`: Controls randomness in the output (default: `0.05`)
- `CTX`: Enable context mode (default: `false`)
- `DEBUG`: Enable debug mode (default: `false`)

### Config File

You can also create a config file at `~/.config/shell-ai/config.json` (Linux/macOS) or `%APPDATA%\shell-ai\config.json` (Windows):

```json
{
  "SHAI_SUGGESTION_COUNT": "3",
  "SHAI_API_PROVIDER": "groq",
  "GROQ_API_KEY": "your-groq-api-key",
  "GROQ_MODEL": "llama-3.3-70b-versatile",
  "SHAI_TEMPERATURE": "0.05"
}
```

## Usage

To use Shell-AI, open your terminal and type:

```bash
shai find all files modified in the last 24 hours
```

Shell-AI will generate several command suggestions, and you can select one to execute.

### Context Mode

Context mode allows Shell-AI to maintain context between commands, which can be useful for complex tasks:

```bash
shai --ctx find all log files
```

In context mode, the output of each command is captured and used as context for the next command.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
