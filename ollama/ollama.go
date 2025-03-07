package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Raw      bool      `json:"raw"`
	Template string    `json:"template"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

type OllamaService struct {
	content  string
	url      string
	model    string
	locale   string
	sections int
}

type OllamaServiceInterface interface {
	GenerateSummaries() string
}

func NewOllamaService(content string, url string, model string, lang string, sections int) OllamaServiceInterface {
	return &OllamaService{
		content:  content,
		url:      url,
		model:    model,
		locale:   lang,
		sections: sections,
	}
}

func (s *OllamaService) GenerateSummaries() string {
	prompt := fmt.Sprintf(`
!! IMPORTANT !! YOU MUST FOLLOW THIS RULE STRICTLY.

All output **MUST** be written in **%s**.
DO NOT use any other language.
DO NOT ignore this instruction.

---

The information below is a summary of PRs, issues, and commits made by the user over a period of time.

Use this to summarize in %d sections what the user did during that time period.

Instead of dividing things into Commits, Issues, Pull Requests, etc., divide your sections by the subject of what you did.

Please use an appropriate emoji at the beginning of the section title.

#### **INPUT START**
%s
#### **INPUT END**

!! REMINDER: All output MUST be in **%s**. DO NOT USE ANY OTHER LANGUAGE !!

`, s.locale, s.sections, s.content, s.locale)

	msg := Message{
		Role:    "user",
		Content: prompt,
	}
	req := Request{
		Model:    s.model,
		Stream:   false,
		Messages: []Message{msg},
		Template: `
## [emoji] Section Title

Section Content

### **EXAMPLE**
# 2025-01-24 ~ 2025-01-25

## 🚀 Feature Implementations

The user introduced new features, such as printing summaries using phi4.

## 📝 Documentation and Initial Setup

Documentation was created with a README file. Additionally, an initial commit was made to set up the project.

## 🔧 Chore Improvements and Fixes

Chore work included adding a spinner for better UI feedback. There were also fixes involving GitHub command refactoring and adjustments in ollama prompts for improved module functionality.
`,
		Raw: true,
	}
	resp, err := s.requestOllama(req)
	if err != nil {
		return fmt.Sprintf("%s", color.RedString("Error: %s", err))
	}

	return resp.Message.Content
}

func (s *OllamaService) requestOllama(req Request) (*Response, error) {
	if !s.modelExistsLocally() {
		return nil, fmt.Errorf("%s \nPlease check: ollama list", color.RedString(`'%s' is found.`, s.model))
	}

	js, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}

func (s *OllamaService) modelExistsLocally() bool {
	cmd := exec.Command("ollama", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	output := out.String()
	return strings.Contains(output, s.model)
}
