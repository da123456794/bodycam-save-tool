package action

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
)

// backupDir 把 src 目录全部内容压缩到 out 文件中
func backupDir(src, out string) {
	fi, err := os.Stat(src)
	if err != nil || !fi.IsDir() {
		fmt.Println(i18n.Tf("备份.源目录不存在", src))
		return
	}

	// 创建临时文件
	dir := filepath.Dir(out)
	base := filepath.Base(out)
	tmpFile, err := os.CreateTemp(dir, base+".*.tmp")
	if err != nil {
		fmt.Println(i18n.Tf("备份.创建压缩包失败", err))
		return
	}
	tmpPath := tmpFile.Name()

	// 归一化 out 与 tmpPath, 用于 Walk 时跳过自身
	cleanOut := cleanAbs(out)
	cleanTmp := cleanAbs(tmpPath)

	// 创建 zip 写入器
	w := zip.NewWriter(tmpFile)
	// 统计压缩的文件数量
	count := 0
	// 统计跳过的文件数量
	skipped := 0

	// 遍历数据目录下的所有文件
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		// 遇到没权限等错误, 跳过但继续
		if err != nil {
			fmt.Println(i18n.Tf("备份.跳过文件", path, err))
			// 加1
			skipped++
			return nil
		}
		if path == src {
			return nil
		}

		// 跳过输出文件与临时文件本身, 不计入 skipped
		abs := cleanAbs(path)
		if strings.EqualFold(abs, cleanOut) || strings.EqualFold(abs, cleanTmp) {
			return nil
		}

		// 计算相对路径
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		// 转换路径为斜杠分隔
		rel = filepath.ToSlash(rel)
		// 写入文件到 zip 包
		if err := addFileToZip(w, path, rel, info); err != nil {
			return err
		}
		// 统计文件数量
		if !info.IsDir() {
			count++
		}
		return nil
	})

	// 关闭 zip 写入器
	if err != nil {
		_ = w.Close()
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}
	// 如果关闭错误
	err = w.Close()
	if err != nil {
		_ = tmpFile.Close()
		// 关闭临时文件失败, 也直接返回
		_ = os.Remove(tmpPath)
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}
	// 如果关闭错误
	err = tmpFile.Close()
	if err != nil {
		// 关闭临时文件失败, 也直接返回
		_ = os.Remove(tmpPath)
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}

	bakPath := out + ".old"
	hasOld := false

	// 如果 out 存在, 先挪到 out + ".old"
	if _, statErr := os.Stat(out); statErr == nil {
		// 清理可能残留的旧 .old 文件
		if rmErr := os.Remove(bakPath); rmErr != nil && !os.IsNotExist(rmErr) {
			fmt.Println(i18n.Tf("备份.清理旧备份失败", bakPath, rmErr))
		}
		// 如果重命名错误
		err = os.Rename(out, bakPath)
		if err != nil {
			_ = os.Remove(tmpPath)
			fmt.Println(i18n.Tf("备份.压缩失败", err))
			return
		}
		hasOld = true
	}

	if err := os.Rename(tmpPath, out); err != nil {
		// 还原旧备份
		if hasOld {
			_ = os.Remove(out)
			_ = os.Rename(bakPath, out)
		}
		_ = os.Remove(tmpPath)
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}

	// 成功, 删除旧备份
	if hasOld {
		_ = os.Remove(bakPath)
	}

	fmt.Println(i18n.T("备份.完成"))
	fmt.Println(i18n.Tf("备份.来源", src))
	fmt.Println(i18n.Tf("备份.文件", out, count))

	// 有跳过的文件, 汇总提示
	if skipped > 0 {
		fmt.Println(i18n.Tf("备份.跳过汇总", skipped))
	}
}

// cleanAbs 返回路径的绝对路径并归一化
// Args:
//
//	p: 路径
//
// Returns:
//
//	string: 绝对路径
func cleanAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return filepath.Clean(p)
	}
	return filepath.Clean(abs)
}

// BackupData 备份游戏数据文件夹
func BackupData() {
	backupDir(config.DataDir(), DataBackupPath())
}

// BackupSaveGames 备份游戏存档文件夹
func BackupSaveGames() {
	backupDir(config.SaveGamesDir(), SaveGamesBackupPath())
}

// BackupWindowsConfig 备份游戏配置文件夹
func BackupWindowsConfig() {
	backupDir(config.WindowsConfigDir(), WindowsConfigBackupPath())
}
