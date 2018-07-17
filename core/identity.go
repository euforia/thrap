package core

import (
	"crypto/sha256"
	"log"
	"math/rand"

	"github.com/euforia/base58"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/pkg/errors"
)

var (
	// ErrIdentityAlreadySigned is used when an identity is already signed
	ErrIdentityAlreadySigned = errors.New("identity already signed")
	// ErrIdentityAlreadyRegistered is used when identity has already been registered
	ErrIdentityAlreadyRegistered = errors.New("identity already registered")
)

// Identity is the caninical interface to interact with identities
type Identity struct {
	store IdentityStorage
	log   *log.Logger
}

// Confirm confirms a identity registration request and completes it.
// In this case the public key field is the signature from the client
func (idt *Identity) Confirm(ident *thrapb.Identity) (*thrapb.Identity, error) {

	sident, err := idt.store.Get(ident.ID)
	if err != nil {
		return nil, err
	}

	if len(sident.Signature) > 0 {
		return nil, ErrIdentityAlreadySigned
	}

	shash := sident.SigHash(sha256.New())

	b58e := base58.Encode(shash)
	idt.log.Printf("Verifying user registration code=%s", b58e)

	if !utils.VerifySignature(ident.PublicKey, shash, ident.Signature) {
		return nil, errors.New("signature verification failed")
	}

	sident.Signature = ident.Signature

	_, err = idt.store.Update(sident)
	if err == nil {
		idt.log.Printf("User registered user=%s", sident.ID)
	}
	return sident, err
}

// Register registers a new identity. It returns an error if the identity exists
// or fails to register
func (idt *Identity) Register(ident *thrapb.Identity) (*thrapb.Identity, []*ActionReport, error) {
	err := ident.Validate()
	if err != nil {
		return nil, nil, err
	}
	ident.Nonce = rand.Uint64()

	er := &ActionReport{}
	er.Data, er.Error = idt.store.Create(ident)
	if er.Error == store.ErrRefExists {
		er.Error = ErrIdentityAlreadyRegistered
	}
	er.Action = NewAction("create", "identity", ident.ID)
	if err == nil {
		idt.log.Printf("User registration request user=%s", ident.ID)
	}

	return ident, []*ActionReport{er}, er.Error
}

// Get returns an Identity by the id
func (idt *Identity) Get(id string) (*thrapb.Identity, error) {
	ident, err := idt.store.Get(id)
	return ident, err
}

func (idt *Identity) Iter(prefix string, f func(*thrapb.Identity) error) error {
	return idt.store.Iter(prefix, f)
}
