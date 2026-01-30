package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v57/github"
)

type URLType string

const (
	URLTypeRepo URLType = "repository"
	URLTypePR   URLType = "pull_request"
)

type GitHubURLInfo struct {
	Owner    string
	Repo     string
	Type     URLType
	Branch   string
	FilePath string
	Line     int
	PRNumber int
}

type PullRequestInfo struct {
	Number    int
	Title     string
	HeadOwner string
	HeadBranch string
	BaseOwner string
	BaseBranch string
}

type GitHubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) *GitHubClient {
	var client *github.Client
	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return &GitHubClient{client: client}
}

func (gc *GitHubClient) GetPullRequest(owner, repo string, number int) (*PullRequestInfo, error) {
	ctx := context.Background()
	pr, _, err := gc.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	return &PullRequestInfo{
		Number:     pr.GetNumber(),
		Title:      pr.GetTitle(),
		HeadOwner:  pr.GetHead().GetUser().GetLogin(),
		HeadBranch: pr.GetHead().GetRef(),
		BaseOwner:  pr.GetBase().GetUser().GetLogin(),
		BaseBranch: pr.GetBase().GetRef(),
	}, nil
}

// ParseGitHubURL 解析各种 GitHub URL 格式
func ParseGitHubURL(url string) (*GitHubURLInfo, error) {
	// 移除尾部斜杠
	url = strings.TrimSuffix(url, "/")

	// 正则表达式匹配不同的 URL 格式
	patterns := []struct {
		regex   *regexp.Regexp
		handler func(matches []string) (*GitHubURLInfo, error)
	}{
		{
			// Pull Request: https://github.com/owner/repo/pull/123
			// Pull Request with files: https://github.com/owner/repo/pull/123/files
			regex: regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/pull/(\d+)`),
			handler: func(matches []string) (*GitHubURLInfo, error) {
				prNum, _ := strconv.Atoi(matches[3])
				return &GitHubURLInfo{
					Owner:    matches[1],
					Repo:     matches[2],
					Type:     URLTypePR,
					PRNumber: prNum,
				}, nil
			},
		},
		{
			// File with line: https://github.com/owner/repo/blob/branch/path/to/file.go#L123
			regex: regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/blob/([^/]+)/(.+?)(?:#L(\d+))?$`),
			handler: func(matches []string) (*GitHubURLInfo, error) {
				info := &GitHubURLInfo{
					Owner:    matches[1],
					Repo:     matches[2],
					Type:     URLTypeRepo,
					Branch:   matches[3],
					FilePath: matches[4],
				}
				if matches[5] != "" {
					info.Line, _ = strconv.Atoi(matches[5])
				}
				return info, nil
			},
		},
		{
			// Tree (directory): https://github.com/owner/repo/tree/branch/path/to/dir
			regex: regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/tree/([^/]+)(?:/(.+))?$`),
			handler: func(matches []string) (*GitHubURLInfo, error) {
				return &GitHubURLInfo{
					Owner:    matches[1],
					Repo:     matches[2],
					Type:     URLTypeRepo,
					Branch:   matches[3],
					FilePath: matches[4],
				}, nil
			},
		},
		{
			// Repository: https://github.com/owner/repo
			regex: regexp.MustCompile(`github\.com/([^/]+)/([^/]+)$`),
			handler: func(matches []string) (*GitHubURLInfo, error) {
				return &GitHubURLInfo{
					Owner: matches[1],
					Repo:  matches[2],
					Type:  URLTypeRepo,
				}, nil
			},
		},
	}

	// 尝试匹配每个模式
	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(url)
		if matches != nil {
			return pattern.handler(matches)
		}
	}

	return nil, fmt.Errorf("unsupported GitHub URL format: %s", url)
}
