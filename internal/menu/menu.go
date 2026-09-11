package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/menu/action"
)

// reader 从标准输入读取输入
var reader = bufio.NewReader(os.Stdin)

// readLine 从标准输入读取一行输入
// Returns:
//
//	line: 输入的行，已去空格
//	ok: 是否成功读取输入
func readLine() (string, bool) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(line), true
}

// pause 等待按下回车
func pause() {
	fmt.Print(i18n.T("通用.按Enter继续"))
	_, _ = reader.ReadString('\n')
}

// Run 主菜单
func Run() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("主菜单.标题"))
		fmt.Println(i18n.T("主菜单.备份"))
		fmt.Println(i18n.T("主菜单.加载"))
		fmt.Println(i18n.T("主菜单.打开游戏文件夹"))
		fmt.Println(i18n.T("主菜单.设置"))
		fmt.Println(i18n.T("通用.退出"))
		fmt.Print(i18n.T("通用.提示"))

		line, ok := readLine()
		if !ok {
			return
		}
		switch line {
		case "1":
			action.Backup()
			pause()
		case "2":
			action.Restore()
			pause()
		case "3":
			action.OpenGameFolder()
			pause()
		case "4":
			settingsMenu()
		case "q", "Q":
			return
		default:
			fmt.Println(i18n.T("通用.无效选项"))
		}
	}
}

// settingsMenu 设置菜单
func settingsMenu() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("设置菜单.标题"))
		fmt.Println(i18n.T("设置菜单.语言"))
		fmt.Println(i18n.T("设置菜单.游戏路径"))
		fmt.Println(i18n.T("设置菜单.数据路径"))
		fmt.Println(i18n.T("设置菜单.恢复默认"))
		fmt.Println(i18n.T("通用.返回"))
		fmt.Println(i18n.T("通用.退出"))
		fmt.Print(i18n.T("通用.提示"))

		line, ok := readLine()
		if !ok {
			return
		}
		switch line {
		case "1":
			languageMenu()
		case "2":
			gamePathMenu()
			pause()
		case "3":
			dataPathMenu()
			pause()
		case "4":
			resetMenu()
			pause()
		case "b", "B":
			return
		case "q", "Q":
			os.Exit(0)
		default:
			fmt.Println(i18n.T("通用.无效选项"))
		}
	}
}

// languageMenu 语言设置菜单
func languageMenu() {
	currentName := i18n.T("语言菜单.中文名")
	if config.Get().Language == "en" {
		currentName = i18n.T("语言菜单.英文名")
	}

	fmt.Println()
	fmt.Println(i18n.T("语言菜单.标题"))
	fmt.Println(i18n.Tf("语言菜单.当前语言", currentName))
	fmt.Println(i18n.T("语言菜单.中文"))
	fmt.Println(i18n.T("语言菜单.英文"))
	fmt.Println(i18n.T("通用.返回"))
	fmt.Print(i18n.T("通用.提示"))

	line, ok := readLine()
	if !ok {
		return
	}
	switch line {
	case "1":
		config.SetLanguage("zh")
		_ = i18n.Load("zh")
		fmt.Println(i18n.Tf("语言菜单.已切换", i18n.T("语言菜单.中文名")))
	case "2":
		config.SetLanguage("en")
		_ = i18n.Load("en")
		fmt.Println(i18n.Tf("语言菜单.已切换", i18n.T("语言菜单.英文名")))
	case "b", "B":
		return
	default:
		fmt.Println(i18n.T("通用.无效选项"))
	}
}

// promptPath 通用路径输入, 返回用户输入和是否继续
// Args:
//
//	titleKey: 标题键
//	currentKey: 当前路径键
//	promptKey: 输入示键
//	current: 当前路径
//
// Returns:
//
//	line: 用户输入的路径
//	ok: 是否成功读取输入
func promptPath(titleKey, currentKey, promptKey, current string) (string, bool) {
	fmt.Println()
	fmt.Println(i18n.T(titleKey))
	fmt.Println(i18n.Tf(currentKey, current))
	fmt.Print(i18n.T(promptKey))

	line, ok := readLine()
	if !ok {
		return "", false
	}
	return line, true
}

// gamePathMenu 游戏本体路径设置菜单
func gamePathMenu() {
	line, ok := promptPath(
		"游戏路径菜单.标题",
		"游戏路径菜单.当前路径",
		"游戏路径菜单.输入提示",
		config.GameDir(),
	)
	if !ok || line == "" {
		return
	}
	config.SetGamePath(line)
	fmt.Println(i18n.Tf("游戏路径菜单.已更新", line))
}

// dataPathMenu 游戏数据文件夹路径设置菜单
// Args:
//
//	titleKey: 标题键
//	currentKey: 当前路径键
//	promptKey: 输入示键
//	current: 当前路径
//
// Returns:
//
//	line: 用户输入的路径
//	ok: 是否成功读取输入
func dataPathMenu() {
	line, ok := promptPath(
		"数据路径菜单.标题",
		"数据路径菜单.当前路径",
		"数据路径菜单.输入提示",
		config.DataDir(),
	)
	if !ok || line == "" {
		return
	}
	// 输入 d 或 default 恢复默认
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetDataPath("")
		fmt.Println(i18n.Tf("数据路径菜单.已恢复默认", config.DataDir()))
		return
	}
	config.SetDataPath(line)
	fmt.Println(i18n.Tf("数据路径菜单.已更新", line))
}

// resetMenu 恢复全部默认设置
func resetMenu() {
	fmt.Println()
	fmt.Println(i18n.T("恢复默认菜单.标题"))
	fmt.Print(i18n.T("恢复默认菜单.确认"))

	line, ok := readLine()
	if !ok {
		return
	}
	if !strings.EqualFold(line, "y") {
		fmt.Println(i18n.T("恢复默认菜单.已取消"))
		return
	}

	// 恢复配置
	config.Reset()

	// 语言可能变回 zh, 重新加载
	_ = i18n.Load(config.Get().Language)

	// 用新语言打印结果
	fmt.Println(i18n.T("恢复默认菜单.已完成"))
	fmt.Println(i18n.Tf("恢复默认菜单.当前语言", i18n.T("语言菜单.中文名")))
	fmt.Println(i18n.Tf("恢复默认菜单.当前游戏路径", config.GameDir()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前数据路径", config.DataDir()))
}
