package main

import (
	"fmt"
	"os"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/orchestrator"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackBuild() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build the stack",
		Action: func(ctx *cli.Context) error {

			st, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := st.Validate(); errs != nil {
				return utils.FlattenErrors(errs)
			}

			oconf := &orchestrator.Config{Provider: "docker"}
			orch, err := orchestrator.New(oconf)
			if err != nil {
				return err
			}

			bldr := orch.(*orchestrator.DockerOrchestrator)

			for _, comp := range st.Components {
				if !comp.IsBuildable() {
					continue
				}
				if comp.Build.Context == "" {
					comp.Build.Context = consts.DefaultBuildContext
				}

				fmt.Printf("Building component.%s\n", comp.ID)
				err = bldr.Build(comp, os.Stdout)
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}
