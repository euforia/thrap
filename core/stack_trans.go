package core

import (
	"context"
	"sync"

	"github.com/euforia/thrap/thrapb"
	"google.golang.org/grpc"
)

// StackTransport implements a stack transport interface
type StackTransport interface {
	Get(addr, id string) (*thrapb.Stack, error)
	Iter(addr, prefix string, f func(*thrapb.Stack) error) error
	Register(addr string, stack *thrapb.Stack) (*thrapb.Stack, error)
}

type remoteStack struct {
	mu    sync.Mutex
	conns map[string]thrapb.ThrapClient
}

func (st *remoteStack) Get(addr, id string) (*thrapb.Stack, error) {
	conn, err := st.getClient(addr)
	if err != nil {
		return nil, err
	}
	return conn.GetStack(context.Background(), &thrapb.Stack{ID: id})
}

func (st *remoteStack) Iter(addr, prefix string, f func(*thrapb.Stack) error) error {
	conn, err := st.getClient(addr)
	if err != nil {
		return err
	}

	in := &thrapb.IterOptions{Prefix: prefix}
	stream, err := conn.IterStacks(context.Background(), in)
	if err != nil {
		return err
	}

	defer stream.CloseSend()

	for {
		stack, err := stream.Recv()
		if err != nil {
			return err
		}

		//fmt.Fprintf(tw, "%s\t%s\n", stack.ID, stack.Version)
		err = f(stack)
		if err != nil {
			return err
		}
	}
}

func (st *remoteStack) Register(remote string, stack *thrapb.Stack) (*thrapb.Stack, error) {
	conn, err := st.getClient(remote)
	if err != nil {
		return nil, err
	}
	return conn.RegisterStack(context.Background(), stack)
}

func (st *remoteStack) getClient(addr string) (thrapb.ThrapClient, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	conn, ok := st.conns[addr]
	if ok {
		return conn, nil
	}

	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	conn = thrapb.NewThrapClient(cc)
	st.conns[addr] = conn

	return conn, nil
}

// local and remote
type stackTrans struct {
	local  StackStore
	remote *remoteStack
}

func (st *stackTrans) Get(remote, id string) (*thrapb.Stack, error) {
	if remote == "" {
		return st.local.Get(id)
	}
	return st.remote.Get(remote, id)
}

func (st *stackTrans) Iter(remote string, prefix string, f func(*thrapb.Stack) error) error {
	if remote == "" {
		return st.local.Iter(prefix, f)
	}

	return st.remote.Iter(remote, prefix, f)
}

func (st *stackTrans) Register(remote string, stack *thrapb.Stack) (*thrapb.Stack, error) {
	if remote == "" {
		stk, _, err := st.local.Register(stack)
		return stk, err
	}
	return st.remote.Register(remote, stack)
}
