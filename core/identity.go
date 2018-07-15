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

// Identity is the caninical interface to interact with identities
type Identity struct {
	store *store.IdentityStore // local store
	log   *log.Logger
}

// Confirm confirms a identity registration request and completes it.
// In this case the public key field is the signature from the client
func (idt *Identity) Confirm(ident *thrapb.Identity) (*thrapb.Identity, error) {

	sident, _, err := idt.store.Get(ident.ID)
	if err != nil {
		return nil, err
	}

	shash := sident.SigHash(sha256.New())

	b58e := base58.Encode(shash)
	idt.log.Printf("Verifying user registration code=%s", b58e)

	if !utils.VerifySignature(ident.PublicKey, shash, ident.Signature) {
		return nil, errors.New("signature verification failed")
	}

	sident.Signature = ident.Signature

	resp, _, err := idt.store.Update(sident)
	if err == nil {
		idt.log.Printf("User registered user=%s", ident.ID)
	}
	return resp, err
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
	er.Data, _, er.Error = idt.store.Create(ident)
	er.Action = NewAction("create", "identity", ident.ID)
	if err == nil {
		idt.log.Printf("User registration request user=%s", ident.ID)
	}

	return ident, []*ActionReport{er}, er.Error
}

// Get returns an Identity by the id
func (idt *Identity) Get(id string) (*thrapb.Identity, error) {
	ident, _, err := idt.store.Get(id)
	return ident, err
}
