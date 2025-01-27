package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type PullRequest struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
	IsDraft   bool   `json:"isDraft"`
}

type Issue struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // 1週間前
	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")   // 昨日

	owner := executeCommand("gh", "repo", "view", "--json", "owner", "--jq", ".owner.login")

	prs := fetchPRs(owner, startDate, endDate)

	issues := fetchIssues(owner, startDate, endDate)

	commits := executeCommand("git", "log", "--pretty=format:%h %s", fmt.Sprintf("--since=%s", startDate), fmt.Sprintf("--until=%s", endDate))

	summary := generateSummary(prs, issues, commits)

	ollamaSummaries(summary)
}

func fetchPRs(owner, startDate, endDate string) []PullRequest {
	output := executeCommand("gh", "pr", "list", "--author", owner, "--json", "title,body,state,createdAt,isDraft", "--search", fmt.Sprintf("created:%s..%s", startDate, endDate))
	var prs []PullRequest
	json.Unmarshal([]byte(output), &prs)
	return prs
}

func fetchIssues(owner, startDate, endDate string) []Issue {
	output := executeCommand("gh", "issue", "list", "--author", owner, "--json", "title,body,state,createdAt", "--search", fmt.Sprintf("created:%s..%s", startDate, endDate))
	var issues []Issue
	json.Unmarshal([]byte(output), &issues)
	return issues
}

func executeCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
	return out.String()
}

func generateSummary(prs []PullRequest, issues []Issue, commits string) string {
	var summary bytes.Buffer
	summary.WriteString("Pull Requests:\n")
	for _, pr := range prs {
		summary.WriteString(fmt.Sprintf("- %s (State: %s, Draft: %t)\n", pr.Title, pr.State, pr.IsDraft))
	}
	summary.WriteString("\nIssues:\n")
	for _, issue := range issues {
		summary.WriteString(fmt.Sprintf("- %s (State: %s)\n", issue.Title, issue.State))
	}
	summary.WriteString("\nCommits:\n")
	summary.WriteString(commits)
	return summary.String()
}

const defaultOllamaURL = "http://localhost:11434/api/chat"

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

func ollamaSummaries(content string) {
	prompt := fmt.Sprintf(`	
The information below is a summary of PRs, issues, and commits made by the user over a period of time.

Use this to summarize in 3 sections what the user did during that time period.

Instead of dividing things into Commits, Issues, Pull Requests, etc., divide your sections by the subject of what you did.

Please use an appropriate emoji at the beginning of the section title.

---
%s
---
`, content)

	msg := Message{
		Role:    "user",
		Content: prompt,
	}
	req := Request{
		Model:    "phi4",
		Stream:   false,
		Messages: []Message{msg},
	}
	resp, err := talkToOllama(defaultOllamaURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(resp.Message.Content)
}

func talkToOllama(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
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
