package ssms

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"net/url"
	"reflect"
	"regexp"
	"server/internal/config"
	"strconv"
	"strings"
)

type Storage struct {
	cfg *config.SQLServerConfig
	db  *sqlx.DB
}

func New(cfg *config.SQLServerConfig) (*Storage, error) {
	const op = "storage.ssms.New"

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

	db, err := sqlx.Connect("sqlserver", u.String())
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %w", op, err)
	}

	return &Storage{
		db:  db,
		cfg: cfg,
	}, nil
}

func replacePositionalParams(query string, args []interface{}) string {
	for i := range args {
		param := fmt.Sprintf("@p%d", i+1)
		query = strings.Replace(query, "?", param, 1)
	}
	return query
}

func nullString(s string) sql.Null[string] {
	nullable := sql.Null[string]{V: s}
	if s != "" {
		nullable.Valid = true
	}
	return nullable
}

func setIfNotZeroNullString(stmt squirrel.UpdateBuilder, columnName string, str string) squirrel.UpdateBuilder {
	if match, err := regexp.MatchString("[Nn][Uu][Ll][Ll]", str); err == nil && match {
		stmt = stmt.Set(columnName, sql.Null[string]{})
	} else {
		stmt = setIfNotZero(stmt, columnName, str)
	}
	return stmt
}
func setIfNotZero(stmt squirrel.UpdateBuilder, columnName string, a any) squirrel.UpdateBuilder {
	if !isZero(a) {
		stmt = stmt.Set(columnName, a)
	}
	return stmt
}
func isZero(a any) bool {
	if reflect.ValueOf(a).Interface() == reflect.Zero(reflect.TypeOf(a)).Interface() {
		return true
	}
	return false
}
