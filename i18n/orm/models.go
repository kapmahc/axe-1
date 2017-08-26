package orm

import "time"

//Model locale model
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Lang      string    `gorm:"not null;type:varchar(8);index;unique_index:idx_locales_lang_code" json:"lang"`
	Code      string    `gorm:"not null;index;type:VARCHAR(255);unique_index:idx_locales_lang_code" json:"code"`
	Message   string    `gorm:"not null;type:varchar(1024)" json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (Model) TableName() string {
	return "locales"
}
