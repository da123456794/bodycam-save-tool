package action

import (
	"path/filepath"

	"bodycam-save-tool/internal/config"
)

// DataBackupPath 返回数据备份文件的完整路径
func DataBackupPath() string {
	return filepath.Join(config.ExeDir(), config.DataBackupName())
}

// SaveGamesBackupPath 返回存档备份文件的完整路径
func SaveGamesBackupPath() string {
	return filepath.Join(config.ExeDir(), config.SaveGamesBackupName())
}

// WindowsConfigBackupPath 返回配置备份文件的完整路径
func WindowsConfigBackupPath() string {
	return filepath.Join(config.ExeDir(), config.WindowsConfigBackupName())
}
