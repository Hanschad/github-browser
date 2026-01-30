package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenInIDE 在指定的 IDE 中打开文件或目录
func OpenInIDE(ide, path string, line int) error {
	var cmd *exec.Cmd

	switch ide {
	case "code", "vscode":
		// VS Code
		if line > 0 {
			cmd = exec.Command("code", "--goto", fmt.Sprintf("%s:%d", path, line))
		} else {
			cmd = exec.Command("code", path)
		}

	case "code-insiders":
		// VS Code Insiders
		if line > 0 {
			cmd = exec.Command("code-insiders", "--goto", fmt.Sprintf("%s:%d", path, line))
		} else {
			cmd = exec.Command("code-insiders", path)
		}

	case "zed":
		// Zed
		if line > 0 {
			cmd = exec.Command("zed", fmt.Sprintf("%s:%d", path, line))
		} else {
			cmd = exec.Command("zed", path)
		}

	case "cursor":
		// Cursor
		if line > 0 {
			cmd = exec.Command("cursor", "--goto", fmt.Sprintf("%s:%d", path, line))
		} else {
			cmd = exec.Command("cursor", path)
		}

	case "idea":
		// IntelliJ IDEA
		if line > 0 {
			cmd = exec.Command("idea", "--line", fmt.Sprintf("%d", line), path)
		} else {
			cmd = exec.Command("idea", path)
		}

	case "pycharm":
		// PyCharm
		if line > 0 {
			cmd = exec.Command("pycharm", "--line", fmt.Sprintf("%d", line), path)
		} else {
			cmd = exec.Command("pycharm", path)
		}

	case "webstorm":
		// WebStorm
		if line > 0 {
			cmd = exec.Command("webstorm", "--line", fmt.Sprintf("%d", line), path)
		} else {
			cmd = exec.Command("webstorm", path)
		}

	case "goland":
		// GoLand
		if line > 0 {
			cmd = exec.Command("goland", "--line", fmt.Sprintf("%d", line), path)
		} else {
			cmd = exec.Command("goland", path)
		}

	case "nvim", "neovim":
		// Neovim
		if runtime.GOOS == "darwin" {
			// macOS: 在新终端窗口打开
			if line > 0 {
				cmd = exec.Command("osascript", "-e",
					fmt.Sprintf(`tell application "Terminal" to do script "nvim +%d %s"`, line, path))
			} else {
				cmd = exec.Command("osascript", "-e",
					fmt.Sprintf(`tell application "Terminal" to do script "nvim %s"`, path))
			}
		} else {
			// Linux: 直接启动（假设在终端中运行）
			if line > 0 {
				cmd = exec.Command("nvim", fmt.Sprintf("+%d", line), path)
			} else {
				cmd = exec.Command("nvim", path)
			}
		}

	case "subl", "sublime":
		// Sublime Text
		if line > 0 {
			cmd = exec.Command("subl", fmt.Sprintf("%s:%d", path, line))
		} else {
			cmd = exec.Command("subl", path)
		}

	default:
		return fmt.Errorf("unsupported IDE: %s", ide)
	}

	// 启动 IDE（不等待）
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %v", ide, err)
	}

	return nil
}
