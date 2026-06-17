package models

import "time"

type APILogs struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Method       string    `gorm:"type:varchar(10);index" json:"method"`
	Operation    string    `gorm:"type:varchar(20);index" json:"operation"`
	Path         string    `gorm:"type:text;index" json:"path"`
	Query        string    `gorm:"type:text" json:"query"`
	StatusCode   int       `gorm:"index" json:"statusCode"`
	Status       string    `gorm:"type:varchar(50);index" json:"status"`
	Message      string    `gorm:"type:text" json:"message"`
	IsError      bool      `gorm:"index" json:"isError"`
	ErrorType    string    `gorm:"type:varchar(50);index" json:"errorType"`
	IPAddress    string    `gorm:"type:varchar(100)" json:"ipAddress"`
	UserAgent    string    `gorm:"type:text" json:"userAgent"`
	UserID       string    `gorm:"type:varchar(100);index" json:"userId"`
	UserEmail    string    `gorm:"type:varchar(255);index" json:"userEmail"`
	LatencyMS    int64     `json:"latencyMs"`
	ErrorMessage string    `gorm:"type:text" json:"errorMessage"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

func (APILogs) TableName() string {
	return "api_logs"
}
