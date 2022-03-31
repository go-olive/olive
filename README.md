# lifesaver


[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/luxcgo/lifesaver?tab=doc)
[![GoReport](https://goreportcard.com/badge/github.com/luxcgo/lifesaver?style=for-the-badge)](https://goreportcard.com/report/github.com/luxcgo/lifesaver)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-ini/ini)

## Save Lives!!

Lives are delicate and fleeting creatures, waiting to be captured by us. ❤

> 全自动录播、投稿工具
>
> 支持抖音直播、虎牙直播、B站直播、油管直播
>
> 支持B站投稿


## Usage

```sh
// install first
go install github.com/luxcgo/lifesaver@latest

// then run
lifesaver -c /path/to/config.toml
```

## Config.toml

template file to reference [config.toml](tmpl/config.toml)

```toml
[UploadConfig]
// 是否上传到 bilibili
Enable = false
// biliup-rs 执行路径
ExecPath = "biliup"
// biliup-rs 配置文件路径，为空的话走默认配置
Filepath = ""

[PlatformConfig]
// 若有录制抖音直播，可在无痕模式非登录状态下找下面的 cookie 填入即可
DouyinCookie = "__ac_nonce=06245c89100e7ab2dd536; __ac_signature=_02B4Z6wo00f01LjBMSAAAIDBwA.aJ.c4z1C44TWAAEx696;"

[[Shows]]
// 平台名，目前支持：
// "bilibili"
// "douyin"
// "huya"
// "youtube"
Platform = "bilibili"
// 房间号，支持字符串类型的房间号
RoomID = "21852"
// 主播名称
StreamerName = "老番茄"
```

## RoodMap

* 支持 go 原生对视频流的抓取，去除 ffmpeg 和 streamlink 的依赖
* 支持 go 原生对 bilibili 的投稿，去除 biliup-rs 的依赖

* 增加对更多直播平台的支持
* 增加对程序运行状况的监控
* 增加网页端

## Credits

* [bililive-go](https://github.com/hr3lxphr6j/bililive-go)
* [biliup-rs](https://github.com/ForgQi/biliup-rs)
* [ffmpeg](https://ffmpeg.org/)
* [streamlink](https://streamlink.github.io/)
