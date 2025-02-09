package server

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/juliendoutre/axon/pkg/v1"
)

func New(version *v1.Version, pg *pgxpool.Pool) (*Server, error) {
	return &Server{
		version:   version,
		pg:        pg,
		jwtParser: jwt.NewParser(),
	}, nil
}

type Server struct {
	v1.UnimplementedAxonServer

	version   *v1.Version
	pg        *pgxpool.Pool
	jwtParser *jwt.Parser
}
