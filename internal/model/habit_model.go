package model

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Info      string    `gorm:"size:255;not null" json:"info"`
	UserID    uint64    `json:"-"`
}

func (Habit) TableName() string {
	return "habits"
}

func (h *Habit) BeforeCreate(tx *gorm.DB) error {
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()
	return nil
}

func (h *Habit) BeforeUpdate(tx *gorm.DB) error {
	h.UpdatedAt = time.Now()
	return nil
}
