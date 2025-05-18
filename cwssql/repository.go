package cwssql

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository is a generic interface for database operations
type IRepository[T any] interface {
	GetAll(equals map[string]any) ([]*T, error)    // Get entities
	Get(equals map[string]any) (*T, error)         // Get entity
	Upsert(entity *T) error                        // Create or Replace
	Delete(entity *T) error                        // Delete entity
	DeleteAll(equals map[string]any) ([]*T, error) // Delete entities
	Refresh(entity *T) error                       // Refresh entity
	Count(equals map[string]any) (int64, error)    // Count entities
	Begin() error                                  // Start a transaction
	Rollback() error                               // Rollback a transiction
	Commit() error                                 // Commit a transiction
}

type Repository[T any] struct {
	IRepository[T]
	session *DBSession
	context context.Context
}

func (r *Repository[T]) isGenericPointer() bool {
	var t T
	return reflect.ValueOf(t).Kind() == reflect.Pointer
}

func (r *Repository[T]) GetDB() *gorm.DB {
	return r.session.GetDb()
}

func (r *Repository[T]) Get(equals map[string]any) (*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entity T
	result := r.session.GetDb().WithContext(r.context).Where(mapToStruct[T](equals)).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entity, result.Error
}

func (r *Repository[T]) GetAll(equals map[string]any) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	result := r.session.GetDb().WithContext(r.context).Where(mapToStruct[T](equals)).Find(&entities)
	return entities, result.Error
}

func (r *Repository[T]) Upsert(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	result := r.session.GetDb().WithContext(r.context).Save(entity)
	return result.Error
}

func (r *Repository[T]) Delete(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	result := r.session.GetDb().WithContext(r.context).Delete(entity)
	return result.Error
}

func (r *Repository[T]) DeleteAll(equals map[string]any) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	result := r.session.GetDb().Clauses(clause.Returning{}).WithContext(r.context).Where(mapToStruct[T](equals)).Delete(&entities)
	return entities, result.Error
}

func (r *Repository[T]) Count(equals map[string]any) (int64, error) {
	if r.isGenericPointer() {
		return 0, errors.New("generic type T must be a struct")
	}
	var count int64
	result := r.session.GetDb().WithContext(r.context).Where(mapToStruct[T](equals)).Count(&count)
	return count, result.Error
}

func (r *Repository[T]) Refresh(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	primaryKeys, err := GetPrimaryKeyValueMap(r.session.GetDb(), entity)
	if err != nil {
		return err
	}
	result := r.session.GetDb().WithContext(r.context).Where(mapToStruct[T](primaryKeys)).First(entity)
	return result.Error
}

func (r *Repository[T]) Begin() error {
	return r.session.Begin()
}

func (r *Repository[T]) Rollback() error {
	return r.session.Rollback()
}

func (r *Repository[T]) Commit() error {
	return r.session.Commit()
}

// NewRepository creates a new generic repository
func NewRepository[T any](ctx context.Context, session *DBSession) Repository[T] {
	// Create a new repository instance
	var repo any = Repository[T]{
		session: session,
		context: ctx,
	}
	return repo.(Repository[T])
}

func mapToStruct[T any](data map[string]any) T {
	var entity T
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &entity)
	return entity
}
