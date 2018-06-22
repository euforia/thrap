package main

import (
	"encoding/json"
	"os"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackDeploy() *cli.Command {
	return &cli.Command{
		Name: "deploy",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"dry"},
				Usage:   "perform a dry run",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			st, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			if errs := st.Validate(); errs != nil {
				return utils.FlattenErrors(errs)
			}
			conf := &orchestrator.Config{
				Provider: "nomad",
			}
			orch, err := orchestrator.New(conf)
			if err != nil {
				return err
			}

			opt := orchestrator.DeployOptions{Dryrun: ctx.Bool("dryrun")}
			_, job, err := orch.Deploy(st, opt)
			if err == nil {
				b, _ := json.MarshalIndent(job, "", "  ")
				os.Stdout.Write(b)
				os.Stdout.Write([]byte("\n"))
			}
			return err
		},
	}
}
