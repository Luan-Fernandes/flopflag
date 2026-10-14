package db

import (
	"fmt"
	"net/url"

	"boliplate-go/internal/config"
)

func PostgresURL(cfg config.DBConfig) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	q := u.Query()
	if cfg.SSLMode != "" {
		q.Set("sslmode", cfg.SSLMode)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
