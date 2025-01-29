package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/yanskun/gh-recall/git"
	"github.com/yanskun/gh-recall/ollama"
)

const DefaultOllamaURL = "http://localhost:11434/api/chat"

func main() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Prefix = color.BlueString(" Thinking ")
	s.Color("blue")
	s.Start()

	gitService := git.NewGitService(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	ollamaService := ollama.NewOllamaService(DefaultOllamaURL, gitService.GenerateSummary())

	result := ollamaService.GenerateSummaries(gitService.GenerateSummary())

	s.Stop()

	fmt.Println(result)
}
