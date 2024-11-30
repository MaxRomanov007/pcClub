package gorm

import "gorm.io/gorm"

func IsFailResult(res *gorm.DB) bool {
	return res.Error != nil || res.RowsAffected == 0
}
