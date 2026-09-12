# Bodycam Save Tool / Bodycam 存档工具

[中文](#中文) | [English](README.en.md)

---

## 中文

一个 Windows 下的命令行小工具，用来备份 / 恢复 Bodycam 的数据，并支持一键打开游戏本体目录（需要自己配置）。

### 功能

| 选项 | 功能 |
|---|---|
| 1 | 备份 Bodycam 存档 |
| 2 | 加载 Bodycam 存档 |
| 3 | 一键打开游戏文件夹 |
| 4 | 设置 | 语言 / 游戏本体路径 / 数据文件夹路径 / 恢复全部默认 |

### 目录与路径

| 名称 | 默认路径 |
|---|---|
| 数据目录 | `%LOCALAPPDATA%\Bodycam` |
| 游戏本体目录 | `D:/SteamLibrary/steamapps/common/bodycam` |
| 备份文件 | exe 同目录下的 `bodycam_backup.zip` |
| 配置文件 | exe 同目录下的 `config.json` |

### 构建

要求 Go 1.27 或更高版本。

```bash
go mod tidy

go run .

go build -o BodycamSaveTool.exe .
```

## 注意
- 请确保在备份前关闭游戏, 否则可能会导致备份失败
- 对于bodycam_backup.zip是覆盖模式
- 不要将exe文件放在data文件夹或游戏文件夹内,否则备份可能会包含自身