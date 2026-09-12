package pathutil

import (
	"os"
	"strings"
)

// IsReachable 检测路径是否可达
// 路径为空, 不存在, 或无权限访问时返回 false
func IsReachable(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// IsEmptyDir 检测路径是否为空目录
// 路径不存在, 不是目录, 或无权限读取时返回 false
func IsEmptyDir(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil || !fi.IsDir() {
		return false
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
