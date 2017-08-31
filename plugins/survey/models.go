package survey

import (
	"time"
)

// Form form
type Form struct {
	ID uint `gorm:"primary_key" json:"id"`

	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`

	Fields  []Field  `json:"fields"`
	Records []Record `json:"records"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Form) TableName() string {
	return "survey_forms"
}

// Field field
type Field struct {
	ID uint `gorm:"primary_key" json:"id"`

	Name      string `json:"name"`
	Label     string `json:"label"`
	Value     string `json:"value"`
	SortOrder int    `json:"sortOrder"`

	FormID uint `json:"formId"`
	Form   Form

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Field) TableName() string {
	return "survey_fields"
}

// Record record
type Record struct {
	ID uint `gorm:"primary_key" json:"id"`

	Value string `json:"value"`

	FormID uint `json:"formId"`
	Form   Form

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Record) TableName() string {
	return "survey_records"
}
