package main

import (
	"fmt"
	"time"

	"github.com/yanskun/gh-recall/git"
	"github.com/yanskun/gh-recall/ollama"
)

const DefaultOllamaURL = "http://localhost:11434/api/chat"

func main() {
	gitService := git.NewGitService(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	ollamaService := ollama.NewOllamaService(DefaultOllamaURL, gitService.GenerateSummary())

	result := ollamaService.GenerateSummaries(gitService.GenerateSummary())

	fmt.Println(result)
}
