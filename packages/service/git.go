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

// Checkout 切换分支或 tag
// ref 可能包含路径（如 "feature/develop/src/file.go"），需要智能解析分支名
func (gc *GitClient) Checkout(repoPath, ref string) error {
	// 1. 先 fetch 确保远程引用是最新的
	gc.Fetch(repoPath)

	// 2. 尝试直接 checkout 整个 ref (适用于简单分支名或 tag)
	if gc.tryCheckout(repoPath, ref) {
		return nil
	}

	// 3. 如果失败，可能是因为 ref 包含路径，尝试从长到短匹配分支名
	// 例如 "feature/develop/src/file" -> 尝试 "feature/develop/src/file", "feature/develop/src", "feature/develop", "feature"
	parts := strings.Split(ref, "/")
	for i := len(parts); i >= 1; i-- {
		candidate := strings.Join(parts[:i], "/")
		if gc.tryCheckout(repoPath, candidate) {
			return nil
		}
	}

	return fmt.Errorf("git checkout failed for ref %s", ref)
}

// tryCheckout 尝试 checkout 指定的 ref，成功返回 true
func (gc *GitClient) tryCheckout(repoPath, ref string) bool {
	// 1. 尝试直接 checkout (本地分支或 tag)
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = repoPath
	if _, err := cmd.CombinedOutput(); err == nil {
		return true
	}

	// 2. 尝试从远程分支创建本地分支
	cmd = exec.Command("git", "checkout", "-b", ref, fmt.Sprintf("origin/%s", ref))
	cmd.Dir = repoPath
	if _, err := cmd.CombinedOutput(); err == nil {
		return true
	}

	// 3. 尝试作为 tag 显式 checkout
	cmd = exec.Command("git", "checkout", fmt.Sprintf("tags/%s", ref))
	cmd.Dir = repoPath
	if _, err := cmd.CombinedOutput(); err == nil {
		return true
	}

	return false
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
