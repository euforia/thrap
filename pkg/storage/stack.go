package storage

import (
	"github.com/euforia/kvdb"
	"github.com/euforia/thrap/thrapb"
)

// KvdbStackStorage implments a badger backed StackStorage interface
type KvdbStackStorage struct {
	table kvdb.Table
}

// Get returns a stack by the id and the chain header.
func (store *KvdbStackStorage) Get(id string) (*thrapb.Stack, error) {
	obj, err := store.table.Get([]byte(id))
	if err == nil {
		if stack, ok := obj.(*thrapb.Stack); ok {
			return stack, nil
		}
		err = ErrInvalidType
	}
	return nil, err
}

// Create tries to write the new stack to the db.  If the stack exists if will
// return an error and abort registration
func (store *KvdbStackStorage) Create(stack *thrapb.Stack) error {
	return store.table.Create([]byte(stack.ID), stack)
}

// Update updates an existing stack
func (store *KvdbStackStorage) Update(stack *thrapb.Stack) error {
	return store.table.Update([]byte(stack.ID), stack)
}

// Iter iterates over each stack from the starting point
func (store *KvdbStackStorage) Iter(start string, callback func(*thrapb.Stack) error) error {
	var prefix []byte
	if start != "" {
		prefix = []byte(start)
	}
	return store.table.Iter(prefix, func(obj kvdb.Object) error {
		proj, ok := obj.(*thrapb.Stack)
		if ok {
			return callback(proj)
		}
		return ErrInvalidType
	})
}

// Delete removes the stack given the id
func (store *KvdbStackStorage) Delete(id string) error {
	return store.table.Delete([]byte(id))
}
