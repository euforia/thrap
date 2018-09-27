package cli

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v2"

	"github.com/euforia/thrap/pkg/api"
	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/thrap"
)

func commandAgent() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Run a server agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bind-addr",
				Aliases: []string{"b"},
				Usage:   "bind address",
				Value:   "0.0.0.0:10000",
			},
			&cli.StringFlag{
				Name:  "conf-dir",
				Usage: "configuration directory",
				Value: "",
			},
		},
		Action: func(ctx *cli.Context) error {
			conf, err := buildConfig(ctx)
			if err != nil {
				return err
			}

			credsFile := filepath.Join(conf.ConfigDir, "creds.hcl")
			conf.Credentials, err = credentials.ReadCredentials(credsFile)
			if err != nil {
				return err
			}

			thp, err := thrap.New(conf)
			if err != nil {
				return err
			}

			srv := api.NewServer(thp, conf.Logger)

			baddr := ctx.String("bind-addr")
			lis, err := net.Listen("tcp", baddr)
			if err != nil {
				return err
			}
			conf.Logger.Println("Starting server:", lis.Addr().String())

			err = srv.Serve(lis)
			return err
		},
	}
}

func buildConfig(ctx *cli.Context) (*thrap.Config, error) {
	conf := &thrap.Config{
		ConfigDir: ctx.String("conf-dir"),
		Logger:    log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds),
	}
	return conf, conf.Validate()
}
