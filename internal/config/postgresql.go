package config

import (
	"net"
	"net/url"
	"os"
)

func PostgresURL() *url.URL {
	pgQuery := url.Values{}
	pgQuery.Add("sslmode", "disable")

	return &url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		User:     url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Path:     os.Getenv("POSTGRES_DB"),
		RawQuery: pgQuery.Encode(),
	}
}

func MigrationsURL() *url.URL {
	return &url.URL{
		Scheme: "file",
		Path:   os.Getenv("MIGRATIONS_PATH"),
	}
}
