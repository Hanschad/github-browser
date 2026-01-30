package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port        int    `json:"port"`
	DefaultIDE  string `json:"defaultIDE"`
	GitHubToken string `json:"githubToken"`
	CacheDir    string `json:"cacheDir"`
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
