package action

import (
	"path/filepath"

	"bodycam-save-tool/internal/config"
)

// backupZipName 备份文件名
const backupZipName = "bodycam_backup.zip"

// ZipPath 备份 zip 的完整路径
func ZipPath() string {
	return filepath.Join(config.ExeDir(), backupZipName)
}
