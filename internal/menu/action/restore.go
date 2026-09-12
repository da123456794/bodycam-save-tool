package action

import (
	"archive/zip"
	"fmt"
	"os"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
)

// Restore 从 exe 同目录的 zip 解压回数据目录
// Args:
//
// overwrite: 冲突时是否全部覆盖(true 表示有冲突也全部覆盖, false 表示有冲突时什么都不做)
func Restore(overwrite bool) {
	// 备份 zip 的完整路径
	z := ZipPath()

	// 打开 zip 包
	r, err := zip.OpenReader(z)
	if err != nil {
		fmt.Println(i18n.Tf("恢复.打不开压缩包", z))
		return
	}

	// 关闭 zip 包
	defer r.Close()

	dst := config.DataDir()
	if err := os.MkdirAll(dst, 0755); err != nil {
		fmt.Println(i18n.Tf("恢复.创建目录失败", err))
		return
	}

	// 冲突预检
	conflicts := findConflicts(dst, r.File)
	if len(conflicts) > 0 && !overwrite {
		fmt.Println(i18n.T("恢复.已取消"))
		return
	}

	// 执行解压
	count := extractAll(dst, r.File)

	fmt.Println(i18n.T("恢复.完成"))
	fmt.Println(i18n.Tf("恢复.文件", z))
	fmt.Println(i18n.Tf("恢复.目标", dst, count))
}

// HasConflict 检查 zip 包解压时是否会与目标目录冲突
// 打开失败或没有冲突都返回 false
// 返回 true 时, conflicts 为冲突文件的完整路径列表
func HasConflict() (bool, []string) {
	z := ZipPath()

	r, err := zip.OpenReader(z)
	if err != nil {
		return false, nil
	}
	defer r.Close()

	dst := config.DataDir()
	conflicts := findConflicts(dst, r.File)
	if len(conflicts) == 0 {
		return false, nil
	}
	return true, conflicts
}
