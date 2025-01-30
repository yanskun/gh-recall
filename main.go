package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/yanskun/gh-recall/config"
	"github.com/yanskun/gh-recall/git"
	"github.com/yanskun/gh-recall/ollama"
)

func main() {
	var helpFlag bool
	var modelVal string
	var localeVal string
	var daysVal int
	var portVal int
	var sectionsVal int

	c := config.LoadConfig()

	pflag.BoolVarP(&helpFlag, "help", "h", false, "Show help for command")
	pflag.StringVarP(&modelVal, "model", "m", c.Model, "Ollama model to use for summarization.")
	pflag.StringVarP(&localeVal, "locale", "l", c.Locale, "Output language for the summary.")
	pflag.IntVarP(&daysVal, "days", "d", c.Days, "Number of days to look back when retrieving data.")
	pflag.IntVarP(&portVal, "port", "p", c.Port, "Port number for Ollama connection.")
	pflag.IntVarP(&sectionsVal, "sections", "s", c.Sections, "Number of sections to display in the summary.")
	pflag.Parse()

	if helpFlag {
		pflag.Usage()
		return
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Prefix = color.BlueString(" Thinking ")
	s.Suffix = "\n"
	s.Color("blue")
	s.Start()

	gitService := git.NewGitService(time.Now().AddDate(0, 0, -daysVal), time.Now().AddDate(0, 0, -1))

	ollamaService := ollama.NewOllamaService(gitService.GenerateSummary(), fmt.Sprintf("http://localhost:%d/api/chat", portVal), modelVal, localeVal, sectionsVal)

	result := ollamaService.GenerateSummaries()

	s.Stop()

	fmt.Println(result)
}
