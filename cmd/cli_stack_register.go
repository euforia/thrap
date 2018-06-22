package main

import (
	"context"
	"fmt"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

func commandStackRegister() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Register a new project",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := stack.Validate(); errs != nil {
				return utils.FlattenErrors(errs)
			}

			taddr := ctx.String("thrap-addr")
			if taddr == "" {
				return errRemoteRequired
			}
			cc, err := grpc.Dial(taddr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			tclient := thrapb.NewThrapClient(cc)
			resp, err := tclient.RegisterStack(context.Background(), stack)
			if err != nil {
				return err
			}

			m := map[string]*thrapb.Stack{
				"stack": resp,
			}

			b, _ := hclencoder.Encode(m)
			fmt.Printf("\n%s\n", b)
			return nil
		},
	}
}
