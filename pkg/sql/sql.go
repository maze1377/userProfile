package sql

import "gorm.io/gorm"

type DbConfig interface {
	Connection() gorm.Dialector
}

type Migrate interface {
	Migrate() error
}
