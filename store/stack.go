package store

import (
	"errors"
	"time"

	"github.com/euforia/thrap/thrapb"
)

// StackStore is the stack storage interface
type StackStore struct {
	store ObjectStorage
}

// NewStackStore returns a StackStore instance
func NewStackStore(objs ObjectStorage) (*StackStore, error) {
	tp := &StackStore{
		store: objs,
	}

	var err error
	if tp.store == nil {
		err = errors.New("storage not provided")
	}

	return tp, err
}

// Get returns a stack by the id and the chain header.
func (thrap *StackStore) Get(id string) (*thrapb.Stack, *thrapb.ChainHeader, error) {
	ref, _, err := thrap.store.GetRef(id, "latest")
	if err != nil {
		return nil, nil, err
	}

	var stack thrapb.Stack
	err = thrap.store.Get(id, ref.DataDigest, &stack)

	return &stack, ref, err
}

// Create tries to write the new stack to the db.  If the stack exists if will
// return an error and abort registration
func (thrap *StackStore) Create(stack *thrapb.Stack) (*thrapb.Stack, *thrapb.ChainHeader, error) {
	// Create empty ref
	ldigest, _, err := thrap.store.CreateRef(stack.ID, "latest")
	if err != nil {
		return nil, nil, err
	}

	sdigest, err := thrap.store.Set(stack.ID, stack)
	if err != nil {
		return nil, nil, err
	}

	ref := &thrapb.ChainHeader{
		Previous:   ldigest,
		DataDigest: sdigest,
		Timestamp:  time.Now().UnixNano(),
	}

	_, err = thrap.store.SetRef(stack.ID, "latest", ref)
	return stack, ref, err
}
