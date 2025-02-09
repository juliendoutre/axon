package server

import (
	"context"

	v1 "github.com/juliendoutre/axon/pkg/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetVersion(_ context.Context, _ *emptypb.Empty) (*v1.Version, error) {
	return s.version, nil
}
