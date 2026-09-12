# Bodycam Save Tool

[English](README.en.md) | [中文](README.md)

A small Windows command-line tool for backing up / restoring Bodycam data, with one-click access to the game install folder (path must be configured manually).

## Features

| Option | Feature
|---|---|
| 1 | Backup Bodycam save |
| 2 | Restore Bodycam save | 
| 3 | Open game folder | 
| 4 | Settings | Language / Game install path / Data folder path / Reset all to default |

## Paths

| Name | Default path |
|---|---|
| Data folder | `%LOCALAPPDATA%\Bodycam` |
| Game install folder | `D:/SteamLibrary/steamapps/common/bodycam` |
| Backup file | `bodycam_backup.zip` next to the exe |
| Config file | `config.json` next to the exe |

## Build

Requires Go 1.27 or later.

```bash
# Development run
go mod tidy
go run .
```
## Notes
 - Close the game before backing up, otherwise the backup may fail.
 - bodycam_backup.zip is overwritten on every backup. Rename it manually if you want to keep multiple versions.
 - Do not place the exe inside the data folder or the game folder, or the backup may include itself.
