package project

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"size:256"`
	Blocked bool `gorm:"default:false"`
}