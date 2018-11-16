package thrap

import (
	"context"
	"log"

	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/pkg/pb"
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

// GetIdentity implements the server-side grpc call
func (s *GRPCService) GetIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".get")

	idt := s.core.Identity()
	return idt.Get(ident.ID)
}

// IterIdentities implements the server-side grpc call
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

	stk, err := s.core.Stack(pb.DefaultProfile())
	if err != nil {
		return nil, err
	}
	stack, _, err := stk.Register(st)
	return stack, err
}

// CommitStack implements the server-side grpc call
func (s *GRPCService) CommitStack(ctx context.Context, stack *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+stack.ID+".commit")

	stk, err := s.core.Stack(pb.DefaultProfile())
	if err == nil {
		return stk.Commit(stack)
	}
	return nil, err
}

// GetStack implements the server-side grpc call
func (s *GRPCService) GetStack(ctx context.Context, stack *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+stack.ID+".get")

	stk, err := s.core.Stack(pb.DefaultProfile())
	if err == nil {
		return stk.Get(stack.ID)
	}
	return nil, err
}

// IterStacks implements the server-side grpc call
func (s *GRPCService) IterStacks(opts *thrapb.IterOptions, stream thrapb.Thrap_IterStacksServer) error {
	s.handleIncomingContext(stream.Context(), "identity.list")

	stk, err := s.core.Stack(pb.DefaultProfile())
	if err != nil {
		return err
	}
	return stk.Iter(opts.Prefix, func(stack *thrapb.Stack) error {
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
