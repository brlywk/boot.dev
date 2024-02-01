package config

import (
	"brlywk/bootdev/webserver/db"
	"time"
)

type ApiConfig struct {
	DbPath    string
	Db        *db.DB
	JwtSecret []byte
	PolkaKey  string
	TokenSettings
}

type TokenSettings struct {
	AccessIssuer     string
	AccessExpiresIn  time.Duration
	RefreshIssuer    string
	RefreshExpiresIn time.Duration
}
