package sql

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteConfig struct {
	FileName string
	InMemory bool
}

func (c SqliteConfig) Connection() gorm.Dialector {
	if c.InMemory {
		return sqlite.Open("file::memory:?cache=shared")
	}

	return sqlite.Open(c.FileName)
}
