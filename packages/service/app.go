package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	AppVersion      = "1.1.0"
	DefaultPort     = 9527
	DefaultCacheDir = ".github-browser/repos"
)

var startedAt = time.Now()

type Service struct {
	config    *Config
	cacheDir  string
	gitClient *GitClient
	ghClient  *GitHubClient
}

type OpenRequest struct {
	URL      string `json:"url" binding:"required"`
	IDE      string `json:"ide"`
	FilePath string `json:"filePath"`
	Line     int    `json:"line"`
}

type OpenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
	Mode    string `json:"mode"`
}

func NewService() (*Service, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v, using defaults", err)
		config = DefaultConfig()
	}

	return NewServiceWithConfig(config)
}

func NewServiceWithConfig(config *Config) (*Service, error) {
	normalized := normalizeConfig(config)
	if err := os.MkdirAll(normalized.CacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Service{
		config:    normalized,
		cacheDir:  normalized.CacheDir,
		gitClient: NewGitClient(normalized.CacheDir),
		ghClient:  NewGitHubClient(normalized.GitHubToken),
	}, nil
}

func normalizeConfig(config *Config) *Config {
	if config == nil {
		config = DefaultConfig()
	}

	normalized := *config
	if normalized.Port == 0 {
		normalized.Port = DefaultPort
	}
	if normalized.DefaultIDE == "" {
		normalized.DefaultIDE = "code"
	}
	if normalized.CacheDir == "" {
		normalized.CacheDir = filepath.Join(os.Getenv("HOME"), DefaultCacheDir)
	}

	return &normalized
}

func (s *Service) Health(mode string) HealthResponse {
	return HealthResponse{
		Status:  "ok",
		Version: AppVersion,
		Uptime:  time.Since(startedAt).Round(time.Second).String(),
		Mode:    mode,
	}
}

func (s *Service) Open(req OpenRequest) (*OpenResponse, error) {
	log.Printf("Received request: %s", req.URL)

	info, err := ParseGitHubURL(req.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid GitHub URL: %w", err)
	}

	log.Printf("Parsed: owner=%s, repo=%s, type=%s", info.Owner, info.Repo, info.Type)

	var repoPath string
	switch info.Type {
	case URLTypeRepo:
		repoPath, err = s.handleRepository(info)
	case URLTypePR:
		repoPath, err = s.handlePullRequest(info)
	default:
		err = fmt.Errorf("unsupported URL type: %s", info.Type)
	}
	if err != nil {
		return nil, err
	}

	targetPath := repoPath
	if req.FilePath != "" {
		targetPath = filepath.Join(repoPath, req.FilePath)
	} else if info.FilePath != "" {
		targetPath = filepath.Join(repoPath, info.FilePath)
	}

	line := req.Line
	if line == 0 && info.Line > 0 {
		line = info.Line
	}

	ide := req.IDE
	if ide == "" {
		ide = s.config.DefaultIDE
	}

	log.Printf("Opening in %s: %s (line: %d)", ide, targetPath, line)
	if err := OpenInIDE(ide, targetPath, line); err != nil {
		return nil, fmt.Errorf("failed to open IDE: %w", err)
	}

	return &OpenResponse{
		Status:  "ok",
		Message: "Opened successfully",
		Path:    repoPath,
	}, nil
}

func (s *Service) UpdateConfig(newConfig *Config) error {
	normalized := normalizeConfig(newConfig)
	if err := os.MkdirAll(normalized.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	if err := SaveConfig(normalized); err != nil {
		return err
	}

	s.config = normalized
	s.cacheDir = normalized.CacheDir
	s.gitClient = NewGitClient(normalized.CacheDir)
	s.ghClient = NewGitHubClient(normalized.GitHubToken)
	return nil
}

func (s *Service) handleRepository(info *GitHubURLInfo) (string, error) {
	repoPath := s.config.GetRepoPath(info.Owner, info.Repo)

	if _, err := os.Stat(repoPath); err == nil {
		log.Printf("Repository exists, updating...")
		if err := s.gitClient.Pull(repoPath); err != nil {
			log.Printf("Warning: git pull failed: %v", err)
		}
	} else {
		log.Printf("Cloning repository...")
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", info.Owner, info.Repo)
		if err := s.gitClient.Clone(repoURL, repoPath); err != nil {
			return "", fmt.Errorf("failed to clone: %w", err)
		}
	}

	if info.Branch != "" {
		log.Printf("Checking out branch/tag: %s", info.Branch)
		if err := s.gitClient.Fetch(repoPath); err != nil {
			log.Printf("Warning: git fetch failed: %v", err)
		}
		if err := s.gitClient.Checkout(repoPath, info.Branch); err != nil {
			return "", fmt.Errorf("failed to checkout %s: %w", info.Branch, err)
		}
	}

	return repoPath, nil
}

func (s *Service) handlePullRequest(info *GitHubURLInfo) (string, error) {
	repoPath := s.config.GetRepoPath(info.Owner, info.Repo)

	if _, err := os.Stat(repoPath); err == nil {
		log.Printf("Repository exists, fetching updates...")
		if err := s.gitClient.Fetch(repoPath); err != nil {
			log.Printf("Warning: git fetch failed: %v", err)
		}
	} else {
		log.Printf("Cloning repository...")
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", info.Owner, info.Repo)
		if err := s.gitClient.Clone(repoURL, repoPath); err != nil {
			return "", fmt.Errorf("failed to clone: %w", err)
		}
	}

	log.Printf("Fetching PR #%d branch...", info.PRNumber)
	prBranchName := fmt.Sprintf("pr-%d", info.PRNumber)
	if err := s.gitClient.FetchPR(repoPath, info.PRNumber, prBranchName); err != nil {
		return "", fmt.Errorf("failed to fetch PR: %w", err)
	}

	log.Printf("Checking out PR branch: %s", prBranchName)
	if err := s.gitClient.Checkout(repoPath, prBranchName); err != nil {
		return "", fmt.Errorf("failed to checkout PR branch: %w", err)
	}

	return repoPath, nil
}
