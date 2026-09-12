package action

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/pathutil"
)

// OpenGameFolder 打开游戏本体目录
// 打开前会检测路径是否存在且不为空
func OpenGameFolder() {
	gameDir := config.GameDir()
	fmt.Println(i18n.Tf("打开文件夹.游戏路径", gameDir))

	// 路径不可达, 直接返回
	if !pathutil.IsReachable(gameDir) {
		fmt.Println(i18n.Tf("打开文件夹.路径不存在", gameDir))
		return
	}

	// 空目录只警告, 依然打开
	if pathutil.IsEmptyDir(gameDir) {
		fmt.Println(i18n.Tf("打开文件夹.路径为空警告", gameDir))
	}

	// explorer 对正斜杠路径兼容性不稳定, 转成本地分隔符后再调用
	localPath := filepath.FromSlash(gameDir)
	if err := exec.Command("explorer", localPath).Start(); err != nil {
		fmt.Println(i18n.Tf("打开文件夹.失败", err))
	}
}
