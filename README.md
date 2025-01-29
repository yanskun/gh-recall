# gh-recall

GitHub CLI commands extension.

A GitHub CLI Extension that retrieves and summarizes your recent activities, including Pull Requests, Issues, and Commits.

```
$ gh recall
# 2025-01-24 ~ 2025-01-28

## 📄 Documentation Enhancements

The user focused on improving project documentation by creating a README file, which serves as an essential guide for understanding the project's purpose and setup. This task was completed on January 28, 2025.

## ✨ Feature Development

A new feature was developed to print summaries using the phi4 tool. This enhancement aimed at providing more detailed insights through summaries, completed on January 27, 2025.

## 🚀 Initial Setup

The user initiated the project with an initial commit on January 24, 2025. This marked the beginning of the development process and laid the foundation for subsequent contributions.
```

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
