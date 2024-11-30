package mssql

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"server/internal/config"
	"server/internal/lib/api/database/mssql"
)

type Storage struct {
	cfg *config.SQLServerConfig
	db  *gorm.DB
}

const (
	AvailablePcStatus = "available"
)

func New(cfg *config.SQLServerConfig) (*Storage, error) {
	const op = "storage.mssql.New"

	db, err := gorm.Open(sqlserver.Open(mssql.GenerateConnString(cfg)), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %w", op, err)
	}

	return &Storage{
		db:  db,
		cfg: cfg,
	}, nil
}
