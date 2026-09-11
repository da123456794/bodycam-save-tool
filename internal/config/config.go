package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// DefaultGamePath 默认游戏本体路径 (正斜杠形式)
const DefaultGamePath = `D:/SteamLibrary/steamapps/common/bodycam`

type Config struct {
	Language string `json:"language"`
	GamePath string `json:"game_path"` // 游戏本体路径
	DataPath string `json:"data_path"` // 游戏数据文件夹路径, 空表示使用默认
}

var current Config

// ExeDir 返回可执行文件所在目录
func ExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	return filepath.Dir(exe)
}

func filePath() string { return filepath.Join(ExeDir(), "config.json") }

// ToSlashPath 归一化路径
//   - 去掉首尾空格
//   - 去掉首尾英文引号
//   - 反斜杠统一转换成正斜杠
func ToSlashPath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, `"`)
	p = strings.ReplaceAll(p, `\`, `/`)
	return p
}

// defaultDataDir 返回系统默认的游戏数据目录
func defaultDataDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = "."
	}
	return ToSlashPath(filepath.Join(base, "Bodycam"))
}

// DefaultDataDir 返回系统默认的游戏数据目录, 供界面显示
func DefaultDataDir() string {
	return defaultDataDir()
}

// DataDir 返回游戏数据文件夹路径
// 若用户自定义过, 返回自定义路径; 否则返回系统默认路径
func DataDir() string {
	if strings.TrimSpace(current.DataPath) != "" {
		return current.DataPath
	}
	return defaultDataDir()
}

// GameDir 返回游戏本体路径
func GameDir() string {
	return current.GamePath
}

// Load 加载配置
func Load() Config {
	current = Config{
		Language: "zh",
		GamePath: ToSlashPath(DefaultGamePath),
		DataPath: "",
	}
	if data, err := os.ReadFile(filePath()); err == nil {
		_ = json.Unmarshal(data, &current)
	}
	// 归一化从配置读出的路径
	current.GamePath = ToSlashPath(current.GamePath)
	current.DataPath = ToSlashPath(current.DataPath)

	if current.Language != "zh" && current.Language != "en" {
		current.Language = "zh"
	}
	if current.GamePath == "" {
		current.GamePath = ToSlashPath(DefaultGamePath)
	}
	return current
}

// Get 返回当前配置
func Get() Config { return current }

// SetLanguage 设置语言
func SetLanguage(lang string) {
	current.Language = lang
	save()
}

// SetGamePath 设置游戏本体路径
func SetGamePath(p string) {
	current.GamePath = ToSlashPath(p)
	save()
}

// SetDataPath 设置游戏数据文件夹路径
// 传空字符串表示恢复默认
func SetDataPath(p string) {
	current.DataPath = ToSlashPath(p)
	save()
}

// Reset 恢复全部默认设置并写回配置
// Returns: 恢复后的配置
func Reset() Config {
	current = Config{
		Language: "zh",
		GamePath: ToSlashPath(DefaultGamePath),
		DataPath: "",
	}
	save()
	return current
}

// save 保存配置
func save() {
	// 序列化配置
	data, _ := json.MarshalIndent(current, "", "  ")
	_ = os.WriteFile(filePath(), data, 0644)
}
