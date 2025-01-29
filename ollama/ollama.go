package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
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
	url     string
	content string
}

type OllamaServiceInterface interface {
	requestOllama(req Request) (*Response, error)
	GenerateSummaries(content string) string
}

func NewOllamaService(url string, content string) OllamaServiceInterface {
	return &OllamaService{
		url:     url,
		content: content,
	}
}

func (s *OllamaService) requestOllama(req Request) (*Response, error) {
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

func (s *OllamaService) GenerateSummaries(content string) string {
	prompt := fmt.Sprintf(`
The information below is a summary of PRs, issues, and commits made by the user over a period of time.

Use this to summarize in 3 sections what the user did during that time period.

Instead of dividing things into Commits, Issues, Pull Requests, etc., divide your sections by the subject of what you did.

Please use an appropriate emoji at the beginning of the section title.

---
%s
---

The format is as follows:

%smarkdown

# Summary Title <!-- Please output the target date from the PR, Issue, or Commit date. -->

## [emoji1] Section1 Title

Section1 Content

## [emoji2] Section2 Title

Section2 Content

## [emoji3] Section3 Title

Section3 Content

%s
`, "```", content, "```")

	msg := Message{
		Role:    "user",
		Content: prompt,
	}
	req := Request{
		Model:    "phi4",
		Stream:   false,
		Messages: []Message{msg},
	}
	resp, err := s.requestOllama(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resp.Message.Content
}
