package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// DefaultGamePath 默认游戏本体路径
const DefaultGamePath = `D:/SteamLibrary/steamapps/common/bodycam`

// DefaultLanguage 默认语言
const DefaultLanguage = "zh"

type Config struct {
	Language string `json:"language"`
	// 游戏本体路径
	GamePath string `json:"game_path"`
	// 游戏数据文件夹路径, 空表示使用默认
	DataPath string `json:"data_path"`
	// 游戏存档文件夹路径, 空表示自数据文件夹
	SaveGamesPath string `json:"save_games_path"`
	// 游戏配置文件夹路径, 空表示自数据文件夹
	WindowsConfigPath string `json:"windows_config_path"`
}

var current Config

// ExeDir 返回可执行文件所在目录
// Returns:
//
//	string: 可执行文件所在目录
func ExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	return filepath.Dir(exe)
}

// filePath 返回配置文件路径
// Returns:
//
//	string: 配置文件路径
func filePath() string {
	return filepath.Join(ExeDir(), "config.json")
}

// ToSlashPath 归一化路径
//   - 去掉首尾空格
//   - 去掉首尾英文引号
//   - 反斜杠统一转换成正斜杠
//
// Args:
//
//	p: 原始路径
//
// Returns:
//
//	string: 归一化后的路径
func ToSlashPath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, `"`)
	p = strings.ReplaceAll(p, `\`, `/`)
	return p
}

// defaultDataDir 返回默认的游戏数据目录
// Returns:
//
//	string: 默认的游戏数据目录
func defaultDataDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = "."
	}
	return ToSlashPath(filepath.Join(base, "Bodycam"))
}

// DataDir 返回游戏数据文件夹路径
// 若用户自定义过, 返回自定义路径; 否则返回默认路径
// Returns:
//
//	string: 游戏数据文件夹路径
func DataDir() string {
	// 若用户自定义过, 返回自定义路径
	if strings.TrimSpace(current.DataPath) != "" {
		return current.DataPath
	}
	return defaultDataDir()
}

// defaultSaveGamesDir 返回的游戏存档文件夹路径
func defaultSaveGamesDir() string {
	return ToSlashPath(filepath.Join(DataDir(), "Saved", "SaveGames"))
}

// defaultWindowsConfigDir 返回的游戏配置文件夹路径
func defaultWindowsConfigDir() string {
	return ToSlashPath(filepath.Join(DataDir(), "Saved", "Config", "Windows"))
}

// SaveGamesDir 返回游戏存档文件夹路径
// 若用户自定义过, 返回自定义路径; 否则返回路径
func SaveGamesDir() string {
	if strings.TrimSpace(current.SaveGamesPath) != "" {
		return current.SaveGamesPath
	}
	return defaultSaveGamesDir()
}

// WindowsConfigDir 返回游戏配置文件夹路径
// 若用户自定义过, 返回自定义路径; 否则返回默认路径
func WindowsConfigDir() string {
	if strings.TrimSpace(current.WindowsConfigPath) != "" {
		return current.WindowsConfigPath
	}
	return defaultWindowsConfigDir()
}

// GameDir 返回游戏本体路径
// Returns:
//
//	string: 游戏本体路径
func GameDir() string {
	return current.GamePath
}

// Load 加载配置
func Load() Config {
	current = Config{
		Language:          DefaultLanguage,
		GamePath:          ToSlashPath(DefaultGamePath),
		DataPath:          "",
		SaveGamesPath:     "",
		WindowsConfigPath: "",
	}
	if data, err := os.ReadFile(filePath()); err == nil {
		_ = json.Unmarshal(data, &current)
	}
	// 统一路径格式
	current.GamePath = ToSlashPath(current.GamePath)
	current.DataPath = ToSlashPath(current.DataPath)
	current.SaveGamesPath = ToSlashPath(current.SaveGamesPath)
	current.WindowsConfigPath = ToSlashPath(current.WindowsConfigPath)

	// 语言只做非空校验
	if strings.TrimSpace(current.Language) == "" {
		current.Language = DefaultLanguage
	}
	// 游戏本体路径做非空校验
	if current.GamePath == "" {
		current.GamePath = ToSlashPath(DefaultGamePath)
	}

	return current
}

// Get 返回当前配置
func Get() Config {
	return current
}

// SetLanguage 设置语言
// Args:
//
//	lang: 语言代码
func SetLanguage(lang string) {
	current.Language = lang
	save()
}

// SetGamePath 设置游戏本体路径
// Args:
//
//	p: 游戏本体路径
func SetGamePath(p string) {
	current.GamePath = ToSlashPath(p)
	save()
}

// SetDataPath 设置游戏数据文件夹路径
// 传空字符串表示恢复默认路径
// Args:
//
//	p: 游戏数据文件夹路径
func SetDataPath(p string) {
	current.DataPath = ToSlashPath(p)
	save()
}

// SetSaveGamesPath 设置游戏存档文件夹路径
// 传空字符串表示恢复默认
func SetSaveGamesPath(p string) {
	current.SaveGamesPath = ToSlashPath(p)
	save()
}

// SetWindowsConfigPath 设置游戏配置文件夹路径
// 传空字符串表示恢复默认
func SetWindowsConfigPath(p string) {
	current.WindowsConfigPath = ToSlashPath(p)
	save()
}

// Reset 恢复全部默认设置
// Returns:
//
//	Config: 新的配置
func Reset() Config {
	current = Config{
		Language:          DefaultLanguage,
		GamePath:          ToSlashPath(DefaultGamePath),
		DataPath:          "",
		SaveGamesPath:     "",
		WindowsConfigPath: "",
	}
	save()
	return current
}

// save 保存配置
// Args:
//
//	config: 要保存的配置
func save() {
	data, _ := json.MarshalIndent(current, "", "  ")
	_ = os.WriteFile(filePath(), data, 0644)
}
