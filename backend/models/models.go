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

type Activity struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	IndicadoPor  string         `json:"indicado_por"`
	Resumen      string         `gorm:"type:text" json:"resumen"`
	Estado       string         `gorm:"default:'pendiente'" json:"estado"` // pendiente, iniciado, finalizado
	AsignadoToID uint           `json:"asignado_to_id"`
	AsignadoTo   User           `gorm:"foreignKey:AsignadoToID" json:"asignado_to"`
	CreadoByID   uint           `gorm:"not null" json:"creado_by_id"`
	CreadoBy     User           `gorm:"foreignKey:CreadoByID" json:"creado_by"`
	Files        []ActivityFile `json:"files,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ActivityFile struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ActivityID   uint           `gorm:"not null" json:"activity_id"`
	FilePath     string         `gorm:"not null" json:"file_path"`
	FileType     string         `gorm:"not null" json:"file_type"` // e.g., image/png, application/pdf
	UploadedByID uint           `json:"uploaded_by_id"`
	UploadedBy   User           `gorm:"foreignKey:UploadedByID" json:"uploaded_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Scopes
func FilterByUser(userID uint, role string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if role == "admin" {
			return db // Admins can see all activities
		}
		// Regular users can see activities assigned to them or created by them
		return db.Where("asignado_to_id = ? OR creado_by_id = ?", userID, userID)
	}
}
