package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed locales/*.json
var localeFS embed.FS

// LangInfo 一种语言的数据结构体
type LangInfo struct {
	Code   string // 语言代码, 也是 locales/<语言代码>.json 的文件名
	Native string // 该语言的本地名, 不随界面语言变化
}

// supportedLangs 支持的语言列表
// 想加新语言: 在这里加一行, 并放一个 locales/<语言代码>.json
var supportedLangs = []LangInfo{
	{Code: "zh", Native: "中文"},
	{Code: "en", Native: "English"},
}

const defaultLang = "zh"

var (
	currentLang = defaultLang
	messages    = map[string]string{}
)

// Available 返回所有支持的语言
// Returns:
//
//	[]LangInfo: 所有支持的语言列表
func Available() []LangInfo {
	return supportedLangs
}

// IsSupported 判断语言代码是否受支持
func IsSupported(code string) bool {
	for _, l := range supportedLangs {
		if l.Code == code {
			return true
		}
	}
	return false
}

// DisplayName 返回 "语言代码(本地名)" 形式的显示名
// 例如 "zh(中文)"
// 未支持的语言直接返回代码本身
// Args:
//
//	code: 语言代码
//
// Returns:
//
//	显示名
func DisplayName(code string) string {
	for _, l := range supportedLangs {
		if l.Code == code {
			return l.Code + "(" + l.Native + ")"
		}
	}
	return code
}

// Load 加载语言
// 传入不支持的语言会回退到默认语言
// Args:
//
//	lang: 要加载的语言代码
//
// Returns:
//
//	error: 如果加载失败, 则返回错误
//	如果成功, 则返回 nil
func Load(lang string) error {
	if !IsSupported(lang) {
		lang = defaultLang
	}
	data, err := localeFS.ReadFile("locales/" + lang + ".json")
	if err != nil {
		return err
	}
	m := map[string]string{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	messages = m
	currentLang = lang
	return nil
}

// Current 返回当前语言代码
func Current() string { return currentLang }

// T 返回翻译后的字符串
// Args:
//
//	key: 翻译键
//
// Returns:
//
//	翻译后的字符串
func T(key string) string {
	if s, ok := messages[key]; ok {
		return s
	}
	return key
}

// Tf 返回格式后的翻译后的字符串
// Args:
//
//	key: 翻译键
//	args: 格式化参数
//
// Returns:
//
//	格式后的翻译后的字符串
func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
