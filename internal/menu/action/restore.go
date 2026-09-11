package action

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
)

// Restore 从 exe 同目录的 zip 解压回数据目录
// 若目标目录里已存在同名文件, 会先列出全部冲突, 再询问是否全部覆盖
func Restore() {
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
	if len(conflicts) > 0 {
		fmt.Println(i18n.Tf("恢复.发现冲突", len(conflicts)))
		for _, p := range conflicts {
			fmt.Println("  " + p)
		}
		fmt.Print(i18n.T("恢复.询问覆盖"))

		var answer string
		_, _ = fmt.Scanln(&answer)
		if strings.ToLower(answer) != "y" {
			fmt.Println(i18n.T("恢复.已取消"))
			return
		}
	}

	// 执行解压
	count := extractAll(dst, r.File)

	fmt.Println(i18n.T("恢复.完成"))
	fmt.Println(i18n.Tf("恢复.文件", z))
	fmt.Println(i18n.Tf("恢复.目标", dst, count))
}
