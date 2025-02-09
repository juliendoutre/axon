package server

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/open-policy-agent/opa/v1/rego"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Authenticator struct {
	Parser *jwt.Parser
	Store  jwt.Keyfunc
	Policy string
}

func (a *Authenticator) Authenticate(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	authorization := md["authorization"]
	if len(authorization) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing bearer token")
	}

	jwtToken, err := a.Parser.Parse(token, a.Store)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid JWT token")
	}

	if !jwtToken.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid JWT token")
	}

	policyEvaluation, err := rego.New(
		rego.Query("data.authz.allowed"),
		rego.Module("policy.rego", a.Policy),
		rego.Input(jwtToken.Claims),
	).Eval(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "policy evaluation")
	}

	if !policyEvaluation.Allowed() {
		return nil, status.Errorf(codes.PermissionDenied, "not allowed")
	}

	return handler(ctx, req)
}
