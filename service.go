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

func (s *GRPCService) handleIncomingContext(ctx context.Context, call string) {
	// TODO: auth
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		s.log.Println(call, ":", md)
	}
}

// RegisterStack implements the server-side grpc call
func (s *GRPCService) RegisterStack(ctx context.Context, st *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+st.ID+".register")

	stk := s.core.Stack()
	stack, _, err := stk.Register(st)
	return stack, err
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
