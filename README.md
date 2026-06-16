# findlib

A CLI tool that finds libraries and APIs by natural language description. Built for use in nvim, VS Code, or any terminal.

```
findlib go "http client with retries"
findlib python "json parsing"
findlib js "state management"
findlib csharp "orm for postgresql"
findlib rust "async runtime"
```

## Installation

```bash
go install github.com/lstbob/findlibcli/cmd/findlib@latest
```

This installs the `findlib` binary to `~/go/bin/`. Make sure that's in your PATH:

```bash
export PATH=$PATH:~/go/bin
# Add to ~/.bashrc or ~/.zshrc to persist
```

Or build manually and move to `/usr/local/bin`:

```bash
git clone https://github.com/lstbob/findlibcli.git
cd findlibcli
go build -o findlib ./cmd/findlib
sudo mv findlib /usr/local/bin/
```

## Setup (optional — works without any configuration)

### For AI-powered search (recommended)

Get a free Gemini API key from [aistudio.google.com/apikey](https://aistudio.google.com/apikey) (no credit card required).

```bash
findlib config gemini "your-api-key-here"
```

### For live documentation snippets

Get a free Context7 API key from [context7.com/dashboard](https://context7.com/dashboard).

```bash
findlib config context7 "your-ctx7-api-key"
```

### Switch LLM provider

```bash
findlib config llm openai
findlib config openai "sk-..."
```

## Usage

```bash
findlib search <language> "<description>"
```

### Flags

| Flag | Description |
|---|---|
| `--json` | Output as JSON |
| `--import-only` | Print only the import path |
| `--no-copy` | Skip clipboard copy |

### nvim integration

```vim
" Quick search — shows docs in terminal
:!findlib go "http client"

" Insert import path into buffer
:r !findlib go "http client" --import-only

" With a keymap
nnoremap <leader>fl :!findlib
```

## How it works

1. **Resolve** — sends your description + language to Gemini (free) to suggest matching libraries
2. **Fetch docs** — looks up usage snippets from Context7
3. **Pick** — if multiple matches, shows an interactive picker
4. **Show** — displays usage docs + copies the import path to clipboard

If no API keys are configured, falls back to searching package registries directly (npm, crates.io, NuGet).
