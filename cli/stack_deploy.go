package cli

import (
	"fmt"

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
			// Load stack version
			lpath, err := utils.GetLocalPath("")
			if err != nil {
				return err
			}
			stack.Version = vcs.GetRepoVersion(lpath).String()
			fmt.Println(stack.ID, stack.Version)

			cr, err := loadCore()
			if err != nil {
				return err
			}

			st, err := cr.Stack(core.DefaultProfile())
			if err != nil {
				return err
			}

			return st.Deploy(stack)
		},
	}
}
