package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackStatus() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show status",
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

			stm, err := cr.Stack(pb.DefaultProfile())
			if err != nil {
				return err
			}

			fmt.Println()
			printStackStatus(stm, stack)
			fmt.Println()

			return nil
		},
	}
}

func printStackStatus(stm *core.Stack, stack *thrapb.Stack) {
	resp := stm.Status(context.Background(), stack)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Component\tImage\tStatus\tDetails\n")
	fmt.Fprintf(tw, "---------\t-----\t------\t-------\n")
	for _, s := range resp {

		// d := s.Details
		// st := d.State

		if s.Error != nil {
			// fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.ID, d.Config.Image, s.Status, s.Error)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.ID, "", s.Status, s.Error)
		} else {
			// fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", s.ID, d.Config.Image, s.Status, s.Details)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", s.ID, "", s.Status, s.Details)
		}

	}
	tw.Flush()
}
