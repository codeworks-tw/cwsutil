package cwssql

import (
	"encoding/json"
	"time"

	"github.com/codeworks-tw/cwsutil/cwsbase"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db_instance *gorm.DB = nil

func NewPostgresDB(db_connection_string string) (*gorm.DB, error) {
	var err error = nil
	if db_instance == nil {
		db_instance, err = gorm.Open(postgres.Open(db_connection_string), &gorm.Config{})
		if err != nil {
			db_instance = nil
			return nil, err
		}
		db, err := db_instance.DB()
		if err != nil {
			db_instance = nil
			return nil, err
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	}
	return db_instance, err
}

func NewSQLiteDB(file_path string) (*gorm.DB, error) {
	var err error = nil
	if db_instance == nil {
		db_instance, err = gorm.Open(sqlite.Open(file_path), &gorm.Config{})
		if err != nil {
			db_instance = nil
			return nil, err
		}
		db, err := db_instance.DB()
		if err != nil {
			db_instance = nil
			return nil, err
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	}
	return db_instance, err
}

type DBSession struct {
	db      *gorm.DB
	dbStack cwsbase.Stack[gorm.DB]
}

// popDbStack restores the previous DB reference after a transaction ends.
func (s *DBSession) popDbStack() {
	if db := s.dbStack.Pop(); db != nil {
		s.db = db
	}
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
	err := s.db.Rollback().Error
	s.popDbStack()
	return err
}

func (s *DBSession) Commit() error {
	result := s.db.Commit()
	if result.Error != nil {
		// Ensure the transaction is properly rolled back and connection returned.
		_ = s.db.Rollback()
	}
	s.popDbStack()
	return result.Error
}

func NewSession(db *gorm.DB) *DBSession {
	return &DBSession{db: db}
}

func GetPrimaryKeyValueMap(db *gorm.DB, model any) (WhereCaluse, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	var inInterface map[string]any
	inrec, _ := json.Marshal(model)
	json.Unmarshal(inrec, &inInterface)
	var wc WhereCaluse
	for _, field := range stmt.Schema.Fields {
		if field.TagSettings["PRIMARYKEY"] == "PRIMARYKEY" {
			value := inInterface[field.DBName]
			if value == nil {
				value = inInterface[field.Name] // field name might affect by json field tag
			}
			if wc == nil {
				wc = Eq(field.Name, value)
				continue
			}
			wc = wc.Eq(field.Name, value)
		}
	}
	return wc, nil
}

func GetPrimaryKeyColumns(db *gorm.DB, model any) ([]clause.Column, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}
	var columns []clause.Column
	for _, field := range stmt.Schema.Fields {
		if field.TagSettings["PRIMARYKEY"] == "PRIMARYKEY" {
			columns = append(columns, clause.Column{Name: field.DBName})
		}
	}
	return columns, nil
}

func GetNonPrimaryKeyAssignments(db *gorm.DB, model any, excludeColumns ...string) (clause.Set, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	var inInterface map[string]any
	inrec, _ := json.Marshal(model)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	var assignments []clause.Assignment
	for _, field := range stmt.Schema.Fields {
		if field.TagSettings["PRIMARYKEY"] != "PRIMARYKEY" && field.DBName != "" && field.DBName != "created_at" {
			if field.DBName == "updated_at" {
				inInterface[field.DBName] = time.Now()
			}
			exclue := false
			for _, excludeColumn := range excludeColumns {
				if field.DBName == excludeColumn || field.Name == excludeColumn {
					exclue = true
					break
				}
			}
			if exclue {
				continue
			}
			value := inInterface[field.DBName]
			if value == nil {
				value = inInterface[field.Name] // field name might affect by json field tag
			}
			assignments = append(assignments, clause.Assignment{
				Column: clause.Column{Name: field.DBName},
				Value:  value,
			})
		}
	}
	return assignments, nil
}

func GetPrimaryKeyAssignments(db *gorm.DB, model any) (clause.Set, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, err
	}

	var inInterface map[string]any
	inrec, _ := json.Marshal(model)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	var assignments []clause.Assignment
	for _, field := range stmt.Schema.Fields {
		if field.TagSettings["PRIMARYKEY"] == "PRIMARYKEY" {
			value := inInterface[field.DBName]
			if value == nil {
				value = inInterface[field.Name] // field name might affect by json field tag
			}
			assignments = append(assignments, clause.Assignment{
				Column: clause.Column{Name: field.DBName},
				Value:  value,
			})
		}
	}
	return assignments, nil
}
