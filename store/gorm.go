package store

import (
	"database/sql"

	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

type Storer interface {
	New(interface{}) error
	Find(interface{}, interface{}, string) error
	ScanRows(*sql.Rows, interface{}) error
	First(interface{}, int, string) error
	Updates(map[string]interface{}, interface{}, interface{}) error
	Delete(interface{}, int) error
	Begin() *gorm.DB
	Table(string) (tx *gorm.DB)
	Save(interface{}) (tx *gorm.DB)
}

func (s *GormStore) New(table interface{}) error {
	return s.db.Create(table).Error
}

func (s *GormStore) Find(table interface{}, where interface{}, joinTable string) error {
	if joinTable != "" {
		return s.db.Preload(joinTable).Where(where).Find(table).Error
	} else {
		return s.db.Where(where).Find(table).Error
	}
}

func (s *GormStore) Table(table string) (tx *gorm.DB) {
	return s.db.Table(table)
}

func (s *GormStore) First(table interface{}, id int, joinTable string) error {
	if joinTable != "" {
		return s.db.Preload(joinTable).First(table, id).Error
	} else {
		return s.db.First(table, id).Error
	}
}

func (s *GormStore) Begin() *gorm.DB {
	return s.db.Begin()
}

func (s *GormStore) Updates(where map[string]interface{}, model interface{}, payload interface{}) error {
	return s.db.Model(model).Where(where).Updates(payload).Error
}

func (s *GormStore) Delete(table interface{}, id int) error {
	return s.db.Delete(table, id).Error
}

func (s *GormStore) ScanRows(rows *sql.Rows, table interface{}) error {
	return s.db.ScanRows(rows, table)
}

func (s *GormStore) Save(table interface{}) (tx *gorm.DB) {
	return s.db.Save(table)
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}
