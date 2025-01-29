package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/yanskun/gh-recall/git"
	"github.com/yanskun/gh-recall/ollama"
	"github.com/yanskun/pflag"
)

const DefaultOllamaURL = "http://localhost:11434/api/chat"
const DefaultOllamaModel = "phi4"

func main() {
	var helpFlag bool
	var modelVal string
	var localeVal string

	pflag.BoolVarP(&helpFlag, "help", "h", false, "Show help for command")
	pflag.StringVarP(&modelVal, "model", "m", DefaultOllamaModel, "Used model for Ollama")
	pflag.StringVarP(&localeVal, "locale", "l", "en", "Output language")
	pflag.Parse()

	if helpFlag {
		pflag.Usage()
		return
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Prefix = color.BlueString(" Thinking ")
	s.Color("blue")
	s.Start()

	gitService := git.NewGitService(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	ollamaService := ollama.NewOllamaService(gitService.GenerateSummary(), DefaultOllamaURL, modelVal, localeVal)

	result := ollamaService.GenerateSummaries()

	s.Stop()

	fmt.Println(result)
}
