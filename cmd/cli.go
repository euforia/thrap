package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"path/filepath"

	"github.com/euforia/thrap/packs"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/vars"

	"github.com/pkg/errors"

	"github.com/euforia/base58"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/thrapb"
	"gopkg.in/urfave/cli.v2"
)

var (
	errRemoteRequired = errors.New("thrap remote required")
)

func newCLI() *cli.App {
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println(ctx.App.Version)
	}

	app := &cli.App{
		Name:     "thrap",
		HelpName: "thrap",
		Version:  version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "thrap-addr",
				Usage:   "thrap registry address",
				EnvVars: []string{"THRAP_ADDR"},
			},
		},
		Commands: []*cli.Command{
			commandConfigure(),
			commandAgent(),
			//commandRegister(),
			commandStack(),
			commandPack(),
			commandVersion(),
		},
	}

	app.HideVersion = true

	return app
}

func commandPack() *cli.Command {
	return &cli.Command{
		Name:  "pack",
		Usage: "Pack operations",
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:  "update",
				Usage: "Update packs",
				Action: func(ctx *cli.Context) error {
					packdir := filepath.Join(consts.DefaultDataDir, consts.PacksDir)
					pks, err := packs.New(packdir)
					if err != nil {
						return err
					}

					return pks.Update()
				},
			},
		},
	}
}

func commandVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(ctx *cli.Context) error {
			fmt.Println(version())
			return nil
		},
	}
}

func commandConfigure() *cli.Command {
	return &cli.Command{
		Name:  "configure",
		Usage: "Configure global settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   vars.VcsID,
				Usage:  "version control `provider`",
				Value:  "github",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:  vars.VcsUsername,
				Usage: "version control `username`",
			},
			&cli.StringFlag{
				Name:  "data-dir",
				Usage: "data `directory`",
				Value: "~/" + consts.WorkDir,
			},
			&cli.BoolFlag{
				Name:  "no-prompt",
				Usage: "do not prompt for input",
			},
		},
		Action: func(ctx *cli.Context) error {
			opts := core.ConfigureOptions{
				VCS: &config.VCSConfig{
					ID:       ctx.String(vars.VcsID),
					Username: ctx.String(vars.VcsUsername),
				},
				DataDir:  ctx.String("data-dir"),
				NoPrompt: ctx.Bool("no-prompt"),
			}

			// Only configures things that are not configured
			return core.ConfigureGlobal(opts)
		},
	}
}

var errNotConfigured = errors.New("thrap not configured. Try running 'thrap configure'")

func commandRegister() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "User registration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "email address",
			},
			&cli.StringFlag{
				Name:    "code",
				Aliases: []string{"c"},
				Usage:   "registration confirmation code",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.String("email") == "" {
				return errors.New("email address required")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			// Load global config
			// gconf, err := thrap.ReadGlobalConfig()
			// if err != nil {
			// 	return err
			// }
			// if gconf.VCS.Username == "" {
			// 	return errNotConfigured
			// }
			//
			// kp, err := thrap.LoadUserKeyPair()
			// if err != nil {
			// 	return errors.Wrap(err, "loading keypair")
			// }
			// pk := kp.PublicKey
			//
			// // Init and check identity
			// ident := thrapb.NewIdentity(ctx.String("email"))
			// ident.PublicKey = append(pk.X.Bytes(), pk.Y.Bytes()...)
			// ident.Meta = map[string]string{
			// 	gconf.VCS.ID + ".username": gconf.VCS.Username,
			// }
			// err = ident.Validate()
			// if err != nil {
			// 	return err
			// }
			//
			// // Check remote addr
			// remoteAddr := ctx.String("thrap-addr")
			// if remoteAddr == "" {
			// 	return errRemoteRequired
			// }
			// cc, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
			// if err != nil {
			// 	return err
			// }
			//
			// tclient := thrapb.NewThrapClient(cc)
			//
			// // Confirm registration request
			// confirmCode := ctx.String("code")
			// if len(confirmCode) > 0 {
			// 	_, err = confirmUserRegistration(tclient, kp, ident, confirmCode)
			// 	if err == nil {
			// 		fmt.Println("Registered!")
			// 	}
			// 	return err
			// }
			//
			// // Submit registration request
			// resp, err := tclient.RegisterIdentity(context.Background(), ident)
			// if err != nil {
			// 	return err
			// }
			//
			// // Generate code
			// h := sha256.New()
			// sh := resp.SigHash(h)
			// out := base58.Encode(sh)
			//
			// fmt.Printf("%s\n", out)

			return errors.New("to be implemented")
		},
	}
}

func confirmUserRegistration(cc thrapb.ThrapClient, kp *ecdsa.PrivateKey, ident *thrapb.Identity, confirmCode string) (*thrapb.Identity, error) {

	code := base58.Decode([]byte(confirmCode))
	r, s, err := ecdsa.Sign(rand.Reader, kp, code)
	if err != nil {
		return nil, err
	}
	ident.Signature = append(r.Bytes(), s.Bytes()...)

	return cc.ConfirmIdentity(context.Background(), ident)
}
