# gh-recall

GitHub CLI commands extension.

A GitHub CLI Extension that retrieves and summarizes your recent activities, including Pull Requests, Issues, and Commits.

![image](https://github.com/user-attachments/assets/c2da796c-e773-4cc4-beaa-709afb8f8abc)


## Requirements

- gh (GitHub CLI) - must be installed
- [Ollama](https://ollama.com/) - must be installed for AI-based summarization.
- Model - Default model: `phi4`, but you can use another model like `"mistral"`, `"llama3"`, etc

### Recommended Model

The default model is phi4, and it is the only model tested so far.
Other models may work, but their performance and compatibility are not guaranteed.

To install the recommended model:

```bash
ollama pull phi4
```

## Install

```bash
gh extension install yanskun/gh-recall
```

## Usage

```bash
gh recall [options]
```

### Configuration (`config.toml`)

You can configure default values using a config.toml file.
The configuration file is stored in:

Linux/macOS: `~/.config/gh-recall/config.toml`
Windows: `C:\Users\YourUser\.config\gh-recall\config.toml`
If no `config.toml` is found, it will be automatically generated with default values.

#### Example config.toml

```
days = 14
locale = "ja"
model = "mistral"
port = 11434
sections = 5
```

### Options

| Option             | Description                                                    | Default |
| ------------------ | -------------------------------------------------------------- | ------- |
| `-h`, `--help`     | Show help for the command.                                     | -       |
| `-d`, `--days`     | Number of days to look back when retrieving data.              | `7`     |
| `-l`, `--locale`   | Output language for the summary (en, ja, etc.).                | `en`    |
| `-m`, `--model`    | Ollama model to use for summarization. (`phi4 is recommended`) | `phi4`  |
| `-p`, `--port`     | Port number for Ollama connection.                             | `11434` |
| `-s`, `--sections` | Number of sections to display in the summary.                  | `3`     |

Priority order:
options > `config.toml` > Default values

### Examples

- Retrieve the last **7 days** of contributions (default):

```bash
gh recall
```

- Retrieve the last **30 days** of contributions:

```bash
gh recall --days 30
```

- Output the summary in **Japanese**:

```bash
gh recall --locale ja
```

- Use a **different Ollama model**:

```bash
gh recall --model mistral
```

- Change the number of sections in the summary:

```bash
gh recall --sections 5
```

- Specify the port number for Ollama:

```bash
gh recall --port 12345
```

## Output Example

When you run:

```bash
gh recall --days 7 --locale en --model phi4 --sections 3
```

You will get an output like this:

```markdown
# 2025-01-24 ~ 2025-01-25

## 🚀 Feature Implementations

The user introduced new features, such as printing summaries using phi4.

## 📝 Documentation and Initial Setup

Documentation was created with a README file. Additionally, an initial commit was made to set up the project.

## 🔧 Chore Improvements and Fixes

Chore work included adding a spinner for better UI feedback. There were also fixes involving GitHub command refactoring and adjustments in ollama prompts for improved module functionality.
```
