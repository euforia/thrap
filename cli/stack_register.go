package cli

import (
	"context"
	"fmt"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandStackRegister() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Register a new stack",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := stack.Validate(); errs != nil {
				return utils.FlattenErrors(errs)
			}

			var st *thrapb.Stack

			// Remote
			raddr := ctx.String("thrap-addr")
			if raddr != "" {
				fmt.Println("Registering with", raddr)
				tclient, err := newThrapClient(ctx)
				if err != nil {
					return err
				}

				st, err = tclient.RegisterStack(context.Background(), stack)
				if err != nil {
					return err
				}

			} else {
				fmt.Println("Registering locally")
				cr, err := loadCore(ctx)
				if err != nil {
					return err
				}

				stk, err := cr.Stack(thrapb.DefaultProfile())
				if err != nil {
					return err
				}

				st, _, err = stk.Register(stack)
				if err != nil {
					return err
				}

			}

			m := map[string]*thrapb.Stack{
				"stack": st,
			}

			b, _ := hclencoder.Encode(m)
			fmt.Printf("\n%s\n", b)
			return nil

		},
	}
}

func commandStackCommit() *cli.Command {
	return &cli.Command{
		Name:  "commit",
		Usage: "Commit stack definition",
		Action: func(ctx *cli.Context) error {

			stack, err := manifest.LoadManifest("")
			if err != nil {
				return err
			}

			tclient, err := newThrapClient(ctx)
			if err != nil {
				return err
			}

			resp, err := tclient.CommitStack(context.Background(), stack)
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
