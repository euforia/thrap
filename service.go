package thrap

import (
	"context"
	"log"

	"github.com/euforia/thrap/thrapb"
	"google.golang.org/grpc/metadata"
)

type GRPCService struct {
	core *Core
	log  *log.Logger
}

func NewService(core *Core, logger *log.Logger) *GRPCService {
	s := &GRPCService{core: core, log: logger}
	if s.log == nil {
		s.log = s.core.log
	}
	return s
}

func (s *GRPCService) handleIncomingContext(ctx context.Context, call string) {
	// TODO: auth
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		s.log.Println(call, ":", md)
	}
}

func (s *GRPCService) RegisterStack(ctx context.Context, st *thrapb.Stack) (*thrapb.Stack, error) {
	s.handleIncomingContext(ctx, "stack."+st.ID+".register")

	stack, _, err := s.core.RegisterStack(st)
	return stack, err
}

func (s *GRPCService) RegisterIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".register")

	nident, _, err := s.core.RegisterIdentity(ident)
	return nident, err
}

func (s *GRPCService) ConfirmIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	s.handleIncomingContext(ctx, "identity."+ident.ID+".confirm")

	nident, err := s.core.ConfirmIdentity(ident)
	return nident, err
}
