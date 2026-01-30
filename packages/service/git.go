package main

import (
	"fmt"
	"os/exec"
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

// Checkout 切换分支或 tag
func (gc *GitClient) Checkout(repoPath, ref string) error {
	// 1. 尝试直接 checkout (适用于本地分支或 tag)
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = repoPath
	if _, err := cmd.CombinedOutput(); err == nil {
		return nil
	} else {
		// 2. 如果失败，尝试作为远程分支处理
		// git checkout -b <ref> origin/<ref>
		cmd = exec.Command("git", "checkout", "-b", ref, fmt.Sprintf("origin/%s", ref))
		cmd.Dir = repoPath
		if _, err := cmd.CombinedOutput(); err == nil {
			return nil
		}

		// 3. 如果还是失败，尝试作为 tag 显式 checkout
		// git checkout tags/<ref>
		cmd = exec.Command("git", "checkout", fmt.Sprintf("tags/%s", ref))
		cmd.Dir = repoPath
		if _, err := cmd.CombinedOutput(); err == nil {
			return nil
		} else {
			return fmt.Errorf("git checkout failed for ref %s", ref)
		}
	}
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
