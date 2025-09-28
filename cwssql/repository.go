package cwssql

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"gorm.io/gorm/clause"
)

// Repository is a generic interface for database operations
type IRepository[T any] interface {
	GetSession() *DBSession                             // Get the gorm db session
	GetContext() context.Context                        // Get the context
	GetAll(whereClause ...WhereCaluse) ([]*T, error)    // Get entities. Note: GetAll() equals Select all records.
	Get(whereClause ...WhereCaluse) (*T, error)         // Get entity
	Upsert(entity *T) error                             // Create or Replace
	Delete(entity *T) error                             // Delete entity
	DeleteAll(whereClause ...WhereCaluse) ([]*T, error) // Delete entities
	Refresh(entity *T) error                            // Refresh entity
	Count(whereClause ...WhereCaluse) (int64, error)    // Count entities.
	Begin() error                                       // Start a transaction
	Rollback() error                                    // Rollback a transiction
	Commit() error                                      // Commit a transiction
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

func (r *Repository[T]) GetSession() *DBSession {
	return r.session
}

func (r *Repository[T]) GetContext() context.Context {
	return r.context
}

func (r *Repository[T]) Get(whereClauses ...WhereCaluse) (*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entity T
	query := r.session.GetDb().WithContext(r.context)
	for _, wc := range whereClauses {
		for key, values := range wc {
			query = query.Where(key, values...)
		}
	}
	result := query.First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, result.Error
}

func (r *Repository[T]) GetAll(whereClauses ...WhereCaluse) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	query := r.session.GetDb().WithContext(r.context)
	for _, wc := range whereClauses {
		for key, values := range wc {
			query = query.Where(key, values...)
		}
	}
	result := query.Find(&entities)
	return entities, result.Error
}

func (r *Repository[T]) Upsert(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	assignments, err := GetNonPrimaryKeyAssignments(r.session.GetDb(), entity)
	if err != nil {
		return err
	}
	columns, err := GetPrimaryKeyColumns(r.session.GetDb(), entity)
	if err != nil {
		return err
	}
	return r.session.GetDb().Clauses(clause.OnConflict{
		Columns:   columns, // The conflicting primary key column(s)
		DoUpdates: assignments,
	}).Create(&entity).Error
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

func (r *Repository[T]) DeleteAll(whereClauses ...WhereCaluse) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	query := r.session.GetDb().Clauses(clause.Returning{}).WithContext(r.context)
	for _, wc := range whereClauses {
		for key, values := range wc {
			query = query.Where(key, values...)
		}
	}
	result := query.Delete(&entities)
	return entities, result.Error
}

func (r *Repository[T]) Count(whereClauses ...WhereCaluse) (int64, error) {
	if r.isGenericPointer() {
		return 0, errors.New("generic type T must be a struct")
	}
	var count int64
	query := r.session.GetDb().WithContext(r.context)
	for _, wc := range whereClauses {
		for key, values := range wc {
			query = query.Where(key, values...)
		}
	}
	result := query.Count(&count)
	return count, result.Error
}

func (r *Repository[T]) Refresh(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	wc, err := GetPrimaryKeyValueMap(r.session.GetDb(), entity)
	if err != nil {
		return err
	}
	result, err := r.Get(wc)
	if err != nil {
		return err
	}
	if result != nil {
		return applyValue(result, entity)
	}
	return nil
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

func applyValue(from any, to any) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}
	return nil
}

// NewRepository creates a new generic repository
func NewRepository[T any](ctx context.Context, session *DBSession) Repository[T] {
	// Create a new repository instance
	repo := Repository[T]{
		session: session,
		context: ctx,
	}
	repo.IRepository = &repo
	return repo
}
