package main

import (
	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/menu"
)

// main 主函数
func main() {
	config.Load()
	_ = i18n.Load(config.Get().Language)
	menu.Run()
}
