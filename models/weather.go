package models

import "time"
	
type Weather struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Division     string    `json:"division" gorm:"size:100;not null;uniqueIndex:ux_weather_division"`
	Lat          float64   `json:"lat" gorm:"not null"`
	Lon          float64   `json:"lon" gorm:"not null"`
	TemperatureC float64   `json:"temperature_c" gorm:"not null"`
	Humidity     int       `json:"humidity" gorm:"not null"`
	Condition    string    `json:"condition" gorm:"size:100;not null"`
	WindKph      float64   `json:"wind_kph" gorm:"not null"`
	VisibilityKm float64   `json:"visibility_km" gorm:"not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null"`
}

func (Weather) TableName() string { return "weather_data" }
