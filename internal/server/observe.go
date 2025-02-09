package server

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/juliendoutre/axon/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Observe(ctx context.Context, input *v1.ObserveInput) (*emptypb.Empty, error) {
	serializedAttributes, err := input.GetAttributes().MarshalJSON()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "serializing attributes")
	}

	md, _ := metadata.FromIncomingContext(ctx)
	token, _, _ := s.jwtParser.ParseUnverified(strings.TrimPrefix(md["authorization"][0], "Bearer "), jwt.MapClaims{})

	serializedClaims, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "serializing claims")
	}

	if _, err := s.pg.Exec(
		ctx,
		"INSERT INTO axon.observations (asset_type, asset_id, attributes, observer_claims) VALUES ($1, $2, $3, $4);",
		input.GetAssetType(),
		input.GetAssetId(),
		serializedAttributes,
		serializedClaims,
	); err != nil {
		return nil, status.Errorf(codes.Internal, "inserting observation")
	}

	return &emptypb.Empty{}, nil
}
