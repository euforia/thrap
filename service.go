package thrap

import (
	"context"
	"log"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/thrapb"
	"google.golang.org/grpc/metadata"
)

// GRPCService implements the server-side grpc service for thrap
type GRPCService struct {
	core *core.Core
	log  *log.Logger
}

// NewService returns a new grpc service with the given core
func NewService(core *core.Core, logger *log.Logger) *GRPCService {
	s := &GRPCService{core: core, log: logger}
	// if s.log == nil {
	// 	s.log = s.core.log
	// }
	return s
}

// RegisterIdentity implements the server-side grpc call
func (s *GRPCService) RegisterIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".register")

	idt := s.core.Identity()
	nident, _, err := idt.Register(ident)
	return nident, err
}

// ConfirmIdentity implements the server-side grpc call
func (s *GRPCService) ConfirmIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".confirm")

	idt := s.core.Identity()
	nident, err := idt.Confirm(ident)
	return nident, err
}

func (s *GRPCService) GetIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".get")

	idt := s.core.Identity()
	return idt.Get(ident.ID)
}

func (s *GRPCService) IterIdentities(opts *thrapb.IterOptions, stream thrapb.Thrap_IterIdentitiesServer) error {
	s.handleIncomingContext(stream.Context(), "identity.list")

	idt := s.core.Identity()
	return idt.Iter(opts.Prefix, func(ident *thrapb.Identity) error {
		return stream.Send(ident)
	})
}

// RegisterStack implements the server-side grpc call
func (s *GRPCService) RegisterStack(ctx context.Context, st *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+st.ID+".register")

	stk := s.core.Stack()
	stack, _, err := stk.Register(st)
	return stack, err
}

func (s *GRPCService) CommitStack(ctx context.Context, stack *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+stack.ID+".commit")

	stk := s.core.Stack()
	return stk.Commit(stack)
}

func (s *GRPCService) GetStack(ctx context.Context, stack *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+stack.ID+".get")

	st := s.core.Stack()
	return st.Get(stack.ID)
}

func (s *GRPCService) IterStacks(opts *thrapb.IterOptions, stream thrapb.Thrap_IterStacksServer) error {
	s.handleIncomingContext(stream.Context(), "identity.list")

	st := s.core.Stack()
	return st.Iter(opts.Prefix, func(stack *thrapb.Stack) error {
		return stream.Send(stack)
	})
}

func (s *GRPCService) handleIncomingContext(ctx context.Context, call string) {
	// TODO: auth
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		s.log.Println(call, ":", md)
	}
}
