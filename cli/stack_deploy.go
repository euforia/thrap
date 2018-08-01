package cli

import (
	"fmt"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/store"
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
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "deployment `profile`",
				Value:   "default",
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

			// Load profiles
			profs, err := store.LoadHCLFileProfileStorage(lpath)
			if err != nil {
				return err
			}

			// Load request profile
			profName := ctx.String("profile")
			prof := profs.Get(profName)
			if prof == nil {
				return fmt.Errorf("profile not found: %s", profName)
			}

			stack.Version = vcs.GetRepoVersion(lpath).String()
			fmt.Println(stack.ID, stack.Version)

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			st, err := cr.Stack(prof)
			if err != nil {
				return err
			}

			return st.Deploy(stack)
		},
	}
}
