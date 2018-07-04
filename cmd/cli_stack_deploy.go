package main

import (
	"fmt"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	"gopkg.in/urfave/cli.v2"
)

func commandStackDeploy() *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: "Deploy stack",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"dry"},
				Usage:   "perform a dry run",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			lpath, _ := utils.GetLocalPath("")
			stack.Version = vcs.GetRepoVersion(lpath).String()
			fmt.Println(stack.ID, stack.Version)

			pconf, err := config.ReadProjectConfig(lpath)
			if err != nil {
				return err
			}

			conf := &core.Config{
				PacksDir:    consts.DefaultPacksDir,
				ThrapConfig: pconf,
			}
			cr, err := core.NewCore(conf)
			if err != nil {
				return err
			}

			st := cr.Stack()
			err = st.Deploy(stack)
			// opt := orchestrator.RequestOptions{
			// 	Dryrun: ctx.Bool("dryrun"),
			// 	Output: os.Stdout,
			// }

			//err = core.DeployStack(st, opt)
			return err

			// conf := &orchestrator.Config{
			// 	Provider: "nomad",
			// }
			// orch, err := orchestrator.New(conf)
			// if err != nil {
			// 	return err
			// }

			// opt := orchestrator.RequestOptions{Dryrun: ctx.Bool("dryrun")}
			// _, job, err := orch.Deploy(st, opt)
			// if err == nil {
			// 	b, _ := json.MarshalIndent(job, "", "  ")
			// 	os.Stdout.Write(b)
			// 	os.Stdout.Write([]byte("\n"))
			// }
			// return err
		},
	}
}
