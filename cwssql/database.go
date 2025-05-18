package cwssql

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db_instance *gorm.DB
var once sync.Once

func NewPostgresDB(db_connection_string string) (*gorm.DB, error) {
	var err error = nil
	once.Do(func() {
		db_instance, err = gorm.Open(postgres.Open(db_connection_string), &gorm.Config{})
		if err != nil {
			return
		}
		db, err := db_instance.DB()
		if err != nil {
			return
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	})
	return db_instance, err
}

func NewSQLiteDB(file_path string) (*gorm.DB, error) {
	var err error = nil
	once.Do(func() {
		db_instance, err = gorm.Open(sqlite.Open(file_path), &gorm.Config{})
		if err != nil {
			return
		}
		db, err := db_instance.DB()
		if err != nil {
			return
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	})
	return db_instance, err
}

type DBSession struct {
	db      *gorm.DB
	dbStack cwsbase.Stack[gorm.DB]
}

func (s *DBSession) GetDb() *gorm.DB {
	return s.db
}

func (s *DBSession) Begin() error {
	tx := s.db.Begin()
	if tx.Error == nil {
		s.dbStack.Push(s.db)
		s.db = tx
	}
	return tx.Error
}

func (s *DBSession) Rollback() error {
	return s.db.Rollback().Error
}

func (s *DBSession) Commit() error {
	defer func() {
		if s.dbStack.Len() > 0 {
			s.db = s.dbStack.Pop()
		}
	}()
	result := s.db.Commit()
	if result.Error != nil {
		err := s.Rollback()
		if err != nil {
			return err
		}
	}
	return result.Error
}

func NewSession(db *gorm.DB) *DBSession {
	return &DBSession{db: db}
}

func GetPrimaryKeyValueMap(db *gorm.DB, model any) (map[string]any, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	var inInterface map[string]any
	inrec, _ := json.Marshal(model)
	json.Unmarshal(inrec, &inInterface)
	result := map[string]any{}
	for _, field := range stmt.Schema.Fields {
		if field.TagSettings["PRIMARYKEY"] == "PRIMARYKEY" {
			result[field.Name] = inInterface[field.Name]
		}
	}
	return result, nil
}
