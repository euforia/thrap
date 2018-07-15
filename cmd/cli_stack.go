package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStack() *cli.Command {
	return &cli.Command{
		Name:  "stack",
		Usage: "Stack operations",
		Subcommands: []*cli.Command{
			commandStackInit(),
			commandStackBuild(),
			commandStackDeploy(),
			commandStackStatus(),
			commandStackLogs(),
			commandStackRegister(),
			// commandStackValidate(),
			commandStackStop(),
			commandStackDestroy(),
			commandStackVersion(),
		},
	}
}

// func commandStackValidate() *cli.Command {
// 	return &cli.Command{
// 		Name:      "validate",
// 		Usage:     "Validate a manifest",
// 		ArgsUsage: "[path to manifest]",
// 		Action: func(ctx *cli.Context) error {

// 			mfile := ctx.Args().Get(0)
// 			mf, err := manifest.LoadManifest(mfile)
// 			if err != nil {
// 				return err
// 			}

// 			core, err := core.NewCore(&core.Config{PacksDir: consts.DefaultPacksDir})
// 			if err != nil {
// 				return err
// 			}

// 			stm := core.Stack()
// 			err = stm.Validate(mf)
// 			if err == nil {
// 				writeHCLManifest(mf, os.Stdout)
// 			}

// 			return err
// 		},
// 	}
// }

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
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			core, err := core.NewCore(&core.Config{DataDir: consts.DefaultDataDir})
			if err != nil {
				return err
			}

			stm := core.Stack()

			return stm.Build(context.Background(), stack)
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

			core, err := core.NewCore(&core.Config{DataDir: consts.DefaultDataDir})
			if err != nil {
				return err
			}

			stm := core.Stack()

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

			core, err := core.NewCore(&core.Config{DataDir: consts.DefaultDataDir})
			if err != nil {
				return err
			}

			stm := core.Stack()

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

			// fmt.Println("Version:", stack.Version)

			core, err := core.NewCore(&core.Config{DataDir: consts.DefaultDataDir})
			if err != nil {
				return err
			}

			stm := core.Stack()
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

			return nil
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

			core, err := core.NewCore(&core.Config{DataDir: consts.DefaultDataDir})
			if err != nil {
				return err
			}

			stm := core.Stack()

			cid := ctx.Args().Get(0)
			c := context.Background()
			if cid == "" {
				return stm.Logs(c, stack, os.Stdout, os.Stderr)
			}

			return stm.Log(c, cid+"."+stack.ID, os.Stdout, os.Stderr)
		},
	}
}
