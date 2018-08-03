package cli

import (
	"context"

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

			report := stm.Stop(context.Background(), stack)
			defaultPrintStackResults(report)

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

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(thrapb.DefaultProfile())
			if err != nil {
				return err
			}

			report := stm.Destroy(context.Background(), stack)
			defaultPrintStackResults(report)

			return nil
		},
	}
}
