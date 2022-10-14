package main

import (
	// "fmt"
	"os"

	ge "github.com/jielyu/graph_engine_go"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	app := &cli.App{
		Name:  "ge_tool",
		Usage: "用于图引擎辅助工作",
		Commands: []*cli.Command{
			{
				Name: "check",
				Action: func(ctx *cli.Context) error {
					jsonFile := ctx.Args().Get(0)
					if len(jsonFile) == 0 {
						log.Fatalf("please specific config file path")
						return nil
					}
					config, err := ge.LoadGraphConfig(jsonFile)
					if err != nil {
						log.Fatalf("failed to load config from %s\r\n", &jsonFile)
						return nil
					}
					log.Printf("load config successfully\r\n, config=%v\r\n", config)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
