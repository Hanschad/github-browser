package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPort     = 9527
	DefaultCacheDir = ".github-browser/repos"
)

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

func main() {
	// 初始化配置
	config, err := LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v, using defaults", err)
		config = DefaultConfig()
	}

	// 创建缓存目录
	cacheDir := filepath.Join(os.Getenv("HOME"), DefaultCacheDir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatalf("Failed to create cache directory: %v", err)
	}

	// 初始化服务
	service := &Service{
		config:    config,
		cacheDir:  cacheDir,
		gitClient: NewGitClient(cacheDir),
		ghClient:  NewGitHubClient(config.GitHubToken),
	}

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 路由
	r.GET("/health", service.handleHealth)
	r.POST("/open", service.handleOpen)
	r.GET("/cache", service.handleListCache)
	r.DELETE("/cache/:repo", service.handleDeleteCache)
	r.GET("/config", service.handleGetConfig)
	r.PUT("/config", service.handleUpdateConfig)

	// 启动服务
	port := config.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("🚀 GitHub Browser service started on http://localhost:%d", port)
	log.Printf("📁 Cache directory: %s", cacheDir)
	log.Printf("💻 Default IDE: %s", config.DefaultIDE)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Service) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"uptime":  time.Since(time.Now()).String(),
	})
}

func (s *Service) handleOpen(c *gin.Context) {
	var req OpenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, OpenResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	log.Printf("📥 Received request: %s", req.URL)

	// 解析 URL
	info, err := ParseGitHubURL(req.URL)
	if err != nil {
		c.JSON(400, OpenResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid GitHub URL: %v", err),
		})
		return
	}

	log.Printf("📦 Parsed: owner=%s, repo=%s, type=%s", info.Owner, info.Repo, info.Type)

	// 处理不同类型
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
		c.JSON(500, OpenResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// 确定要打开的文件路径
	var targetPath string
	if req.FilePath != "" {
		targetPath = filepath.Join(repoPath, req.FilePath)
	} else if info.FilePath != "" {
		targetPath = filepath.Join(repoPath, info.FilePath)
	} else {
		targetPath = repoPath
	}

	// 确定行号
	line := req.Line
	if line == 0 && info.Line > 0 {
		line = info.Line
	}

	// 确定 IDE
	ide := req.IDE
	if ide == "" {
		ide = s.config.DefaultIDE
	}

	// 打开 IDE
	log.Printf("🚀 Opening in %s: %s (line: %d)", ide, targetPath, line)
	if err := OpenInIDE(ide, targetPath, line); err != nil {
		c.JSON(500, OpenResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to open IDE: %v", err),
		})
		return
	}

	c.JSON(200, OpenResponse{
		Status:  "ok",
		Message: "Opened successfully",
		Path:    repoPath,
	})
}

func (s *Service) handleRepository(info *GitHubURLInfo) (string, error) {
	repoPath := filepath.Join(s.cacheDir, fmt.Sprintf("%s-%s", info.Owner, info.Repo))

	// 克隆或更新
	if _, err := os.Stat(repoPath); err == nil {
		log.Printf("📦 Repository exists, updating...")
		if err := s.gitClient.Pull(repoPath); err != nil {
			log.Printf("⚠️  Warning: git pull failed: %v", err)
		}
	} else {
		log.Printf("📥 Cloning repository...")
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", info.Owner, info.Repo)
		if err := s.gitClient.Clone(repoURL, repoPath); err != nil {
			return "", fmt.Errorf("failed to clone: %v", err)
		}
	}

	// 如果指定了分支或 tag，切换到该分支/tag
	if info.Branch != "" {
		log.Printf("🔀 Checking out branch/tag: %s", info.Branch)
		// 先 fetch 确保有最新的远程分支
		if err := s.gitClient.Fetch(repoPath); err != nil {
			log.Printf("⚠️  Warning: git fetch failed: %v", err)
		}
		if err := s.gitClient.Checkout(repoPath, info.Branch); err != nil {
			// 可能是 tag，尝试 checkout tag
			log.Printf("⚠️  Branch checkout failed, trying as tag: %v", err)
			if err := s.gitClient.CheckoutTag(repoPath, info.Branch); err != nil {
				log.Printf("⚠️  Warning: failed to checkout branch/tag: %v", err)
			}
		}
	}

	return repoPath, nil
}

func (s *Service) handlePullRequest(info *GitHubURLInfo) (string, error) {
	repoPath := filepath.Join(s.cacheDir, fmt.Sprintf("%s-%s", info.Owner, info.Repo))

	// 克隆或更新主仓库
	if _, err := os.Stat(repoPath); err == nil {
		log.Printf("📦 Repository exists, fetching updates...")
		if err := s.gitClient.Fetch(repoPath); err != nil {
			log.Printf("⚠️  Warning: git fetch failed: %v", err)
		}
	} else {
		log.Printf("📥 Cloning repository...")
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", info.Owner, info.Repo)
		if err := s.gitClient.Clone(repoURL, repoPath); err != nil {
			return "", fmt.Errorf("failed to clone: %v", err)
		}
	}

	// 使用 git fetch 直接获取 PR 分支（无需 GitHub API）
	// GitHub 支持 refs/pull/<PR_NUMBER>/head 格式
	log.Printf("📥 Fetching PR #%d branch...", info.PRNumber)
	prBranchName := fmt.Sprintf("pr-%d", info.PRNumber)
	if err := s.gitClient.FetchPR(repoPath, info.PRNumber, prBranchName); err != nil {
		return "", fmt.Errorf("failed to fetch PR: %v", err)
	}

	log.Printf("🔀 Checking out PR branch: %s", prBranchName)
	if err := s.gitClient.Checkout(repoPath, prBranchName); err != nil {
		return "", fmt.Errorf("failed to checkout PR branch: %v", err)
	}

	return repoPath, nil
}

func (s *Service) handleListCache(c *gin.Context) {
	entries, err := os.ReadDir(s.cacheDir)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var repos []map[string]interface{}
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(s.cacheDir, entry.Name())
			info, _ := entry.Info()
			repos = append(repos, map[string]interface{}{
				"name":     entry.Name(),
				"path":     path,
				"modified": info.ModTime(),
			})
		}
	}

	c.JSON(200, gin.H{
		"repos": repos,
		"count": len(repos),
	})
}

func (s *Service) handleDeleteCache(c *gin.Context) {
	repo := c.Param("repo")
	repoPath := filepath.Join(s.cacheDir, repo)

	if err := os.RemoveAll(repoPath); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "message": "Cache deleted"})
}

func (s *Service) handleGetConfig(c *gin.Context) {
	c.JSON(200, s.config)
}

func (s *Service) handleUpdateConfig(c *gin.Context) {
	var newConfig Config
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 更新配置
	s.config = &newConfig

	// 保存到文件
	configPath := filepath.Join(os.Getenv("HOME"), ".github-browser", "config.json")
	data, _ := json.MarshalIndent(newConfig, "", "  ")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "config": s.config})
}
