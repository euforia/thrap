package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStack() *cli.Command {
	return &cli.Command{
		Name:  "stack",
		Usage: "Stack operations",
		Subcommands: []*cli.Command{
			commandStackList(),
			commandStackInit(),
			commandStackRegister(),
			commandStackCommit(),
			commandStackBuild(),
			commandStackArtifacts(),
			commandStackDeploy(),
			commandStackStatus(),
			commandStackLogs(),
			commandStackStop(),
			commandStackDestroy(),
			commandStackVersion(),
		},
	}
}

func commandStackVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show stack version",
		Action: func(ctx *cli.Context) error {
			stack, err := manifest.LoadManifest("")
			if err == nil {
				fmt.Println(stack.Version)
			}

			return err
		},
	}
}

func commandStackBuild() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build stack components",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "pub",
				Aliases: []string{"p"},
				Usage:   "publish artifacts",
			},
		},
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(thrapb.DefaultProfile())
			if err != nil {
				return err
			}

			dd, _ := utils.GetLocalPath("")
			opt := core.BuildOptions{
				Workdir: dd,
				Publish: ctx.Bool("pub"),
			}

			return stm.Build(context.Background(), stack, opt)
		},
	}
}

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
			for _, r := range report {
				if r.Error == nil {
					fmt.Println(r.Action.String())
				} else {
					fmt.Println(r.Action.String(), r.Error)
				}
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

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(thrapb.DefaultProfile())
			if err != nil {
				return err
			}

			report := stm.Destroy(context.Background(), stack)
			for _, r := range report {
				if r.Error != nil {
					fmt.Println(r.Error)
				}
			}

			return nil
		},
	}
}

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

			stm, err := cr.Stack(thrapb.DefaultProfile())
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

			stm, err := cr.Stack(thrapb.DefaultProfile())
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

func printStackStatus(stm *core.Stack, stack *thrapb.Stack) {
	resp := stm.Status(context.Background(), stack)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Component\tImage\tStatus\tDetails\n")
	fmt.Fprintf(tw, "---------\t-----\t------\t-------\n")
	for _, s := range resp {

		d := s.Details
		st := d.State

		if s.Error != nil {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.ID, d.Config.Image, st.Status, s.Error)
		} else {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", s.ID, d.Config.Image, st.Status, d.NetworkSettings.Ports)
		}

	}
	tw.Flush()
}

func commandStackLogs() *cli.Command {
	return &cli.Command{
		Name:      "logs",
		Usage:     "Show stack runtime logs",
		ArgsUsage: "[component]",
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

			cid := ctx.Args().Get(0)
			c := context.Background()
			if cid == "" {
				return stm.Logs(c, stack, os.Stdout, os.Stderr)
			}

			return stm.Log(c, cid+"."+stack.ID, os.Stdout, os.Stderr)
		},
	}
}

func commandStackList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all stacks",
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "filter by `prefix`",
			},
		},
		Action: func(ctx *cli.Context) error {
			tclient, err := newThrapClient(ctx)
			if err != nil {
				return err
			}

			stream, err := tclient.IterStacks(context.Background(), &thrapb.IterOptions{
				Prefix: ctx.String("prefix"),
			})
			if err != nil {
				return err
			}

			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)
			fmt.Fprintf(tw, "ID\tVERSION\n")
			for {
				stack, err := stream.Recv()
				if err != nil {
					defer tw.Flush()
					if err == io.EOF {
						return stream.CloseSend()
					}
					return err
				}

				fmt.Fprintf(tw, "%s\t%s\n", stack.ID, stack.Version)
			}
		},
	}
}
