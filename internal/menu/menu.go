package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"bodycam-save-tool/internal/config"
	"bodycam-save-tool/internal/i18n"
	"bodycam-save-tool/internal/menu/action"
	"bodycam-save-tool/internal/version"
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
		fmt.Println(i18n.Tf("通用.版本", version.Number))
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
			backupMenu()
		case "2":
			restoreMenu()
		case "3":
			openFolderMenu()
		case "4":
			settingsMenu()
		case "q", "Q":
			return
		default:
			fmt.Println(i18n.T("通用.无效选项"))
		}
	}
}

// ========== 备份菜单 ==========

// backupMenu 备份菜单
// 除了三个备份动作, 还承担备份文件名设置
func backupMenu() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("备份菜单.标题"))
		fmt.Println(i18n.T("备份菜单.数据"))
		fmt.Println(i18n.T("备份菜单.存档"))
		fmt.Println(i18n.T("备份菜单.配置"))
		fmt.Println(i18n.T("备份菜单.设置数据名"))
		fmt.Println(i18n.T("备份菜单.设置存档名"))
		fmt.Println(i18n.T("备份菜单.设置配置名"))
		fmt.Println(i18n.T("通用.返回"))
		fmt.Println(i18n.T("通用.退出"))
		fmt.Print(i18n.T("通用.提示"))

		line, ok := readLine()
		if !ok {
			return
		}
		switch line {
		case "1":
			action.BackupData()
			pause()
		case "2":
			action.BackupSaveGames()
			pause()
		case "3":
			action.BackupWindowsConfig()
			pause()
		case "4":
			dataBackupNameMenu()
			pause()
		case "5":
			saveGamesBackupNameMenu()
			pause()
		case "6":
			windowsConfigBackupNameMenu()
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

// ========== 加载菜单 ==========

// restoreMenu 加载菜单
// 除了三个加载动作, 还承担三个目标路径的设置
func restoreMenu() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("加载菜单.标题"))
		fmt.Println(i18n.T("加载菜单.数据"))
		fmt.Println(i18n.T("加载菜单.存档"))
		fmt.Println(i18n.T("加载菜单.配置"))
		fmt.Println(i18n.T("加载菜单.设置数据路径"))
		fmt.Println(i18n.T("加载菜单.设置存档路径"))
		fmt.Println(i18n.T("加载菜单.设置配置路径"))
		fmt.Println(i18n.T("通用.返回"))
		fmt.Println(i18n.T("通用.退出"))
		fmt.Print(i18n.T("通用.提示"))

		line, ok := readLine()
		if !ok {
			return
		}
		switch line {
		case "1":
			runRestore(action.DataBackupPath(), config.DataDir())
			pause()
		case "2":
			runRestore(action.SaveGamesBackupPath(), config.SaveGamesDir())
			pause()
		case "3":
			runRestore(action.WindowsConfigBackupPath(), config.WindowsConfigDir())
			pause()
		case "4":
			dataPathMenu()
			pause()
		case "5":
			saveGamesPathMenu()
			pause()
		case "6":
			windowsConfigPathMenu()
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

// runRestore 处理恢复流程: 检查冲突 -> 询问 -> 执行
func runRestore(zipPath, dst string) {
	hasConflict, conflicts, err := action.HasConflict(zipPath, dst)
	if err != nil {
		fmt.Println(i18n.Tf("恢复.打不开压缩包", zipPath))
		return
	}
	if !hasConflict {
		action.Restore(zipPath, dst, false)
		return
	}

	fmt.Println(i18n.Tf("恢复.发现冲突", len(conflicts)))
	for _, p := range conflicts {
		fmt.Println("  " + p)
	}
	fmt.Print(i18n.T("恢复.询问覆盖"))

	line, ok := readLine()
	if !ok {
		return
	}
	if !strings.EqualFold(line, "y") {
		fmt.Println(i18n.T("恢复.已取消"))
		return
	}

	action.Restore(zipPath, dst, true)
}

// ========== 打开游戏文件夹菜单 ==========

// openFolderMenu 打开游戏文件夹菜单
// 承载游戏本体路径的设置
func openFolderMenu() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("打开文件夹菜单.标题"))
		fmt.Println(i18n.T("打开文件夹菜单.打开"))
		fmt.Println(i18n.T("打开文件夹菜单.设置路径"))
		fmt.Println(i18n.T("通用.返回"))
		fmt.Println(i18n.T("通用.退出"))
		fmt.Print(i18n.T("通用.提示"))

		line, ok := readLine()
		if !ok {
			return
		}
		switch line {
		case "1":
			action.OpenGameFolder()
			pause()
		case "2":
			gamePathMenu()
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

// ========== 设置菜单 ==========

// settingsMenu 设置菜单
// 只保留语言与恢复默认, 路径和备份文件名都挪到对应菜单
func settingsMenu() {
	for {
		fmt.Println()
		fmt.Println(i18n.T("设置菜单.标题"))
		fmt.Println(i18n.T("设置菜单.语言"))
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
	langs := i18n.Available()

	fmt.Println()
	fmt.Println(i18n.T("语言菜单.标题"))
	fmt.Println(i18n.Tf("语言菜单.当前语言", i18n.DisplayName(config.Get().Language)))
	for i, l := range langs {
		fmt.Printf("%d. %s\n", i+1, i18n.DisplayName(l.Code))
	}
	fmt.Println(i18n.T("通用.返回"))
	fmt.Print(i18n.T("通用.提示"))

	line, ok := readLine()
	if !ok {
		return
	}
	if strings.EqualFold(line, "b") {
		return
	}

	// 数字对应语言列表下标
	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(langs) {
		fmt.Println(i18n.T("通用.无效选项"))
		return
	}

	// -1 是为了匹配语言列表的索引
	code := langs[idx-1].Code
	// 设置语言
	config.SetLanguage(code)
	_ = i18n.Load(code)
	fmt.Println(i18n.Tf("语言菜单.已切换", i18n.DisplayName(code)))
}

// ========== 通用输入 ==========

// prompt 通用输入
// Args:
//
//	titleKey: 标题键
//	currentKey: 当前路径键
//	promptKey: 输入提示键
//	current: 当前路径
//
// Returns:
//
//	line: 用户输入的路径
//	ok: 是否成功读取输入
func prompt(titleKey, currentKey, promptKey, current string) (string, bool) {
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

// ========== 路径设置(由加载菜单调用) ==========

// gamePathMenu 游戏本体路径设置菜单
func gamePathMenu() {
	line, ok := prompt(
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
func dataPathMenu() {
	line, ok := prompt(
		"数据路径菜单.标题",
		"数据路径菜单.当前路径",
		"数据路径菜单.输入提示",
		config.DataDir(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetDataPath("")
		fmt.Println(i18n.Tf("数据路径菜单.已恢复默认", config.DataDir()))
		return
	}
	config.SetDataPath(line)
	fmt.Println(i18n.Tf("数据路径菜单.已更新", line))
}

// saveGamesPathMenu 游戏存档文件夹路径设置菜单
func saveGamesPathMenu() {
	line, ok := prompt(
		"存档路径菜单.标题",
		"存档路径菜单.当前路径",
		"存档路径菜单.输入提示",
		config.SaveGamesDir(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetSaveGamesPath("")
		fmt.Println(i18n.Tf("存档路径菜单.已恢复默认", config.SaveGamesDir()))
		return
	}
	config.SetSaveGamesPath(line)
	fmt.Println(i18n.Tf("存档路径菜单.已更新", line))
}

// windowsConfigPathMenu 游戏配置文件夹路径设置菜单
func windowsConfigPathMenu() {
	line, ok := prompt(
		"配置路径菜单.标题",
		"配置路径菜单.当前路径",
		"配置路径菜单.输入提示",
		config.WindowsConfigDir(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetWindowsConfigPath("")
		fmt.Println(i18n.Tf("配置路径菜单.已恢复默认", config.WindowsConfigDir()))
		return
	}
	config.SetWindowsConfigPath(line)
	fmt.Println(i18n.Tf("配置路径菜单.已更新", line))
}

// ========== 备份文件名设置(由备份菜单调用) ==========

// dataBackupNameMenu 数据备份文件名设置菜单
func dataBackupNameMenu() {
	line, ok := prompt(
		"数据备份名菜单.标题",
		"数据备份名菜单.当前",
		"数据备份名菜单.输入提示",
		config.DataBackupName(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetDataBackupName("")
		fmt.Println(i18n.Tf("数据备份名菜单.已恢复默认", config.DataBackupName()))
		return
	}
	config.SetDataBackupName(line)
	fmt.Println(i18n.Tf("数据备份名菜单.已更新", config.DataBackupName()))
}

// saveGamesBackupNameMenu 存档备份文件名设置菜单
func saveGamesBackupNameMenu() {
	line, ok := prompt(
		"存档备份名菜单.标题",
		"存档备份名菜单.当前",
		"存档备份名菜单.输入提示",
		config.SaveGamesBackupName(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetSaveGamesBackupName("")
		fmt.Println(i18n.Tf("存档备份名菜单.已恢复默认", config.SaveGamesBackupName()))
		return
	}
	config.SetSaveGamesBackupName(line)
	fmt.Println(i18n.Tf("存档备份名菜单.已更新", config.SaveGamesBackupName()))
}

// windowsConfigBackupNameMenu 配置备份文件名设置菜单
func windowsConfigBackupNameMenu() {
	line, ok := prompt(
		"配置备份名菜单.标题",
		"配置备份名菜单.当前",
		"配置备份名菜单.输入提示",
		config.WindowsConfigBackupName(),
	)
	if !ok || line == "" {
		return
	}
	if strings.EqualFold(line, "d") || strings.EqualFold(line, "default") {
		config.SetWindowsConfigBackupName("")
		fmt.Println(i18n.Tf("配置备份名菜单.已恢复默认", config.WindowsConfigBackupName()))
		return
	}
	config.SetWindowsConfigBackupName(line)
	fmt.Println(i18n.Tf("配置备份名菜单.已更新", config.WindowsConfigBackupName()))
}

// ========== 恢复默认 ==========

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

	// 语言可能变回默认, 重新加载
	_ = i18n.Load(config.Get().Language)

	// 用新语言打印结果
	fmt.Println(i18n.T("恢复默认菜单.已完成"))
	fmt.Println(i18n.Tf("恢复默认菜单.当前语言", i18n.DisplayName(config.Get().Language)))
	fmt.Println(i18n.Tf("恢复默认菜单.当前游戏路径", config.GameDir()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前数据路径", config.DataDir()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前存档路径", config.SaveGamesDir()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前配置路径", config.WindowsConfigDir()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前数据备份名", config.DataBackupName()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前存档备份名", config.SaveGamesBackupName()))
	fmt.Println(i18n.Tf("恢复默认菜单.当前配置备份名", config.WindowsConfigBackupName()))
}
