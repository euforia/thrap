package cli

import (
	"context"
	"io"

	"github.com/euforia/thrap/thrapb"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func commandIdentity() *cli.Command {
	return &cli.Command{
		Name:  "identity",
		Usage: "Identity operations",
		Subcommands: []*cli.Command{
			commandIdentityRegister(),
			commandIdentityShow(),
			commandIdentityList(),
		},
	}
}

func commandIdentityShow() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show identity details",
		Action: func(ctx *cli.Context) error {
			identID := ctx.Args().Get(0)
			if len(identID) == 0 {
				return errors.New("identity required")
			}

			tclient, err := newThrapClient(ctx)
			if err != nil {
				return err
			}

			ident, err := tclient.GetIdentity(context.Background(), &thrapb.Identity{ID: identID})
			if err != nil {
				return err
			}

			writeJSON(ident)
			return nil
		},
	}
}

func commandIdentityList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List identities",
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

			stream, err := tclient.IterIdentities(context.Background(), &thrapb.IterOptions{
				Prefix: ctx.String("prefix"),
			})
			if err != nil {
				return err
			}

			for {
				ident, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						return stream.CloseSend()
					}
					return err
				}

				writeJSON(ident)
			}
		},
	}
}
