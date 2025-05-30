package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CustomTime is a custom time type that formats time in a consistent way
type CustomTime time.Time

// MarshalJSON implements the json.Marshaler interface
func (t CustomTime) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(tm.Format("2006-01-02 15:04:05"))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tm, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = CustomTime(tm)
	return nil
}

// Value implements the driver.Valuer interface
func (t CustomTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan implements the sql.Scanner interface
func (t *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*t = CustomTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = CustomTime(v)
		return nil
	case string:
		tm, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*t = CustomTime(tm)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomTime", value)
	}
}

// String returns the time in the format "2006-01-02 15:04:05"
func (t CustomTime) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt CustomTime     `json:"created_at"`
	UpdatedAt CustomTime     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Pagination represents common pagination parameters
type Pagination struct {
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
	Total    int64 `json:"total"`
}

func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	return p.PageSize
}
