package main

import (
	"context"
	"fmt"

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
			// commandStackRegister(),
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

			core, err := core.NewCore(&core.Config{PacksDir: consts.DefaultPacksDir})
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

			core, err := core.NewCore(&core.Config{PacksDir: consts.DefaultPacksDir})
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

			core, err := core.NewCore(&core.Config{PacksDir: consts.DefaultPacksDir})
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
