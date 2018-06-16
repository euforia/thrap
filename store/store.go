package store

import (
	"bytes"
	"errors"
	"hash"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/euforia/thrap/thrapb"
)

var (
	errPrevHashMismatch = errors.New("previous hash mismatch")
	errObjectExists     = errors.New("object exists")
	errRefExists        = errors.New("ref exists")
)

const (
	objKeyPrefix = "objects/"
	refKeyPrefix = "refs/"
)

// Object is a datastructure that is hashable and protobuf friendly
type Object interface {
	Hash(h hash.Hash) []byte
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
}

func NewBadgerDB(datadir string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = datadir
	opts.ValueDir = datadir
	opts.SyncWrites = true
	return badger.Open(opts)
}

type BadgerStore struct {
	db         *badger.DB
	prefix     []byte
	hf         func() hash.Hash
	zeroDigest []byte
}

// NewBadgerStore returns a badger backed manifest store
func NewBadgerStore(db *badger.DB, hf func() hash.Hash, prefix string) *BadgerStore {
	b := &BadgerStore{
		db:         db,
		prefix:     []byte(prefix),
		hf:         hf,
		zeroDigest: make([]byte, hf().Size()),
	}
	if b.prefix[len(b.prefix)-1] != '/' {
		b.prefix = append(b.prefix, '/')
	}

	return b
}

func (store *BadgerStore) CreateRef(ref string) ([]byte, *thrapb.ChainHeader, error) {
	h := store.hf()
	ch := &thrapb.ChainHeader{
		Previous:   make([]byte, h.Size()),
		Timestamp:  time.Now().UnixNano(),
		DataDigest: make([]byte, h.Size()),
	}
	b, _ := ch.Marshal()
	digest := ch.Hash(h)

	refkey := store.getRefOpaque([]byte(ref))

	er := store.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(refkey); err == nil {
			return errRefExists
		}
		err := txn.Set(store.getObjOpaque(digest), b)
		if err == nil {
			err = txn.Set(refkey, digest)
		}
		return err
	})

	return digest, ch, er
}

func (store *BadgerStore) IterRefChain(ref string, f func(header *thrapb.ChainHeader) error) error {
	refkey := store.getRefOpaque([]byte(ref))
	er := store.db.View(func(txn *badger.Txn) error {

		item, err := txn.Get(refkey)
		if err != nil {
			return err
		}
		digest, err := item.Value()
		if err != nil {
			return err
		}

		for err == nil && !store.isZeroDigest(digest) {
			var robj *thrapb.ChainHeader
			robj, err = store.getRefObj(txn, digest)
			if err == nil {
				if err = f(robj); err == nil {
					digest = robj.Previous
				}
			}

		}
		return err
	})

	return er
}

func (store *BadgerStore) isZeroDigest(in []byte) bool {
	return bytes.Compare(store.zeroDigest, in) == 0
}

func (store *BadgerStore) SetRef(ref string, ch *thrapb.ChainHeader) ([]byte, error) {

	refkey := store.getRefOpaque([]byte(ref))
	h := store.hf()

	var digest []byte
	err := store.db.Update(func(txn *badger.Txn) error {
		// Check specified previous exists
		refobj, er := store.getRefObj(txn, ch.Previous)
		if er != nil {
			return er
		}

		ch.Height = refobj.Height + 1
		b, er := ch.Marshal()
		if er != nil {
			return er
		}

		digest = ch.Hash(h)
		objkey := store.getObjOpaque(digest)

		er = txn.Set(objkey, b)
		if er == nil {
			er = txn.Set(refkey, digest)
		}
		return er
	})

	return digest, err
}

func (store *BadgerStore) GetRef(ref string) (*thrapb.ChainHeader, error) {
	var (
		refkey = store.getRefOpaque([]byte(ref))
		pb     *thrapb.ChainHeader
	)

	er := store.db.View(func(txn *badger.Txn) error {
		var err error
		pb, err = store.getRef(txn, refkey)
		return err
	})

	return pb, er
}

func (store *BadgerStore) getRef(txn *badger.Txn, refkey []byte) (*thrapb.ChainHeader, error) {
	item, err := txn.Get(refkey)
	if err != nil {
		return nil, err
	}
	digest, err := item.Value()
	if err != nil {
		return nil, err
	}

	return store.getRefObj(txn, digest)

}

func (store *BadgerStore) DeleteRef(ref string) error {
	key := store.getRefOpaque([]byte(ref))
	return store.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (store *BadgerStore) Get(digest []byte, out Object) error {
	key := store.getObjOpaque(digest)

	err := store.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.Value()
		if err == nil {
			err = out.Unmarshal(val)
		}

		return err
	})

	return err
}

func (store *BadgerStore) Set(obj Object) ([]byte, error) {
	h := store.hf()
	digest := obj.Hash(h)
	objkey := store.getObjOpaque(digest)

	err := store.db.Update(func(txn *badger.Txn) error {

		if _, err := txn.Get(objkey); err == nil {
			return errObjectExists
		}

		val, err := obj.Marshal()
		if err == nil {
			err = txn.Set(objkey, val)
		}

		return err
	})

	return digest, err
}

// Delete deletes an object by the given digest
func (store *BadgerStore) Delete(digest []byte) error {
	return store.db.Update(func(txn *badger.Txn) error {
		key := store.getObjOpaque(digest)
		return txn.Delete(key)
	})
}

func (store *BadgerStore) getRefOpaque(key []byte) []byte {
	return append(store.prefix, append([]byte(refKeyPrefix), key...)...)
}

func (store *BadgerStore) getObjOpaque(digest []byte) []byte {
	return append(store.prefix, append([]byte(objKeyPrefix), digest...)...)
}

func (store *BadgerStore) getRefObj(txn *badger.Txn, digest []byte) (*thrapb.ChainHeader, error) {
	rkey := store.getObjOpaque(digest)
	item, err := txn.Get(rkey)
	if err != nil {
		return nil, err
	}

	pb := new(thrapb.ChainHeader)

	val, err := item.Value()
	if err == nil {
		err = pb.Unmarshal(val)
	}

	return pb, err
}
