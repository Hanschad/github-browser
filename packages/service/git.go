package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitClient struct {
	cacheDir string
}

func NewGitClient(cacheDir string) *GitClient {
	return &GitClient{cacheDir: cacheDir}
}

// Clone 克隆仓库
func (gc *GitClient) Clone(repoURL, targetPath string) error {
	// 不使用 --single-branch，以便可以 checkout 其他分支
	cmd := exec.Command("git", "clone",
		"--filter=blob:none",
		repoURL,
		targetPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// Pull 更新仓库
func (gc *GitClient) Pull(repoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// Checkout 切换分支
func (gc *GitClient) Checkout(repoPath, branch string) error {
	// 先尝试直接 checkout
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果失败，可能是远程分支，尝试创建本地分支
		cmd = exec.Command("git", "checkout", "-b", branch, fmt.Sprintf("origin/%s", branch))
		cmd.Dir = repoPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git checkout failed: %v\nOutput: %s", err, string(output))
		}
	}
	return nil
}

// AddRemote 添加远程仓库
func (gc *GitClient) AddRemote(repoPath, name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 检查是否是因为 remote 已存在
		if strings.Contains(string(output), "already exists") {
			return nil
		}
		return fmt.Errorf("git remote add failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// FetchBranch 获取远程分支
func (gc *GitClient) FetchBranch(repoPath, remote, branch string) error {
	cmd := exec.Command("git", "fetch", remote, branch)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// CheckoutRemoteBranch 从远程分支创建并切换到本地分支
func (gc *GitClient) CheckoutRemoteBranch(repoPath, remote, remoteBranch, localBranch string) error {
	// 先删除可能存在的本地分支
	cmd := exec.Command("git", "branch", "-D", localBranch)
	cmd.Dir = repoPath
	cmd.Run() // 忽略错误

	// 创建并切换到新分支
	cmd = exec.Command("git", "checkout", "-b", localBranch, fmt.Sprintf("%s/%s", remote, remoteBranch))
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout remote branch failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// GetCurrentBranch 获取当前分支名
func (gc *GitClient) GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch failed: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// Fetch 获取所有远程更新
func (gc *GitClient) Fetch(repoPath string) error {
	cmd := exec.Command("git", "fetch", "--all", "--tags")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// CheckoutTag 切换到指定 tag
func (gc *GitClient) CheckoutTag(repoPath, tag string) error {
	cmd := exec.Command("git", "checkout", fmt.Sprintf("tags/%s", tag))
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout tag failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// FetchPR 获取 PR 分支（使用 GitHub 的 refs/pull/<number>/head 格式）
func (gc *GitClient) FetchPR(repoPath string, prNumber int, localBranch string) error {
	// 使用 git fetch origin pull/<PR_NUMBER>/head:<local_branch>
	refspec := fmt.Sprintf("pull/%d/head:%s", prNumber, localBranch)
	cmd := exec.Command("git", "fetch", "origin", refspec)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch PR failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}
