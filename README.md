# gh-recall

GitHub CLI commands extension.

A GitHub CLI Extension that retrieves and summarizes your recent activities, including Pull Requests, Issues, and Commits.

## Required

- [Ollama](https://ollama.com/) - must be installed for AI-based summarization.
  - [phi4](https://ollama.com/library/phi4) - model is required for the summarization to work. You can install it with

## Install

```shell
gh extension install yanskun/gh-recall
```

## Usage

```shell
gh recall [options]
```

### Options

| Option           | Description                                       | Default |
| ---------------- | ------------------------------------------------- | ------- |
| `-h`, `--help`   | Show help for the command.                        | -       |
| `-d`, `--days`   | Number of days to look back when retrieving data. | `7`     |
| `-l`, `--locale` | Output language for the summary (en, ja, etc.).   | `en`    |
| `-m`, `--model`  | Ollama model to use for summarization.            | `phi4`  |

### Examples

- Retrieve the last **7 days** of contributions (default):

```shell
gh recall
```

- Retrieve the last **30 days** of contributions:

```shell
gh recall --days 30
```

- Output the summary in **Japanese**:

```shell
gh recall --locale ja
```

- Use a **different Ollama model**:

```shell
gh recall --model mistral
```

## Output Example

When you run:

```bash
gh recall --days 7 --locale en --model phi4
```

You will get an output like this:

```markdown
# Summary 2025-01-20 ~ 2025-01-27

## 🚀 Implemented New Features
- Added dark mode support in the UI.
- Implemented API rate limiting for better security.

## 🛠 Fixed Bugs & Issues
- Fixed a critical bug in authentication flow.
- Resolved memory leaks in the background worker.

## 📖 Documentation & Code Refactoring
- Updated the README with clearer installation steps.
- Refactored database connection handling for better performance.
```
