package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username" binding:"required"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Status    int8           `gorm:"default:1;not null" json:"status"`
	Habits    []Habit        `gorm:"foreignkey:UserID" json:"habits"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if u.Status == 0 {
		u.Status = 1
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
