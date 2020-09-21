package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/windard/godown/fetch"
	"os"
)

var version = "0.0.1"


func main() {
	var poolSize int64
	var chunkSize int64
	app := &cli.App{
		Name: "godown",
		Usage: "Goroutine Download For Golang",
		UsageText: "godown [global options] argument",
		Version: version,
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
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
