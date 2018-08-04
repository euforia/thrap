package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackStop() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop stack components",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := stack.Validate(); len(errs) > 0 {
				return utils.FlattenErrors(errs)
			}

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(thrapb.DefaultProfile())
			if err != nil {
				return err
			}

			var stop bool
			utils.PromptUntilNoError("Are you sure you want to stop "+stack.ID+" [y/N] ? ",
				os.Stdout, os.Stdin, func(in []byte) error {
					s := string(in)
					switch s {
					case "y", "Y", "yes", "Yes":
						stop = true
					}
					return nil
				})

			if stop {
				report := stm.Stop(context.Background(), stack)
				defaultPrintStackResults(report)
			} else {
				fmt.Println("Exiting!")
			}

			return nil
		},
	}
}

func commandStackDestroy() *cli.Command {
	return &cli.Command{
		Name:  "destroy",
		Usage: "Destroy stack components",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := stack.Validate(); len(errs) > 0 {
				return utils.FlattenErrors(errs)
			}

			_, prof, err := loadProfile(ctx)
			if err != nil {
				return err
			}

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(prof)
			if err != nil {
				return err
			}

			var destroy bool
			utils.PromptUntilNoError("Are you sure you want to destroy "+stack.ID+" [y/N] ? ",
				os.Stdout, os.Stdin, func(in []byte) error {
					s := string(in)
					switch s {
					case "y", "Y", "yes", "Yes":
						destroy = true
					}
					return nil
				})

			if destroy {
				report := stm.Destroy(context.Background(), stack)
				defaultPrintStackResults(report)
			} else {
				fmt.Println("Exiting!")
			}

			return nil
		},
	}
}
