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
go install github.com/lstbob/findlibcli/cmd/findlib@main
```

If the proxy cache hasn't updated yet, `@main` pulls directly from the branch. Once the tag propagates, `@latest` will also work.

The binary installs to `~/go/bin/findlib`. Make sure `~/go/bin` is at the **front** of your PATH:

```bash
export PATH=$HOME/go/bin:$PATH
echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
```

Or build from source:

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
findlib <language> "<description>"
# or
findlib search <language> "<description>"
```

### Flags

| Flag | Description |
|---|---|
| `--json` | Output as JSON (uses first result, no picker) |
| `--import-only` | Print only the import path/name |
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

1. **Resolve** — sends description + language to Gemini (free) to suggest matching libraries
2. **Fetch docs** — looks up usage snippets from Context7 (if API key configured)
3. **Pick** — if multiple matches, shows an interactive picker (TTY) or numbered list (pipe)
4. **Show** — displays usage docs + copies the import path to clipboard

If no API keys are configured, falls back to searching package registries directly.

### Supported registries (offline mode)

| Language | Registry |
|---|---|
| JavaScript / TypeScript | npm |
| Python | npm (fallback, no PyPI API) |
| Go | npm (fallback) |
| Rust | crates.io |
| C# / .NET | NuGet |

## Config

Configuration is stored at `~/.config/findlib/config.json`. View or set values:

```bash
findlib config              # view all
findlib config gemini "key" # set Gemini API key
findlib config llm groq     # switch provider
```
