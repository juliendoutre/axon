package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

func PostgresURL() (*url.URL, error) {
	pgQuery := url.Values{}
	pgQuery.Add("sslmode", "require")

	password, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_PATH"))
	if err != nil {
		return nil, fmt.Errorf("reading PG password file: %w", err)
	}

	return &url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		User:     url.UserPassword(os.Getenv("POSTGRES_USER"), string(password)),
		Path:     os.Getenv("POSTGRES_DB"),
		RawQuery: pgQuery.Encode(),
	}, nil
}

func MigrationsURL() url.URL {
	return url.URL{
		Scheme: "file",
		Path:   os.Getenv("MIGRATIONS_PATH"),
	}
}
