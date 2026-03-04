// models/user.go
package models
import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Role         string    `gorm:"type:varchar(20);default:'user'"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


//table Name overrides the default table name for User model
func (User) TableName() string {
	return "users_tbl"
}