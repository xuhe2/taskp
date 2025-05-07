package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() *Database {
	return &Database{
		DB: nil,
	}
}

func (d *Database) InitFromDSN(dsn string) error {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.DB = db

	d.migrate()

	return nil
}

func (d *Database) migrate() {
	d.AutoMigrate(&TaskRecord{})
}
