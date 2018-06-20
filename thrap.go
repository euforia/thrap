package thrap

import (
	"context"

	"github.com/euforia/thrap/thrapb"
)

type GRPCService struct {
	core *Core
}

func NewService(core *Core) *GRPCService {
	return &GRPCService{core: core}
}

func (s *GRPCService) RegisterStack(ctx context.Context, st *thrapb.Stack) (*thrapb.Stack, error) {
	stack, _, err := s.core.RegisterStack(st)
	return stack, err
}

func (s *GRPCService) RegisterIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	nident, _, err := s.core.RegisterIdentity(ident)
	return nident, err
}

func (s *GRPCService) ConfirmIdentity(ctx context.Context, ident *thrapb.Identity) (*thrapb.Identity, error) {
	nident, err := s.core.ConfirmIdentity(ident)
	return nident, err
}
