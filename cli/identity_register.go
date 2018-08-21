package cli

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/euforia/base58"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/pkg/config"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
)

func commandIdentityRegister() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Identity / user registration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "email `address`",
			},
			&cli.StringFlag{
				Name:    "code",
				Aliases: []string{"c"},
				Usage:   "registration `confirmation` code",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.String("email") == "" {
				return errors.New("email address required")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			cfile, err := homedir.Expand(filepath.Join(consts.DefaultDataDir, consts.ConfigFile))
			if err != nil {
				return err
			}
			conf, err := config.ReadConfig(cfile)
			if err != nil {
				return err
			}
			vcsp := conf.DefaultVCS()
			if vcsp.Username == "" {
				return errNotConfigured
			}

			kp, err := utils.LoadECDSAKeyPair(filepath.Join(consts.DefaultDataDir, "ecdsa256"))
			if err != nil {
				return errors.Wrap(err, "loading keypair")
			}
			pk := kp.PublicKey

			// Init and check identity
			ident := thrapb.NewIdentity(ctx.String("email"))
			ident.PublicKey = append(pk.X.Bytes(), pk.Y.Bytes()...)
			ident.Meta = map[string]string{
				vcsp.ID + ".username": vcsp.Username,
			}
			err = ident.Validate()
			if err != nil {
				return err
			}

			tclient, err := newThrapClient(ctx)
			if err != nil {
				return err
			}

			// Confirm registration request
			if confirmCode, ok := isConfirmRequest(ctx); ok {

				idt, err := confirmUserRegistration(tclient, kp, ident, confirmCode)
				if err == nil {
					b, _ := json.MarshalIndent(idt, "", "  ")
					fmt.Printf("%s\n", b)
				}
				return err

			}

			// Submit registration request
			resp, err := tclient.RegisterIdentity(context.Background(), ident)
			if err != nil {
				return err
			}

			// Generate code
			h := sha256.New()
			sh := resp.SigHash(h)
			code := base58.Encode(sh)
			fmt.Printf("Code: %s\n", code)
			return nil
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

	id, err := cc.ConfirmIdentity(context.Background(), ident)

	return id, err
}

func isConfirmRequest(ctx *cli.Context) (string, bool) {
	code := ctx.String("code")
	return code, len(code) > 0
}
