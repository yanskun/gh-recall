package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/yanskun/gh-recall/git"
	"github.com/yanskun/gh-recall/ollama"
)

const DefaultOllamaURL = "http://localhost:11434/api/chat"
const DefaultOllamaModel = "phi4"

func main() {
	var helpFlag bool
	var modelVal string
	var localeVal string
	var daysVal int

	pflag.BoolVarP(&helpFlag, "help", "h", false, "Show help for command")
	pflag.StringVarP(&modelVal, "model", "m", DefaultOllamaModel, "Used model for Ollama")
	pflag.StringVarP(&localeVal, "locale", "l", "en", "Output language")
	pflag.IntVarP(&daysVal, "days", "d", 7, "Number of days to look back")
	pflag.Parse()

	if helpFlag {
		pflag.Usage()
		return
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Prefix = color.BlueString(" Thinking ")
	s.Color("blue")
	s.Start()

	gitService := git.NewGitService(time.Now().AddDate(0, 0, -daysVal), time.Now().AddDate(0, 0, -1))

	ollamaService := ollama.NewOllamaService(gitService.GenerateSummary(), DefaultOllamaURL, modelVal, localeVal)

	result := ollamaService.GenerateSummaries()

	s.Stop()

	fmt.Println(result)
}
