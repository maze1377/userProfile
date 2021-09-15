package sql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSL      bool
}

func (c PostgresConfig) Connection() gorm.Dialector {
	sslMode := "enable"
	if !c.SSL {
		sslMode = "disable"
	}

	return postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Database,
		c.Password,
		sslMode,
	))
}
