package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bianquiviri/control-horario/models"
	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migration
	err = DB.AutoMigrate(&models.User{}, &models.TimeLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connected and migrated")
	
	SeedDB(false) // Auto-seed if empty
}

func SeedDB(force bool) {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 && !force {
		return // DB already seeded
	}

	log.Println("Starting database seeding...")
	gofakeit.Seed(0)

	// 1. Create Admin
	hashedAdminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Email:    "admin@admin.com",
		Password: string(hashedAdminPassword),
		Name:     "Administrador",
		Role:     "admin",
	}
	// Use Upsert or check to avoid duplicates if forced
	DB.Where(models.User{Email: admin.Email}).FirstOrCreate(&admin)
	log.Println("Admin user ensured: admin@admin.com / admin123")

	// 2. Create 100 Users
	hashedUserPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	
	for i := 0; i < 100; i++ {
		user := models.User{
			Email:    gofakeit.Email(),
			Password: string(hashedUserPassword),
			Name:     gofakeit.Name(),
			Role:     "user",
		}
		
		if err := DB.Create(&user).Error; err != nil {
			continue
		}

		numLogs := gofakeit.Number(5, 20)
		for j := 0; j < numLogs; j++ {
			daysAgo := gofakeit.Number(1, 30)
			startHour := gofakeit.Number(8, 10)
			startMinute := gofakeit.Number(0, 59)
			
			startTime := time.Now().AddDate(0, 0, -daysAgo)
			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startHour, startMinute, 0, 0, startTime.Location())
			
			durationHours := float64(gofakeit.Number(4, 9)) + gofakeit.Float64Range(0.0, 0.99)
			endTime := startTime.Add(time.Duration(durationHours * float64(time.Hour)))
			
			workMode := "Presencial"
			if gofakeit.Bool() {
				workMode = "Remoto"
			}

			timeLog := models.TimeLog{
				UserID:    user.ID,
				StartTime: startTime,
				EndTime:   &endTime,
				Duration:  durationHours,
				WorkMode:  workMode,
				Comment:   gofakeit.Sentence(5),
			}
			
			DB.Create(&timeLog)
		}
	}
	log.Println("Database seeding completed.")
}
