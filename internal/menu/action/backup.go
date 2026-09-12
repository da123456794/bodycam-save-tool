package action

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
)

// Backup 把数据目录全部内容压缩到 exe 同目录
func Backup() {
	src := config.DataDir()
	fi, err := os.Stat(src)
	if err != nil || !fi.IsDir() {
		fmt.Println(i18n.Tf("备份.存档目录不存在", src))
		return
	}

	// 创建压缩包
	out := ZipPath()
	f, err := os.Create(out)
	if err != nil {
		fmt.Println(i18n.Tf("备份.创建压缩包失败", err))
		return
	}

	// 创建 zip 写入器
	w := zip.NewWriter(f)
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
		_ = f.Close()
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}
	// 如果关闭错误
	err = w.Close()
	if err != nil {
		_ = f.Close()
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}
	// 如果关闭错误
	err = f.Close()
	if err != nil {
		fmt.Println(i18n.Tf("备份.压缩失败", err))
		return
	}

	fmt.Println(i18n.T("备份.完成"))
	fmt.Println(i18n.Tf("备份.来源", src))
	fmt.Println(i18n.Tf("备份.文件", out, count))

	// 有跳过的文件, 汇总提示
	if skipped > 0 {
		fmt.Println(i18n.Tf("备份.跳过汇总", skipped))
	}
}
