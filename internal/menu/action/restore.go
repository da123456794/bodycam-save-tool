package action

import (
	"archive/zip"
	"fmt"
	"os"

	"bodycam-save-tool/internal/i18n"
)

// Restore 从 zipPath 解压到 dst
// Args:
// zipPath: 压缩包路径
// dst: 目标路径
// overwrite: 冲突时是否全部覆盖(true 表示有冲突也全部覆盖, false 表示有冲突时什么都不做)
func Restore(zipPath, dst string, overwrite bool) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println(i18n.Tf("恢复.打不开压缩包", zipPath))
		return
	}

	// 关闭 zip 包
	defer r.Close()

	if err := os.MkdirAll(dst, 0755); err != nil {
		fmt.Println(i18n.Tf("恢复.创建目录失败", err))
		return
	}

	// 冲突预检
	// 只有不允许覆盖时, 才需要检查冲突
	if !overwrite {
		conflicts := findConflicts(dst, r.File)
		if len(conflicts) > 0 {
			fmt.Println(i18n.T("恢复.已取消"))
			return
		}
	}

	// 执行解压
	count, failed := extractAll(dst, r.File)

	fmt.Println(i18n.T("恢复.完成"))
	fmt.Println(i18n.Tf("恢复.文件", zipPath))
	fmt.Println(i18n.Tf("恢复.目标", dst, count))

	// 有失败的文件, 列出全部
	if len(failed) > 0 {
		fmt.Println(i18n.Tf("恢复.失败列表", len(failed)))
		for _, s := range failed {
			fmt.Println("  " + s)
		}
	}

}

// HasConflict 检查 zipPath 解压到 dst 时是否会有文件冲突
// Args:
//
//	zipPath: 压缩包路径
//	dst: 目标路径
//
// Returns:
//
//	bool: 是否会有文件冲突
//	[]string: 冲突的文件列表
func HasConflict(zipPath, dst string) (bool, []string) {

	// 打开 zip 包
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return false, nil
	}
	defer r.Close()

	conflicts := findConflicts(dst, r.File)
	if len(conflicts) == 0 {
		return false, nil
	}
	return true, conflicts
}
