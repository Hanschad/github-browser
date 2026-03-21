package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	nativeHostMode := flag.Bool("native-host", false, "run as a Chrome Native Messaging host")
	flag.Parse()

	service, err := NewService()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	if *nativeHostMode || launchedByBrowser(flag.Args()) {
		if err := RunNativeHost(service); err != nil {
			log.Fatalf("Native host failed: %v", err)
		}
		return
	}

	runHTTPServer(service)
}

func runHTTPServer(service *Service) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", service.handleHealth)
	r.POST("/open", service.handleOpen)
	r.GET("/cache", service.handleListCache)
	r.DELETE("/cache/:repo", service.handleDeleteCache)
	r.GET("/config", service.handleGetConfig)
	r.PUT("/config", service.handleUpdateConfig)

	port := service.config.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("GitHub Browser service started on http://localhost:%d", port)
	log.Printf("Cache directory: %s", service.cacheDir)
	log.Printf("Default IDE: %s", service.config.DefaultIDE)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func launchedByBrowser(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "chrome-extension://") || strings.HasPrefix(arg, "edge-extension://") {
			return true
		}
	}

	return false
}

func (s *Service) handleHealth(c *gin.Context) {
	c.JSON(200, s.Health("http"))
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

	response, err := s.Open(req)
	if err != nil {
		c.JSON(500, OpenResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, response)
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

	if err := s.UpdateConfig(&newConfig); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "config": s.config})
}
