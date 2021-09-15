package sql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabase(config DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open(config.Connection(), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
