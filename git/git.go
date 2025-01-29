package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cli/go-gh/v2"
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

type GitService struct {
	startDate string
	endDate   string
}

type GitServiceInterface interface {
	getPullRequests() []PullRequest
	getIssues() []Issue
	getCommits() string
	GenerateSummary() string
}

func NewGitService(startDate time.Time, endDate time.Time) GitServiceInterface {
	return &GitService{
		startDate: startDate.Format("2006-01-02"),
		endDate:   endDate.Format("2006-01-02"),
	}
}

func (s *GitService) getPullRequests() []PullRequest {
	query := fmt.Sprintf("created:%s..%s", s.startDate, s.endDate)
	result, _, err := gh.Exec("pr", "list", "--author", "@me", "--state", "all", "--json", "title,body,state,createdAt,isDraft", "--search", query)
	if err != nil {
		log.Fatal(err)
	}

	var pullRequests []PullRequest
	err = json.Unmarshal([]byte(result.Bytes()), &pullRequests)

	return pullRequests
}

func (s *GitService) getIssues() []Issue {
	query := fmt.Sprintf("created:%s..%s", s.startDate, s.endDate)
	result, _, err := gh.Exec("issue", "list", "--author", "@me", "--state", "all", "--json", "title,body,state,createdAt", "--search", query, "--repo", "cli/cli")

	if err != nil {
		log.Fatal(err)
	}

	var issues []Issue
	err = json.Unmarshal([]byte(result.Bytes()), &issues)

	return issues
}

func (s *GitService) getCommits() string {
	userName := executeCommand("git", "config", "--global", "user.name")

	return executeCommand("git", "log", "--pretty=format:%as: %s %b",
		fmt.Sprintf("--since=%s", s.startDate),
		fmt.Sprintf("--until=%s", s.endDate),
		fmt.Sprintf("--author=%s", userName),
		"--no-merges")
}

func (s *GitService) GenerateSummary() string {
	prs := s.getPullRequests()
	issues := s.getIssues()
	commits := s.getCommits()

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

func executeCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}

	return strings.TrimSpace(out.String())
}
