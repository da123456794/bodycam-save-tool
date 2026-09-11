package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	currentLang = "zh"
	messages    = map[string]string{}
)

func Load(lang string) error {
	if lang != "zh" && lang != "en" {
		lang = "zh"
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

func Current() string { return currentLang }

func T(key string) string {
	if s, ok := messages[key]; ok {
		return s
	}
	return key
}

func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
