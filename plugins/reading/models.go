package reading

import (
	"time"

	"github.com/kapmahc/axe/plugins/nut"
)

// Book book
type Book struct {
	ID uint `gorm:"primary_key" json:"id"`

	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Lang        string    `json:"lang"`
	File        string    `json:"-"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Cover       string    `json:"cover"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Book) TableName() string {
	return "reading_books"
}

// Note note
type Note struct {
	ID uint `gorm:"primary_key" json:"id"`

	Type string
	Body string

	UserID uint
	User   nut.User
	BookID uint
	Book   Book

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Note) TableName() string {
	return "reading_notes"
}