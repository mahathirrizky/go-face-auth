package models

import "gorm.io/gorm"

type AttendanceLocation struct {
	gorm.Model
	CompanyID uint   `gorm:"not null" json:"company_id"`
	Name      string `gorm:"not null" json:"name"`
	Latitude  float64 `gorm:"not null" json:"latitude"`
	Longitude float64 `gorm:"not null" json:"longitude"`
	Radius    uint   `gorm:"not null" json:"radius"`
}
