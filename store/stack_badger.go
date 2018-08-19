package store

import (
	"bytes"
	"errors"

	"github.com/dgraph-io/badger"
	"github.com/euforia/thrap/thrapb"
	"github.com/gogo/protobuf/proto"
)

const defaultStackPrefix = "/stack/"

var (
	// ErrStackExists is used when a stack already exists
	ErrStackExists = errors.New("stack exists")
	// ErrStackNotFound is used when an stack is not found in the store
	ErrStackNotFound = errors.New("stack not found")
)

// BadgerStackStorage implments a badger backed StackStorage interface
type BadgerStackStorage struct {
	db *badger.DB
}

// NewBadgerStackStorage returns a new badger backed stack storage
func NewBadgerStackStorage(db *badger.DB) *BadgerStackStorage {
	return &BadgerStackStorage{db: db}
}

func (store *BadgerStackStorage) getOpaqueKey(k string) []byte {
	return []byte(defaultStackPrefix + k)
}

// Get returns a stack by the id and the chain header.
func (store *BadgerStackStorage) Get(id string) (*thrapb.Stack, error) {
	var (
		key   = store.getOpaqueKey(id)
		stack *thrapb.Stack
	)

	err := store.db.View(func(txn *badger.Txn) error {
		var err error
		stack, err = store.getStack(txn, key)
		return err
	})

	return stack, err
}

// Create tries to write the new stack to the db.  If the stack exists if will
// return an error and abort registration
func (store *BadgerStackStorage) Create(stack *thrapb.Stack) (*thrapb.Stack, error) {
	key := store.getOpaqueKey(stack.ID)
	val, err := proto.Marshal(stack)
	if err != nil {
		return nil, err
	}

	err = store.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err == nil {
			return ErrStackExists
		}

		return txn.Set(key, val)
	})
	return stack, err
}

// Update updates an existing stack
func (store *BadgerStackStorage) Update(stack *thrapb.Stack) (*thrapb.Stack, error) {
	key := store.getOpaqueKey(stack.ID)
	val, err := proto.Marshal(stack)
	if err != nil {
		return nil, err
	}

	err = store.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return ErrStackNotFound
		}

		return txn.Set(key, val)
	})

	return stack, err
}

// List lists all stack ids.  It takes an option prefix to limit results to it
func (store *BadgerStackStorage) List(prefix string) ([]string, error) {
	pre := store.getOpaqueKey(prefix)
	out := make([]string, 0)

	err := store.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		iter := txn.NewIterator(opts)
		defer iter.Close()

		for iter.Seek(pre); iter.ValidForPrefix(pre); iter.Next() {
			item := iter.Item()
			key := item.Key()
			out = append(out, string(bytes.TrimSuffix(key, pre)))
		}

		return nil
	})
	return out, err
}

// Iter iterates over each stack from the starting point
func (store *BadgerStackStorage) Iter(start string, callback func(*thrapb.Stack) error) error {
	prefix := store.getOpaqueKey(start)

	return store.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		for iter.Seek(prefix); iter.ValidForPrefix(prefix); iter.Next() {
			item := iter.Item()
			stack, err := stackFromItem(item)
			if err != nil {
				return err
			}
			if err = callback(stack); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete removes the stack given the id
func (store *BadgerStackStorage) Delete(id string) (*thrapb.Stack, error) {
	var (
		key   = store.getOpaqueKey(id)
		stack *thrapb.Stack
	)

	err := store.db.Update(func(txn *badger.Txn) error {
		var err error
		stack, err = store.getStack(txn, key)
		if err != nil {
			return ErrStackNotFound
		}

		return txn.Delete(key)
	})

	return stack, err
}

func (store *BadgerStackStorage) getStack(txn *badger.Txn, key []byte) (*thrapb.Stack, error) {
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	return stackFromItem(item)
}

func stackFromItem(item *badger.Item) (*thrapb.Stack, error) {
	val, err := item.Value()
	if err != nil {
		return nil, err
	}

	var stack thrapb.Stack
	err = proto.Unmarshal(val, &stack)

	return &stack, err
}
