package mssql

import (
	"fmt"
	"net/url"
	"server/internal/config"
	"strconv"
)

func GenerateConnString(cfg *config.SQLServerConfig) string {
	query := url.Values{}
	query.Add("app name", cfg.AppName)
	query.Add("database", cfg.Database)
	query.Add("encrypt", strconv.FormatBool(cfg.Encrypt))
	query.Add("TrustServerCertificate", strconv.FormatBool(cfg.TrustServerCertificate))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port),
		RawQuery: query.Encode(),
	}

	return u.String()
}
