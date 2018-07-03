package main

import (
	"context"
	"fmt"
	"os"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"gopkg.in/urfave/cli.v2"
)

const defaultPacksDir = "~/.thrap/packs"

func commandStack() *cli.Command {
	return &cli.Command{
		Name:  "stack",
		Usage: "Stack operations",
		Subcommands: []*cli.Command{
			commandStackBuild(),
			commandStackDeploy(),
			commandStackInit(),
			// commandStackRegister(),
			commandStackValidate(),
			commandStackVersion(),
		},
	}
}

func commandStackValidate() *cli.Command {
	return &cli.Command{
		Name:      "validate",
		Usage:     "Validate a manifest",
		ArgsUsage: "[path to manifest]",
		Action: func(ctx *cli.Context) error {

			mfile := ctx.Args().Get(0)
			mf, err := manifest.LoadManifest(mfile)
			if err != nil {
				return err
			}

			core, err := core.NewCore(&core.Config{PacksDir: defaultPacksDir})
			if err != nil {
				return err
			}

			stm := core.Stack()
			err = stm.Validate(mf)
			if err == nil {
				// return utils.FlattenErrors(errs)
				writeHCLManifest(mf, os.Stdout)
			}

			return err
		},
	}
}

func commandStackVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show stack version",
		Action: func(ctx *cli.Context) error {
			// lpath, err := utils.GetLocalPath("")
			// if err == nil {
			// 	fmt.Println(vcs.GetRepoVersion(lpath))
			// }
			// return err
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
		Usage: "Build the stack",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			core, err := core.NewCore(&core.Config{PacksDir: defaultPacksDir})
			if err != nil {
				return err
			}

			stm := core.Stack()
			return stm.Build(context.Background(), stack)
		},
	}
}
