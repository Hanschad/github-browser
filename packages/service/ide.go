package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type ideConfig struct {
	cmd      string
	args     []string
	gotoFlag string
}

var ides = map[string]ideConfig{
	"code":          {cmd: "code", args: []string{"$PATH"}, gotoFlag: "--goto"},
	"vscode":        {cmd: "code", args: []string{"$PATH"}, gotoFlag: "--goto"},
	"code-insiders": {cmd: "code-insiders", args: []string{"$PATH"}, gotoFlag: "--goto"},
	"zed":           {cmd: "zed", args: []string{"$PATH:$LINE"}},
	"cursor":        {cmd: "cursor", args: []string{"$PATH"}, gotoFlag: "--goto"},
	"idea":          {cmd: "idea", args: []string{"--line", "$LINE", "$PATH"}},
	"pycharm":       {cmd: "pycharm", args: []string{"--line", "$LINE", "$PATH"}},
	"webstorm":      {cmd: "webstorm", args: []string{"--line", "$LINE", "$PATH"}},
	"goland":        {cmd: "goland", args: []string{"--line", "$LINE", "$PATH"}},
	"subl":          {cmd: "subl", args: []string{"$PATH:$LINE"}},
	"sublime":       {cmd: "subl", args: []string{"$PATH:$LINE"}},
	"nvim":          {cmd: "nvim", args: []string{"+$LINE", "$PATH"}},
	"neovim":        {cmd: "nvim", args: []string{"+$LINE", "$PATH"}},
}

// OpenInIDE 在指定的 IDE 中打开文件或目录
func OpenInIDE(ideName, path string, line int) error {
	config, ok := ides[ideName]
	if !ok {
		return fmt.Errorf("unsupported IDE: %s", ideName)
	}

	// 特殊处理 Neovim 在 macOS 下的情况
	if (ideName == "nvim" || ideName == "neovim") && runtime.GOOS == "darwin" {
		cmdStr := fmt.Sprintf("nvim %s", path)
		if line > 0 {
			cmdStr = fmt.Sprintf("nvim +%d %s", line, path)
		}
		return exec.Command("osascript", "-e",
			fmt.Sprintf(`tell application "Terminal" to do script "%s"`, cmdStr)).Start()
	}

	var args []string
	if line > 0 {
		if config.gotoFlag != "" {
			args = append(args, config.gotoFlag, fmt.Sprintf("%s:%d", path, line))
		} else {
			for _, arg := range config.args {
				arg = strings.ReplaceAll(arg, "$PATH", path)
				arg = strings.ReplaceAll(arg, "$LINE", fmt.Sprintf("%d", line))
				args = append(args, arg)
			}
		}
	} else {
		// 不带行号时，移除所有 $LINE 相关的参数
		for _, arg := range config.args {
			if strings.Contains(arg, "$LINE") {
				// 如果参数只包含 $LINE 或者是以 $LINE 开头的标志，跳过它
				if arg == "$LINE" || strings.HasPrefix(arg, "--line") {
					continue
				}
				// 否则（如 $PATH:$LINE），只保留 $PATH
				arg = strings.ReplaceAll(arg, ":$LINE", "")
				arg = strings.ReplaceAll(arg, "+$LINE", "")
			}
			arg = strings.ReplaceAll(arg, "$PATH", path)
			args = append(args, arg)
		}
	}

	cmd := exec.Command(config.cmd, args...)
	return cmd.Start()
}
