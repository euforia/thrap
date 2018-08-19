package identity

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"strings"

	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/pkg/storage"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/pkg/errors"
)

var (
	// ErrDomainNotInWhitelist is used when a user email domain is not whitelisted
	ErrDomainNotInWhitelist = errors.New("domain not in whitelist")
	// ErrIdentityAlreadySigned is used when an identity is already signed
	ErrIdentityAlreadySigned = errors.New("identity already signed")
	// ErrIdentityAlreadyRegistered is used when identity has already been registered
	ErrIdentityAlreadyRegistered = errors.New("identity already registered")
)

// Config is the config used to initialize the Identity manager
type Config struct {
	Whitelist []string
	Storage   storage.IdentityStorage
	Logger    *log.Logger
}

// Identity is the canonical interface to interact with identities
type Identity struct {
	// identity store
	store storage.IdentityStorage
	// enable email domain whitelisting
	whitelistEnabled bool
	// allowed email domains
	whitelist []string
	// user supplied logger
	log *log.Logger
}

func New(conf *Config) *Identity {
	ident := &Identity{
		store:            conf.Storage,
		whitelistEnabled: len(conf.Whitelist) > 0,
		whitelist:        conf.Whitelist,
		log:              conf.Logger,
	}

	// Discard logs if logger not supplied
	if ident.log == nil {
		ident.log = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return ident
}

func (idt *Identity) isWhitelisted(email string) bool {
	if !idt.whitelistEnabled {
		return true
	}
	for _, v := range idt.whitelist {
		if strings.HasSuffix(email, v) {
			return true
		}
	}
	return false
}

// Register registers a new identity. It returns an error if the identity exists
// or fails to register otherwise the created Identity and keypair are returned
func (idt *Identity) Register(ident *thrapb.Identity) (*thrapb.Identity, *ecdsa.PrivateKey, error) {
	kp, err := utils.GenerateECDSAKeyPair()
	if err != nil {
		return nil, nil, err
	}

	ident.PublicKey = append(kp.PublicKey.X.Bytes(), kp.PublicKey.Y.Bytes()...)

	err = ident.Validate()
	if err != nil {
		return nil, nil, err
	}

	if !idt.isWhitelisted(ident.Email) {
		return nil, nil, ErrDomainNotInWhitelist
	}

	ident.Nonce = mrand.Uint64()

	// Sign
	code := ident.SigHash(sha256.New())
	r, s, err := ecdsa.Sign(crand.Reader, kp, code)
	if err != nil {
		return nil, nil, err
	}
	ident.Signature = append(r.Bytes(), s.Bytes()...)

	// Add to store
	err = idt.store.Create(ident)
	if err == nil {
		idt.log.Printf("User registration request user=%s", ident.ID)
		return ident, kp, nil
	}

	if err == kvdb.ErrExists {
		err = ErrIdentityAlreadyRegistered
	}

	return ident, kp, err
}

// Get returns an Identity by the id
func (idt *Identity) Get(id string) (*thrapb.Identity, error) {
	ident, err := idt.store.Get(id)
	return ident, err
}

func (idt *Identity) Delete(id string) error {
	return idt.store.Delete(id)
}

// Iter iterates over each identity with the matching prefix
func (idt *Identity) Iter(prefix string, f func(*thrapb.Identity) error) error {
	return idt.store.Iter(prefix, f)
}
