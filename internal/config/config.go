package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// DefaultLanguage 默认语言
const DefaultLanguage = "zh"

// DefaultGamePath 默认游戏本体路径
const DefaultGamePath = `D:/SteamLibrary/steamapps/common/bodycam`

// 默认备份文件名
const (
	DefaultDataBackupName          = "bodycam_data.zip"
	DefaultSaveGamesBackupName     = "bodycam_savegames.zip"
	DefaultWindowsConfigBackupName = "bodycam_windows_config.zip"
)

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

	// 数据备份文件名, 空表示使用默认
	DataBackupName string `json:"data_backup_name"`
	// 存档备份文件名, 空表示使用默认
	SaveGamesBackupName string `json:"save_games_backup_name"`
	// 配置备份文件名, 空表示使用默认
	WindowsConfigBackupName string `json:"windows_config_backup_name"`
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

// sanitizeFileName 只保留文件名部分, 去掉任何目录分隔
func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = filepath.Base(name)
	if name == "." || name == ".." {
		return ""
	}
	return name
}

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

// DataBackupName 返回数据备份文件名
func DataBackupName() string {
	if strings.TrimSpace(current.DataBackupName) != "" {
		return current.DataBackupName
	}
	return DefaultDataBackupName
}

// SaveGamesBackupName 返回存档备份文件名
func SaveGamesBackupName() string {
	if strings.TrimSpace(current.SaveGamesBackupName) != "" {
		return current.SaveGamesBackupName
	}
	return DefaultSaveGamesBackupName
}

// WindowsConfigBackupName 返回配置备份文件名
func WindowsConfigBackupName() string {
	if strings.TrimSpace(current.WindowsConfigBackupName) != "" {
		return current.WindowsConfigBackupName
	}
	return DefaultWindowsConfigBackupName
}

// Load 加载配置
func Load() Config {
	current = Config{
		Language:                DefaultLanguage,
		GamePath:                ToSlashPath(DefaultGamePath),
		DataPath:                "",
		SaveGamesPath:           "",
		WindowsConfigPath:       "",
		DataBackupName:          "",
		SaveGamesBackupName:     "",
		WindowsConfigBackupName: "",
	}
	if data, err := os.ReadFile(filePath()); err == nil {
		_ = json.Unmarshal(data, &current)
	}
	// 统一路径格式
	current.GamePath = ToSlashPath(current.GamePath)
	current.DataPath = ToSlashPath(current.DataPath)
	current.SaveGamesPath = ToSlashPath(current.SaveGamesPath)
	current.WindowsConfigPath = ToSlashPath(current.WindowsConfigPath)

	// 备份文件名做安全过滤
	current.DataBackupName = sanitizeFileName(current.DataBackupName)
	current.SaveGamesBackupName = sanitizeFileName(current.SaveGamesBackupName)
	current.WindowsConfigBackupName = sanitizeFileName(current.WindowsConfigBackupName)

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

// SetDataBackupName 设置数据备份文件名, 传空字符串恢复默认
func SetDataBackupName(name string) {
	current.DataBackupName = sanitizeFileName(name)
	save()
}

// SetSaveGamesBackupName 设置存档备份文件名, 传空字符串恢复默认
func SetSaveGamesBackupName(name string) {
	current.SaveGamesBackupName = sanitizeFileName(name)
	save()
}

// SetWindowsConfigBackupName 设置配置备份文件名, 传空字符串恢复默认
func SetWindowsConfigBackupName(name string) {
	current.WindowsConfigBackupName = sanitizeFileName(name)
	save()
}

// Reset 恢复全部默认设置
// Returns:
//
//	Config: 新的配置
func Reset() Config {
	current = Config{
		Language:                DefaultLanguage,
		GamePath:                ToSlashPath(DefaultGamePath),
		DataPath:                "",
		SaveGamesPath:           "",
		WindowsConfigPath:       "",
		DataBackupName:          "",
		SaveGamesBackupName:     "",
		WindowsConfigBackupName: "",
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
