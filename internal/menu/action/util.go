package action

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bodycam-save-tool/internal/i18n"
)

// safeJoin 安全拼接目标路径, 防止 Zip Slip
// Args:
//
//	root: 根路径
//	name: 目标路径
//
// Returns:
//
//	string: 完全路径
//	bool: 是否成功拼接
//	如果失败, 则返回空字符串和 false
//	如果成功, 则返回完全路径和 true
func safeJoin(root, name string) (string, bool) {
	if name == "" {
		return "", false
	}
	// 拼接路径
	target := filepath.Join(root, name)
	// 检查路径是否越界
	cleanRoot := filepath.Clean(root)
	// 只允许目标在 root 之下, 防止 Zip Slip
	if target != cleanRoot && !strings.HasPrefix(target, cleanRoot+string(os.PathSeparator)) {
		return "", false
	}
	return target, true
}

// addFileToZip 把单个文件或目录写入 zip 文件
// Args:
//
//	w: zip.Writer 实例
//	srcPath: 源文件或目录路径
//	relPath: 相对路径
//	info: 文件或目录信息
//
// Returns:
//
//	error: 如果写入失败, 则返回错误
//	如果成功, 则返回 nil
func addFileToZip(w *zip.Writer, srcPath, relPath string, info os.FileInfo) error {
	// 检查路径是否是文件夹
	if info.IsDir() {
		_, err := w.Create(relPath + "/")
		return err
	}
	// 创建 zip 头
	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// 设置名字
	hdr.Name = relPath
	// 设置压缩方法
	hdr.Method = zip.Deflate
	// 创建 zip 写入器
	fw, err := w.CreateHeader(hdr)
	if err != nil {
		return err
	}
	// 打开文件
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	// 最后关闭文件
	defer in.Close()
	// 复制文件
	_, err = io.Copy(fw, in)

	return err
}

// extractZipFile 把单个 zip 条目解压到 dstPath
// 先写到 dstPath + ".tmp", 成功后再重命名覆盖, 避免中途失败留下半个文件
// Args:
//
//	f: zip.File 实例
//	dstPath: 目标路径
//
// Returns:
//
//	error: 如果解压失败, 则返回错误
//	如果成功, 则返回 nil
func extractZipFile(f *zip.File, dstPath string) error {
	// 检查目标路径是否存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// 打开 zip 文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	// 最后关闭 zip 文件
	defer rc.Close()

	// 写临时文件
	tmpPath := dstPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	// 出错时清理临时文件
	defer func() {
		if err != nil {
			_ = out.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	// 复制文件
	if _, err = io.Copy(out, rc); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}

	// 重命名覆盖
	err = os.Rename(tmpPath, dstPath)
	if err == nil {
		return nil
	}

	// Windows 上 rename 不能覆盖已存在的文件, 先删目标再 rename
	removeErr := os.Remove(dstPath)
	if removeErr == nil {
		err = os.Rename(tmpPath, dstPath)
	}
	return err
}

// findConflicts 扫描所有会与目标目录冲突的文件, 返回完整路径列表
// 只检查文件, 不检查目录
// Args:
//
//	root: 根路径
//	files: zip.File 实例列表
//
// Returns:
//
//	[]string: 冲突文件路径列表
func findConflicts(root string, files []*zip.File) []string {
	var conflicts []string
	for _, f := range files {
		if f.FileInfo().IsDir() {
			continue
		}
		target, ok := safeJoin(root, f.Name)
		if !ok {
			continue
		}
		if _, err := os.Stat(target); err == nil {
			conflicts = append(conflicts, target)
		}
	}
	return conflicts
}

// extractAll 把 zip 中所有条目解压到 root
// Args:
//
//	root: 根路径
//	files: zip.File 实例列表
//
// Returns:
//
//	count: 成功解压的文件数
//	failed: 失败的文件列表, 每项格式为 "<路径>: <原因>"
func extractAll(root string, files []*zip.File) (int, []string) {
	count := 0
	var failed []string
	for _, f := range files {
		target, ok := safeJoin(root, f.Name)
		if !ok {
			failed = append(failed, f.Name+": "+i18n.T("恢复.路径越界"))
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				failed = append(failed, target+": "+err.Error())
			}
			continue
		}
		if err := extractZipFile(f, target); err != nil {
			failed = append(failed, target+": "+err.Error())
			continue
		}
		count++
	}
	return count, failed
}
