package main

import (
	"log"
	"os"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/pipeline"
	"github.com/urfave/cli"
)

func main() {
	app := cli.App{
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "skipVerifyTLS",
			},
			cli.StringFlag{
				Name:  "schema",
				Usage: "http schema",
				Value: "https",
			},
			cli.StringFlag{
				Name:  "host",
				Value: pipeline.DefaultHost,
			},
			cli.StringFlag{
				Name:  "ref",
				Usage: "for change the branch",
				Value: "master",
			},
			cli.StringFlag{
				Name:  "urlPrefix",
				Value: "/",
			},
			cli.IntFlag{
				Name:  "projectID",
				Usage: "project ID",
				Value: 0,
			},
			cli.StringFlag{
				Name:  "token",
				Value: "",
			},
			cli.StringSliceFlag{
				Name:  "variables",
				Usage: "variable1:value,variable2:value",
				Value: nil,
			},
		},
		Commands: []cli.Command{
			{
				Name: "run",
				Action: func(ctx *cli.Context) {
					pipeln := pipeline.New(ctx.GlobalBool("skipVerifyTLS"))
					schema, err := schemaConverter(ctx.GlobalString("schema"))
					if err != nil {
						log.Fatal(err)
						return
					}

					pipelineResp, err := pipeln.Run(
						schema,
						ctx.GlobalString("host"),
						ctx.GlobalString("token"),
						ctx.GlobalString("ref"),
						ctx.GlobalInt("projectID"),
						ctx.GlobalStringSlice("variables"),
					)
					if err != nil {
						log.Fatal(err)
						return
					}

					log.Printf("RESP: %+v", pipelineResp)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Panicf("start application failed: %s", err)
	}
}
