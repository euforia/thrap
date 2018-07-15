package thrapb

import (
	"encoding/binary"
	"errors"
	"hash"
	"math/rand"
	"sort"
	"strings"
)

// NewIdentity returns a new Identity with given email address
func NewIdentity(email string) *Identity {
	return &Identity{
		Email: email,
		Nonce: rand.Uint64(),
	}
}

// SigHash returns the hash to be used to sign the object. This is everything
// except the signature itself
func (ident *Identity) SigHash(h hash.Hash) []byte {
	ident.hash(h)
	return h.Sum(nil)
}

// Hash returns the hash for the identity object
func (ident *Identity) Hash(h hash.Hash) []byte {
	ident.hash(h)
	h.Write(ident.Signature)
	return h.Sum(nil)
}

// Hash of the whole object. does not perform the final Sum(nil)
func (ident *Identity) hash(h hash.Hash) {
	binary.Write(h, binary.BigEndian, ident.Nonce)
	h.Write([]byte(ident.ID))
	h.Write([]byte(ident.Email))
	h.Write(ident.PublicKey)

	keys := make([]string, 0, len(ident.Meta))
	for k := range ident.Meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(ident.Meta[k]))
	}

	//h.Write(ident.Signature)

	//return h.Sum(nil)
}

// Validate validates the identity fields.  If the id is not set, it defaults
// to the email address
func (ident *Identity) Validate() error {
	if len(ident.Email) == 0 {
		return errors.New("email missing")
	}
	if !strings.Contains(ident.Email, "@") {
		return errors.New("email invalid")
	}
	if len(ident.PublicKey) == 0 {
		return errors.New("public key missing")
	}

	if ident.ID == "" {
		ident.ID = ident.Email
	}

	return nil
}
