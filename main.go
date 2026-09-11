package main

import (
	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/menu"
)

// main 主函数
func main() {
	// 加载配置
	cfg := config.Load()
	// 加载语言
	err := i18n.Load(cfg.Language)
	if err != nil {
		// 加载失败，使用默认语言
		cfg.Language = "zh"
		_ = i18n.Load("zh")
	}
	// 运行菜单
	menu.Run()
}
