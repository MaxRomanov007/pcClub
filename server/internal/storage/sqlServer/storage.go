package sqlServer

import (
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	sql "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"net/url"
	"reflect"
	"server/internal/config"
	"strconv"
	"strings"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrReferenceNotExists = errors.New("reference does not exist")
	ErrTooLong            = errors.New("field too long")
	ErrNullPointer        = errors.New("null pointer")
)

type Storage struct {
	cfg *config.SQLServerConfig
	db  *sqlx.DB
}

func New(cfg *config.SQLServerConfig) (*Storage, error) {
	const op = "storage.sqlServer.New"

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

func handleError(err error) error {
	if driverErr, ok := err.(sql.Error); ok {
		switch driverErr.Number {
		case 2627: // Нарушение уникального индекса
			return ErrAlreadyExists
		case 547: // Нарушение внешнего ключа
			return ErrReferenceNotExists
		case 8152: // Строка слишком длинная для столбца
			return ErrTooLong
		case 515: // Попытка вставить NULL в столбец, который не допускает NULL
			return ErrNullPointer
		default:
			return err
		}
	}
	return err
}

func replacePositionalParams(query string, args []interface{}) string {
	for i := range args {
		param := fmt.Sprintf("@p%d", i+1)
		query = strings.Replace(query, "?", param, 1)
	}
	return query
}

func addPrefixToSlice(strings []string, prefix string) []string {
	result := make([]string, len(strings))
	for i, str := range strings {
		result[i] = prefix + str
	}
	return result
}

func appendStringSlices(slices ...[]string) []string {
	var result []string
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// GetDBTags возвращает слайс структурных тегов `db` из структуры.
func GetDBTags(s interface{}) []string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	tags := make([]string, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag != "" && tag != "-" {
			tags = append(tags, tag)
		}
	}

	return tags
}
