package cwssql

import (
	"encoding/json"
	"time"

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
			// When field type is json/jsonb, the value deserialized from JSON
			// becomes []any or map[string]any. GORM will expand these as SQL
			// records/tuples instead of a JSON string. Re-marshal back to a
			// JSON string so GORM treats it as a scalar value.
			switch value.(type) {
			case []any, map[string]any:
				if jsonBytes, err := json.Marshal(value); err == nil {
					value = string(jsonBytes)
				}
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
