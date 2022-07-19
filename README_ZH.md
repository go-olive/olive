<p align="center">
  <img src="https://raw.githubusercontent.com/go-olive/brand-kit/main/banner/banner-01.png" />
</p>

[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/go-olive/olive?tab=doc)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-olive/olive/goreleaser?style=for-the-badge)](https://github.com/go-olive/olive/actions/workflows/release.yml)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-olive/olive)
[![Github All Releases](https://img.shields.io/github/downloads/go-olive/olive/total.svg?style=for-the-badge)](https://github.com/go-olive/olive/releases)

**olive** 是一款强大的直播录制引擎。它会时刻监控着主播的直播状态，并在主播上线时自动开启录制。帮助您捕捉到每一场直播内容。

## 主要特性

* 小巧
* 易于使用
* 高效
* 可拓展
* 可定制
* 跨平台

## 安装部署

您可以通过以下 2 种方式来安装 **olive**：

* 源码安装

    `go install github.com/go-olive/olive/src/cmd/olive@latest`

* [二进制安装](https://github.com/go-olive/olive/releases)

## 快速开始

只需要传入直播间网址就可以让 **olive** 开始工作。

```sh
$ olive -u https://www.huya.com/518512
```

## 使用指南

通过使用配置文件启动 **olive** , 该文件为您提供了更多的选项。

模板文件参考: [config.toml](src/tmpl/config.toml)

```sh
$ olive -f /path/to/config.toml
```

### 最小配置项

只需要填写这几个配置项，就可以开始对一个直播间进行录制了。

```toml
[[Shows]]
# 平台名称
Platform = "bilibili"
# 房间号
RoomID = "21852"
# 主播名称
StreamerName = "old-tomato"
```

### 定制文件名称

新增配置项 `OutTmpl`

* 日期: `{{ now | date \"2006-01-02 15-04-05\"}}`

* 主播名称: `{{ .StreamerName }}`

* 直播标题: `{{ .RoomName }}`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
# 文件名称将会是 `[2022-04-24 02-02-32][old-tomato][Hi!]`
OutTmpl = "[{{ now | date \"2006-01-02 15-04-05\"}}][{{ .StreamerName }}][{{ .RoomName }}]"
```

### 定制文件保存位置

新增配置项 `SaveDir`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
SaveDir = "/Users/luxcgo/Videos"
```

### 定制文件下载器

新增配置项 `Parser`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
# 使用 `ffmpeg` 作为文件下载器
Parser = "ffmpeg"
```

参考表

| 下载器     | 类型   | 平台                |
| ---------- | ------ | ------------------- |
| streamlink | 第三方 | YouTube/Twitch      |
| yt-dlp     | 第三方 | YouTube             |
| ffmpeg     | 第三方 | 除了 YouTube/Twitch |
| flv        | 原生   | 除了 YouTube/Twitch |

> 如果您需要使用这些第三方下载器，请务必手动将它们下载到您的本地环境中。
>
> 原生下载器已内置到 **olive** ，您无需下载就可畅快使用。

### 录制结束后执行定制化命令

增加 `Shows.PostCmds` 配置项，在任何一个 `Shows` 的下面都可以自定义的配置多个命令。

在录制结束的时候会自动依次执行，若中途有命令执行失败，则提前退出。

**olive** 内部提供了几个已经实现好开箱即用的命令（ 在`[[Shows.PostCmds]]` 下增加 `Path` 配置项 ）

- `olivearchive`: 将文件移动到当前路径下的 `archive` 文件夹中。
- `olivetrash`: 将文件删除（不可恢复）。
- `olivebiliup`: 若有配置 `UploadConfig`，则会根据配置自动上传至哔哩哔哩，若上传失败会执行`olivearchive`。
    - 这条命令需要本地安装好 [biliup-rs](https://github.com/ForgQi/biliup-rs) ，并将 `UploadConfig` 中的 `ExecPath` 设置为可执行文件的路径。
- `oliveshell`: 将常规终端指令切分成字符串数组，并配置到 `Args` 中。
    - **olive** 内置了文件路径作为环境变量，并可以通过 `$FILE_PATH` 获取。注意环境变量只有在 shell 环境中才会正确解析，如 `/bin/sh -c "echo '$FILE_PATH'"` ，只执行 `echo '$FILE_PATH'` 则很可能获取不到。

配置文件样例

```toml
[UploadConfig]
Enable = true
ExecPath = "/xxx/biliup"
Filepath = ""

[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "test"
OutTmpl = "[test][{{ now | date \"2006-01-02 15-04-05\"}}].flv"
[[Shows.PostCmds]]
Path = "oliveshell"
Args = ["/bin/sh", "-c", "echo '$FILE_PATH'"]
[[Shows.PostCmds]]
Path = "olivebiliup"
[[Shows.PostCmds]]
Path = "olivetrash"
```

模拟执行配置文件

1. 录制结束
2. 执行自定义命令`/bin/sh -c "echo '$FILE_PATH'"`
3. 若上条命令执行成功，执行内置命令`olivebiliup`
4. 若上条命令执行成功，执行内置命令`olivetrash`

### 切分文件

当以下任意一个条件满足是时， **olive** 会新创建一个文件用于录制。

* 最大视频时长: `Duration`
    * 一个时间段字符串是一个序列，每个片段包含可选的正负号、十进制数、可选的小数部分和单位后缀，如"300ms"、"-1.5h"、"2h45m"。
    * 合法的单位有"ns"、"us" /"µs"、"ms"、"s"、"m"、"h"。
* 最大视频大小(字节): `Filesize`

```toml
[[Shows]]
Platform = "huya"
RoomID = "518512"
StreamerName = "250"
[Shows.SplitRule]
# 10 MB
FileSize = 10000000
# 1 minute
Duration = "1m"
[[Shows.PostCmds]]
Path = "oliveshell"
Args = ["/bin/sh", "-c", "echo $FILE_PATH"]
```

## 直播平台

| Platform |
| -------- |
| bilibili |
| douyin   |
| huya     |
| kuaishou |
| lang     |
| tiktok   |
| twitch   |
| youtube  |

**olive** 依赖 **[olivetv](https://github.com/go-olive/tv)** 来支持上述网站的直播录制。如果您的不在上述列表中，欢迎在 **[olivetv](https://github.com/go-olive/tv)** 中提交 issue 或 pr 。


## Config.toml

一个包含所有特性的配置文件。

```toml
LogLevel = 5
SnapRestSeconds = 15
SplitRestSeconds = 60
CommanderPoolSize = 1

[UploadConfig]
Enable = true
ExecPath = "biliup"
Filepath = ""

[PlatformConfig]
DouyinCookie = "__ac_nonce=06245c89100e7ab2dd536; __ac_signature=_02B4Z6wo00f01LjBMSAAAIDBwA.aJ.c4z1C44TWAAEx696;"
KuaishouCookie = "did=web_d86297aa2f579589b8abc2594b0ea985"

[[Shows]]
Platform = "huya"
RoomID = "518512"
StreamerName = "250"

[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
SaveDir = "/Users/luxcgo/Videos"
Parser = "flv"
OutTmpl = "[{{ now | date \"2006-01-02 15-04-05\"}}][{{ .StreamerName }}]"
[Shows.SplitRule]
# 1 GB
FileSize = 1024000000
# 1 hour
Duration = "1h"
[[Shows.PostCmds]]
Path = "oliveshell"
Args = ["/bin/sh", "-c", "echo '$FILE_PATH'"]
[[Shows.PostCmds]]
Path = "olivebiliup"
[[Shows.PostCmds]]
Path = "olivetrash"
```

## 授权许可

本项目采用 MIT 开源授权许可证，完整的授权说明已放置在 [LICENSE](https://github.com/go-olive/olive/blob/main/LICENSE) 文件中。