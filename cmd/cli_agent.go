package main

import (
	"log"
	"net"

	"github.com/euforia/thrap"
	"github.com/euforia/thrap/thrapb"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

func commandAgent() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Run a server agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "bind-addr",
				Usage: "bind address",
			},
			&cli.StringFlag{
				Name:  "adv-addr",
				Usage: "advertise address",
			},
		},
		Action: func(ctx *cli.Context) error {
			core, err := thrap.NewCore(nil)
			if err != nil {
				return err
			}

			srv := grpc.NewServer()
			svc := thrap.NewService(core)
			thrapb.RegisterThrapServer(srv, svc)

			baddr := ctx.String("bind-addr")
			lis, err := net.Listen("tcp", baddr)
			if err != nil {
				return err
			}
			log.Println("Starting server:", lis.Addr().String())

			err = srv.Serve(lis)
			return err
		},
	}
}
