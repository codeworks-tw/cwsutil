package cwssql

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IRepository is a generic interface for database operations using GORM
// T represents the entity type that this repository manages
type IRepository[T any] interface {
	GetGorm(whereClauses ...WhereCaluse) *gorm.DB       // Get the GORM DB object with applied where clauses
	GetSession() *DBSession                             // Get the CWS database session
	GetContext() context.Context                        // Get the current context
	GetAll(whereClause ...WhereCaluse) ([]*T, error)    // Get all entities matching the where clauses
	Get(whereClause ...WhereCaluse) (*T, error)         // Get a single entity matching the where clauses
	Upsert(entity *T, excludeColumns ...string) error   // Create or update entity (insert or update on conflict)
	Update(entity *T, excludeColumns ...string) error   // Update entity in database
	Delete(entity *T) error                             // Delete a specific entity
	DeleteAll(whereClause ...WhereCaluse) ([]*T, error) // Delete all entities matching the where clauses
	Refresh(entity *T) error                            // Refresh entity with latest data from database
	Count(whereClause ...WhereCaluse) (int64, error)    // Count entities matching the where clauses
	Begin() error                                       // Start a database transaction
	Rollback() error                                    // Rollback the current transaction
	Commit() error                                      // Commit the current transaction
}

// Repository is a concrete implementation of IRepository interface
// It provides GORM-based database operations for entities of type T
type Repository[T any] struct {
	IRepository[T]
	session *DBSession      // Database session for connection management
	context context.Context // Context for database operations
}

// isGenericPointer checks if the generic type T is a pointer type
// Returns true if T is a pointer, false if T is a struct
// This is used to validate that the repository operates on struct types, not pointers
func (r *Repository[T]) isGenericPointer() bool {
	var t T
	return reflect.ValueOf(t).Kind() == reflect.Pointer
}

// GetGorm returns a GORM DB instance with applied where clauses and context
// whereClauses: Optional where conditions to apply to the query
func (r *Repository[T]) GetGorm(whereClauses ...WhereCaluse) *gorm.DB {
	query := r.session.GetDb().WithContext(r.context)
	for _, wc := range whereClauses {
		for key, values := range wc {
			query = query.Where(key, values...)
		}
	}
	return query
}

// GetSession returns the database session associated with this repository
func (r *Repository[T]) GetSession() *DBSession {
	return r.session
}

// GetContext returns the context associated with this repository
func (r *Repository[T]) GetContext() context.Context {
	return r.context
}

// Get retrieves a single entity from the database matching the given where clauses
// Returns the first matching entity or an error if not found or if T is a pointer type
func (r *Repository[T]) Get(whereClauses ...WhereCaluse) (*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entity T
	result := r.GetGorm(whereClauses...).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, result.Error
}

// GetAll retrieves all entities from the database matching the given where clauses
// Returns a slice of pointers to entities or an error if T is a pointer type
func (r *Repository[T]) GetAll(whereClauses ...WhereCaluse) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	result := r.GetGorm(whereClauses...).Find(&entities)
	return entities, result.Error
}

// Upsert performs an insert or update operation (create or replace)
// If the entity exists (based on primary key), it updates the record
// If the entity doesn't exist, it creates a new record
// excludeColumns: Column names to exclude from the update operation
func (r *Repository[T]) Upsert(entity *T, excludeColumns ...string) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	assignments, err := GetNonPrimaryKeyAssignments(r.GetGorm(), entity, excludeColumns...)
	if err != nil {
		return err
	}
	columns, err := GetPrimaryKeyColumns(r.GetGorm(), entity)
	if err != nil {
		return err
	}
	return r.GetGorm().Clauses(clause.Returning{}, clause.OnConflict{
		Columns:   columns, // The conflicting primary key column(s)
		DoUpdates: assignments,
	}).Create(&entity).Error
}

func (r *Repository[T]) Update(entity *T, excludeColumns ...string) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	pAssignments, err := GetPrimaryKeyAssignments(r.GetGorm(), entity)
	if err != nil {
		return err
	}
	if len(pAssignments) == 0 {
		return errors.New("no primary key values found")
	}
	statement := r.GetGorm().Clauses(clause.Returning{})
	for _, assignment := range pAssignments {
		statement.Where(assignment.Column.Name+" = ?", assignment.Value)
	}
	result := statement.Updates(entity)
	if result.Error != nil {
		return result.Error
	}
	// If no rows were updated, it means the entity was not found
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete removes the specified entity from the database
// The entity must have primary key values set to identify which record to delete
func (r *Repository[T]) Delete(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	result := r.GetGorm().Delete(entity)
	return result.Error
}

// DeleteAll removes all entities from the database matching the given where clauses
// Returns the deleted entities and any error that occurred during deletion
func (r *Repository[T]) DeleteAll(whereClauses ...WhereCaluse) ([]*T, error) {
	if r.isGenericPointer() {
		return nil, errors.New("generic type T must be a struct")
	}
	var entities []*T
	result := r.GetGorm(whereClauses...).Delete(&entities)
	return entities, result.Error
}

// Count returns the number of entities matching the given where clauses
// Returns the count as int64 and any error that occurred during counting
func (r *Repository[T]) Count(whereClauses ...WhereCaluse) (int64, error) {
	if r.isGenericPointer() {
		return 0, errors.New("generic type T must be a struct")
	}
	var count int64
	result := r.GetGorm(whereClauses...).Count(&count)
	return count, result.Error
}

// Refresh reloads the entity with the latest data from the database
// Uses the entity's primary key to fetch the current state and updates the provided entity
func (r *Repository[T]) Refresh(entity *T) error {
	if entity == nil {
		return errors.New("entity must not be nil")
	}
	if r.isGenericPointer() {
		return errors.New("generic type T must be a struct")
	}
	wc, err := GetPrimaryKeyValueMap(r.GetGorm(), entity)
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

// Begin starts a new database transaction
func (r *Repository[T]) Begin() error {
	return r.session.Begin()
}

// Rollback rolls back the current transaction, undoing all changes made within the transaction
func (r *Repository[T]) Rollback() error {
	return r.session.Rollback()
}

// Commit commits the current transaction, making all changes permanent
func (r *Repository[T]) Commit() error {
	return r.session.Commit()
}

// applyValue copies values from one struct to another using JSON marshaling/unmarshaling
// This is used internally to update entity values after refresh operations
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

// NewRepository creates a new generic repository instance for entity type T
// ctx: Context to be used for database operations
// session: Database session for connection management
// Returns a configured Repository instance ready for database operations
func NewRepository[T any](ctx context.Context, session *DBSession) Repository[T] {
	// Create a new repository instance with provided session and context
	repo := Repository[T]{
		session: session,
		context: ctx,
	}
	// Set the interface reference to enable method calls
	repo.IRepository = &repo
	return repo
}
