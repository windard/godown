package main

import (
	"github.com/urfave/cli/v2"
	"github.com/windard/godown/fetch"
	"github.com/windard/godown/server"
	"log"
	"os"
)

var version = "0.1.0"


func main() {
	var poolSize int64
	var chunkSize int64

	var host, port string
	var path, root string
	var listDirectory bool

	app := &cli.App{
		Name: "godown",
		Usage: "Goroutine Download For Golang",
		UsageText: "godown [global options] command [command options] argument",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "download from server",
				UsageText: "godown download [command options] argument",
				Flags: []cli.Flag {
					&cli.Int64Flag{
						Aliases: 	[]string{"p"},
						Name:        "poolSize",
						Value:       20,
						Usage:       "pool size for the fetch",
						Destination: &poolSize,
					},
					&cli.Int64Flag{
						Aliases: 	[]string{"c"},
						Name:        "chunkSize",
						Value:       1024 * 1024,
						Usage:       "chunk size for the fetch",
						Destination: &chunkSize,
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						_ = cli.ShowAppHelp(c)
						return nil
					}

					requestUrl := c.Args().Get(0)
					fetch.GoroutineDownload(requestUrl, poolSize, chunkSize)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start static server",
				UsageText: "godown server [command options]",
				HideHelp: true,
				Flags: []cli.Flag {
					&cli.StringFlag{
						Aliases: 	[]string{"h"},
						Name:        "host",
						Value:       "0.0.0.0",
						Usage:       "server host",
						Destination: &host,
					},
					&cli.StringFlag{
						Aliases: 	[]string{"p"},
						Name:        "port",
						Value:       "8080",
						Usage:       "server port",
						Destination: &port,
					},
					&cli.StringFlag{
						Aliases: 	[]string{"r"},
						Name:        "root",
						Value:       ".",
						Usage:       "server root",
						Destination: &root,
					},
					&cli.StringFlag{
						Name:        "path",
						Value:       "/",
						Usage:       "server path",
						Destination: &path,
					},
					&cli.BoolFlag{
						Aliases: 	[]string{"l"},
						Name:        "list",
						Value:       false,
						Usage:       "list server directory",
						Destination: &listDirectory,
					},
				},
				Action: func(c *cli.Context) error {
					server.StaticServerFileSystem(host, port, path, root, listDirectory)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
