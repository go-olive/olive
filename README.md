<p align="center">
  <img src="https://raw.githubusercontent.com/go-olive/brand-kit/main/banner/banner-01.png" />
</p>

[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/go-olive/olive?tab=doc)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-olive/olive/goreleaser?style=for-the-badge)](https://github.com/go-olive/olive/actions/workflows/release.yml)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-olive/olive)
[![Github All Releases](https://img.shields.io/github/downloads/go-olive/olive/total.svg?style=for-the-badge)](https://github.com/go-olive/olive/releases)

[简体中文](README_ZH.md)

**olive** is a powerful engine which monitors streamers status and automatically records when they're online. Help you catch every live stream.

If you have new features or find bugs, please go to the [issues](https://github.com/go-olive/olive/issues) section.<br/>If you get some questions or ideas, please go to the [discussions](https://github.com/go-olive/olive/discussions) section.

## Feature

- Small
- Easy-to-use
- Efficient
- Extensible
- Customizable
- Cross-platform

## Installation

- build from source

  `go install github.com/go-olive/olive/src/cmd/olive@latest`

- download from [**releases**](https://github.com/go-olive/olive/releases)

## Quickstart

Get **olive** to work simply by passing the live url.

```sh
$ olive -u https://www.huya.com/518512
```

## Usage

Start **olive** by using config file which provides you with more options.

template file to reference: [config.toml](src/tmpl/config.toml)

```sh
$ olive -f /path/to/config.toml
```

### Minimal configuration

With these minimal configs, you are already good to go.

```toml
[[Shows]]
# platform name
Platform = "bilibili"
# room id
RoomID = "21852"
# streamer name
StreamerName = "old-tomato"
```

### Custom video file name

Add config `OutTmpl`

- Date: `{{ now | date \"2006-01-02 15-04-05\"}}`

- Streame Name: `{{ .StreamerName }}`

- Stream Title: `{{ .RoomName }}`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
# The file name will be `[2022-04-24 02-02-32][old-tomato][Hi!]`
OutTmpl = "[{{ now | date \"2006-01-02 15-04-05\"}}][{{ .StreamerName }}][{{ .RoomName }}]"
```

### Custom video save location

Add config `SaveDir`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
SaveDir = "/Users/luxcgo/Videos"
```

### Custom video downloader

Add config `Parser`

```toml
[[Shows]]
Platform = "bilibili"
RoomID = "21852"
StreamerName = "old-tomato"
# Use `ffmpeg` as video downloader
Parser = "ffmpeg"
```

reference table

| Parser     | Type        | Platform                  |
| ---------- | ----------- | ------------------------- |
| streamlink | third-party | YouTube/Twitch            |
| yt-dlp     | third-party | YouTube                   |
| ffmpeg     | third-party | Other than YouTube/Twitch |
| flv        | native      | Other than YouTube/Twitch |

> You have to manually download the third-party `Parser` locally in order to use them.
>
> The deault `Parser` use `flv` which has already beed embedded into the olive , no need to download.

### Exec cmds after recording

Add config `Shows.PostCmds`, you can add a series of commands under any `[[Shows]]`.

The commands will be executed automatically when the live ends , and if any command fails to execute in the middle, it will exit early.

**olive** provides several out-of-box commands that have been implemented internally. (set config `Path` under `[[Shows.PostCmds]]`)

- `olivearchive`: Move the file to the archive folder under the current directory.
- `olivetrash`: Delete the file (unrecoverable).
- `olivebiliup`: If `UploadConfig` is configured, it will automatically upload to `bilibili` according to the configuration, if the upload fails it will execute `olivearchive`.
  - this requires to install [biliup-rs](https://github.com/ForgQi/biliup-rs) locally, and set `ExecPath` as the excutable filepath.
- `oliveshell`: split normal shell commands as an array of strings, and put them in config `Args` .
  - **olive** has embedded video file path as an environment variable which can be referenced by `$FILE_PATH` .

Config example:

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
Args = ["/bin/zsh", "-c", "echo $FILE_PATH"]
[[Shows.PostCmds]]
Path = "olivebiliup"
[[Shows.PostCmds]]
Path = "olivetrash"
```

Simulation:

1. Live ends.
2. Execute the custom command `/bin/sh -c "echo $FILE_PATH" `.
3. If the last command is executed successfully, execute the built-in command `olivebiliup `.
4. If the last command is executed successfully, execute the built-in command `olivetrash `.

### Split video files

When any of the following conditions is met, **olive** will create a new file to record.

- maximum video duration: `Duration`
  - A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
  - Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
- maximum video filesize (byte): `Filesize`

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

## Supported platforms

### native

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

**olive** relies on **[olivetv](https://github.com/go-olive/tv)** to support above sites.

If yours is not on the list above, welcome to leave your site at [discussion](https://github.com/go-olive/olive/discussions/50), or you can create a pr at **[olivetv](https://github.com/go-olive/tv)**.

### streamlink

Since streamlink is so powerful and supportive of many more sites, there's no need to reinvent the wheel. I'm using streamlink as a plugin in **olive**. As long as the site is supported in streamlink, you can use it in **olive** as well.

Examples are as follows:

```toml
[[Shows]]
Platform = "streamlink"
RoomID = "https://twitcasting.tv/c:aiueo033000"
StreamerName = "c:aiueo033000"
OutTmpl = "[{{ .StreamerName }}][{{ now | date \"2006-01-02 15-04-05\"}}]"

[[Shows]]
Platform = "streamlink"
RoomID = "https://live.nicovideo.jp/watch/lv337855989"
StreamerName = "ラムミ"
OutTmpl = "[{{ .StreamerName }}][{{ now | date \"2006-01-02 15-04-05\"}}]"

[[Shows]]
Platform = "streamlink"
RoomID = "https://play.afreecatv.com/030b1004/241795118"
StreamerName = "덩이￼"
OutTmpl = "[{{ .StreamerName }}][{{ now | date \"2006-01-02 15-04-05\"}}]"
```

## Config.toml

A config file with every feature involved.

```toml
[[Envs]]
Key = "https_proxy"
Value = "http://127.0.0.1:7890"
[[Envs]]
Key = "http_proxy"
Value = "http://127.0.0.1:7890"

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
Args = ["/bin/zsh", "-c", "echo $FILE_PATH"]
[[Shows.PostCmds]]
Path = "olivebiliup"
[[Shows.PostCmds]]
Path = "olivetrash"
```

## RoadMap

- Add docker image
- Add mock test
- Add web ui
- Add prometheus and grafana

## License

This project is under the MIT License. See the [LICENSE](https://github.com/go-olive/olive/blob/main/LICENSE) file for the full license text.
