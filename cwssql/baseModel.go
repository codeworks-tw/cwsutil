package cwssql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type BaseIdModel struct {
	Id string `gorm:"type:text;primaryKey;default:uuid_generate_v4()"`
}

type BaseTimeModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// JSONB Interface for JSONB Field
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

type BaseJsonBModel struct {
	Attributes JSONB `gorm:"type:jsonb" json:"attributes"`
}
