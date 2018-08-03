package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStack() *cli.Command {
	return &cli.Command{
		Name:  "stack",
		Usage: "Stack operations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "`profile` to use",
				Value:   "local",
			},
		},
		Subcommands: []*cli.Command{
			commandStackList(),
			commandStackInit(),
			commandStackRegister(),
			commandStackEnsure(),
			commandStackCommit(),
			commandStackBuild(),
			commandStackArtifacts(),
			commandStackDeploy(),
			commandStackStatus(),
			commandStackLogs(),
			commandStackStop(),
			commandStackDestroy(),
			commandStackVersion(),
			commandProfile(),
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
				Name:  "pub",
				Usage: "publish artifacts",
			},
		},
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

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

			cr, err := loadCore(ctx)
			if err != nil {
				return err
			}

			stm, err := cr.Stack(prof)
			if err != nil {
				return err
			}

			// lpath, _ := utils.GetLocalPath("")
			opt := core.BuildOptions{
				Workdir: lpath,
				Publish: ctx.Bool("pub"),
			}

			return stm.Build(context.Background(), stack, opt)
		},
	}
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

			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
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

func commandStackEnsure() *cli.Command {
	return &cli.Command{
		Name:  "ensure",
		Usage: "Ensure resources exist",
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

			st, err := cr.Stack(prof)
			if err != nil {
				return err
			}

			results := st.EnsureResources(stack)
			results.Print(os.Stdout)

			return nil
		},
	}
}

func defaultPrintStackResults(results []*thrapb.ActionResult) {
	fmt.Println()
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, " \tComponent\tStatus\tDetails\n")
	fmt.Fprintf(tw, " \t---------\t------\t-------\n")
	for _, r := range results {
		if r.Error == nil {
			fmt.Fprintf(tw, " \t%s\tsucceeded\t \n", r.Resource)
		} else {
			fmt.Fprintf(tw, " \t%s\tfailed\t%v\n", r.Resource, r.Error)
		}
	}
	tw.Flush()
	fmt.Println()
}
