package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/models"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/trigger"
	"github.com/urfave/cli"
)

const (
	waitTimeout      = 10 * time.Minute
	thresholdTimeout = 2 * time.Second
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
				Value: trigger.DefaultHost,
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
				Name:  "privateToken",
				Usage: "your personal private token",
				Value: "",
			},
			cli.StringFlag{
				Name:  "token",
				Usage: "project token",
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
					schema, err := schemaConverter(ctx.GlobalString("schema"))
					if err != nil {
						log.Fatal(err)
						return
					}
					trigger := trigger.New(
						ctx.GlobalBool("skipVerifyTLS"),
						schema,
						ctx.GlobalString("host"),
						ctx.GlobalString("privateToken"),
						ctx.GlobalString("token"),
						ctx.GlobalString("ref"),
						ctx.GlobalInt("projectID"),
						ctx.GlobalStringSlice("variables"),
					)

					triggerRunResp, err := trigger.RunPipeline()
					if err != nil {
						log.Fatal(err)
					}

					go func() {
						shutdown()
						time.Sleep(waitTimeout)
						os.Exit(1)
					}()

					var previousPipelineStatus models.PipelineStatus
					for {
						triggerCompletionResp, err := trigger.PollForCompletion(triggerRunResp.ID)
						if err != nil {
							log.Fatal(err)
						}

						switch triggerCompletionResp.Status {
						case models.Pending, models.Running:
							if previousPipelineStatus != triggerCompletionResp.Status {
								log.Printf("pipeline status: %s\n", triggerCompletionResp.Status)
								previousPipelineStatus = triggerCompletionResp.Status
							}
						case models.Failed:
							log.Fatalln("pipeline failed!")
						case models.Success:
							log.Println("pipeline success!")
							os.Exit(0)
						default:
							log.Printf("unknown status: %s", triggerCompletionResp.Status)
							os.Exit(1)
						}

						time.Sleep(thresholdTimeout)
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("start application failed: %s", err)
	}
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutting down...")
	os.Exit(0)
}
