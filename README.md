# GoDown

Go 语言多线程下载器
> Goroutine Download For Golang

[![](https://img.shields.io/travis/windard/godown)](https://travis-ci.com/github/windard/godown)
[![](https://img.shields.io/tokei/lines/github/windard/godown)](https://github.com/windard/godown)
[![](https://img.shields.io/github/release/windard/godown.svg)](https://github.com/windard/godown/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/windard/godown)](https://goreportcard.com/report/github.com/windard/godown)
[![](https://img.shields.io/github/license/windard/godown)](https://github.com/windard/godown/blob/master/LICENSE)
[![](https://img.shields.io/badge/author-windard-359BE1)](https://windard.com)

```shell script
NAME:
   godown - Goroutine Download For Golang

USAGE:
   godown [global options] command [command options] argument

VERSION:
   0.2.2

COMMANDS:
   download, d  download from server
   server, s    start static server
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## download

多线程下载
> 默认 20 个并发请求，每个请求分块大小 10M

```shell script
NAME:
   godown download - download from server

USAGE:
   godown download [command options] argument

OPTIONS:
   --poolSize value, -p value   pool size for the fetch (default: 20)
   --chunkSize value, -c value  chunk size for the fetch (default: 1048576)
   --help, -h                   show help (default: false)
```

## server

静态服务器
> 默认监听 `0.0.0.0:8080` ，下载根目录为当前目录, 支持文件上传，或者离线下载

```shell script
NAME:
   godown server - start static server

USAGE:
   godown server [command options]

OPTIONS:
   --host value, -h value  server host (default: "0.0.0.0")
   --port value, -p value  server port (default: "8080")
   --root value, -r value  server root (default: ".")
   --path value            server path (default: "/")
   --list, -l              list server directory (default: false)
```
