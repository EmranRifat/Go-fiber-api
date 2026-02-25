// models/user.go
package models

import "time"

// type User struct {
// 	ID           uint      `json:"id"            gorm:"primaryKey"`
// 	Name         string    `json:"name"          gorm:"size:120;not null"`
// 	Email        string    `json:"email"         gorm:"size:255;not null;uniqueIndex:ux_users_email_lower"`
// 	PasswordHash string    `json:"-"             gorm:"size:255;not null"` // do not expose
// 	CreatedAt    time.Time `json:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at"`
// }
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


//table Name overrides the default table name for User model
func (User) TableName() string {
	return "users_tbl"
}