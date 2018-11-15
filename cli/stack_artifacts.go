package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackArtifacts() *cli.Command {
	return &cli.Command{
		Name:    "artifacts",
		Aliases: []string{"art"},
		Usage:   "List stack artifacts",
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
			printStackArtifacts(stm, stack)
			fmt.Println()

			return nil
		},
	}
}

func printStackArtifacts(stm *core.Stack, stack *thrapb.Stack) {
	imgs := stm.Artifacts(stack)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Name\tID\tCreated\tSize\n")
	fmt.Fprintf(tw, "----\t--\t-------\t----\n")
	for _, img := range imgs {
		for _, tag := range img.Tags {
			d := time.Now().Sub(time.Unix(img.Created, 0)).Round(time.Second)
			smb := img.DataSize / (1024 * 1024)
			fmt.Fprintf(tw, "%s\t%s\t%s ago\t%d MB\n", tag, img.ID.Hex()[:12], d, smb)
		}

	}
	tw.Flush()
}
