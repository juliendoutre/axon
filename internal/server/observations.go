package server

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/juliendoutre/axon/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Server) CountObservations(
	ctx context.Context,
	input *v1.CountObservationsInput,
) (*v1.CountObservationsOutput, error) {
	row := s.pg.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM axon.observations WHERE timestamp >= $1 AND timestamp <= $2;",
		input.GetFrom().AsTime().Format(time.RFC3339),
		input.GetTo().AsTime().Format(time.RFC3339),
	)

	var count uint64
	if err := row.Scan(&count); err != nil {
		return nil, status.Errorf(codes.Internal, "counting observations")
	}

	return &v1.CountObservationsOutput{
		Count: count,
	}, nil
}

func (s *Server) ListObservations(
	ctx context.Context,
	input *v1.ListObservationsInput,
) (*v1.ListObservationsOutput, error) {
	rows, err := s.pg.Query(
		ctx,
		`SELECT id, timestamp, asset_type, asset_id, attributes, observer_claims
FROM axon.observations
WHERE timestamp >= $1 AND timestamp <= $2 ORDER BY timestamp DESC OFFSET $3 LIMIT $4;`,
		input.GetFrom().AsTime().Format(time.RFC3339),
		input.GetTo().AsTime().Format(time.RFC3339),
		input.GetPage()*input.GetPageSize(),
		input.GetPageSize(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "listing observations")
	}

	defer rows.Close()

	observations := []*v1.Observation{}

	for rows.Next() {
		var id, assetType, assetID string
		var timestamp time.Time
		var attributes, claims []byte

		if err := rows.Scan(
			&id,
			&timestamp,
			&assetType,
			&assetID,
			&attributes,
			&claims,
		); err != nil {
			return nil, status.Errorf(codes.Internal, "scanning observation")
		}

		observation := &v1.Observation{
			Id:         id,
			Timestamp:  timestamppb.New(timestamp),
			AssetType:  assetType,
			AssetId:    assetID,
			Attributes: &structpb.Struct{},
			Claims:     &structpb.Struct{},
		}

		if err := observation.Attributes.UnmarshalJSON(attributes); err != nil { //nolint:protogetter
			return nil, status.Errorf(codes.Internal, "scanning observation")
		}

		if err := observation.Claims.UnmarshalJSON(claims); err != nil { //nolint:protogetter
			return nil, status.Errorf(codes.Internal, "scanning observation")
		}

		observations = append(observations, observation)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "listing observations")
	}

	output := &v1.ListObservationsOutput{
		Observations: observations,
	}

	if uint32(len(observations)) == input.GetPageSize() {
		output.NextPage = input.GetPage() + 1
	}

	return output, nil
}
