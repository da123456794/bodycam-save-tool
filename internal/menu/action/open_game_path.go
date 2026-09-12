package action

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
)

// OpenGameFolder 打开游戏本体目录
func OpenGameFolder() {
	gameDir := config.GameDir()
	fmt.Println(i18n.Tf("打开文件夹.游戏路径", gameDir))

	// explorer 对正斜杠路径兼容性不稳定, 转成本地分隔符后再调用
	localPath := filepath.FromSlash(gameDir)
	if err := exec.Command("explorer", localPath).Start(); err != nil {
		fmt.Println(i18n.Tf("打开文件夹.失败", err))
	}
}
