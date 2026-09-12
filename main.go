package main

import (
	"fmt"
	"os"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/menu"
)

// main 主函数
func main() {
	// 加载配置
	config.Load()
	// 配置里的语言可能不受支持
	if !i18n.IsSupported(config.Get().Language) {
		// 校正一下
		config.SetLanguage(config.DefaultLanguage)
	}
	// 加载语言包
	if err := i18n.Load(config.Get().Language); err != nil {
		// 语言文件损坏, 无法继续
		fmt.Fprintln(os.Stderr, "load language failed:", err)
		os.Exit(1)
	}
	// 运行菜单
	menu.Run()
}
