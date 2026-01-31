package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PathMapping 定义 GitHub 路径到本地目录的映射
// Pattern 支持：
//   - "owner" - 匹配特定用户/组织的所有仓库
//   - "owner/repo" - 匹配特定仓库
//   - "*" - 默认匹配所有
type PathMapping struct {
	Pattern   string `json:"pattern"`   // GitHub 路径模式，如 "microsoft" 或 "microsoft/vscode"
	LocalPath string `json:"localPath"` // 本地目录路径
}

type Config struct {
	Port         int           `json:"port"`
	DefaultIDE   string        `json:"defaultIDE"`
	GitHubToken  string        `json:"githubToken"`
	CacheDir     string        `json:"cacheDir"`
	PathMappings []PathMapping `json:"pathMappings,omitempty"` // 路径映射规则
}

func DefaultConfig() *Config {
	return &Config{
		Port:       DefaultPort,
		DefaultIDE: "code", // VS Code
		CacheDir:   filepath.Join(os.Getenv("HOME"), DefaultCacheDir),
	}
}

func LoadConfig() (*Config, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".github-browser", "config.json")

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".github-browser")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetRepoPath 根据 owner 和 repo 返回本地仓库路径
// 按优先级匹配：owner/repo > owner > * > 默认 cacheDir
func (c *Config) GetRepoPath(owner, repo string) string {
	exactMatch := owner + "/" + repo

	// 优先匹配 owner/repo
	for _, m := range c.PathMappings {
		if m.Pattern == exactMatch {
			return expandPath(m.LocalPath)
		}
	}

	// 其次匹配 owner
	for _, m := range c.PathMappings {
		if m.Pattern == owner {
			return filepath.Join(expandPath(m.LocalPath), repo)
		}
	}

	// 匹配通配符 *
	for _, m := range c.PathMappings {
		if m.Pattern == "*" {
			return filepath.Join(expandPath(m.LocalPath), owner+"-"+repo)
		}
	}

	// 默认使用 cacheDir
	cacheDir := c.CacheDir
	if cacheDir == "" {
		cacheDir = filepath.Join(os.Getenv("HOME"), DefaultCacheDir)
	}
	return filepath.Join(cacheDir, owner+"-"+repo)
}

// expandPath 展开路径中的 ~ 为 home 目录
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		return filepath.Join(os.Getenv("HOME"), path[1:])
	}
	return path
}
