package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"`
	Name          string         `json:"name"`
	Role          string         `gorm:"default:'user';not null" json:"role"`
	ResetToken    string         `json:"-"`
	ResetTokenExp *time.Time     `json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	TimeLogs      []TimeLog      `json:"time_logs,omitempty"`
}

type TimeLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	StartTime     time.Time      `gorm:"not null" json:"start_time"`
	EndTime       *time.Time     `json:"end_time"`
	Duration      float64        `json:"duration_hours"` // In hours
	WorkMode      string         `json:"work_mode"`      // remote or on-site
	Comment       string         `json:"comment"`
	ClosedByAdmin bool           `gorm:"default:false" json:"closed_by_admin"`
	AdminID       *uint          `json:"admin_id"` // ID of the admin who closed/modified the log
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
