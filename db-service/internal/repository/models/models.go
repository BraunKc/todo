package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}
