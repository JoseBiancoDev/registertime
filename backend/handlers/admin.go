package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bianquiviri/control-horario/models"
	"github.com/bianquiviri/control-horario/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	// Only preload active logs to check session status
	utils.DB.Preload("TimeLogs", "end_time IS NULL").Find(&users)
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("default123"), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := utils.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetUserLogs(c *gin.Context) {
	userID := c.Param("id")
	var logs []models.TimeLog
	utils.DB.Where("user_id = ?", userID).Order("start_time desc").Find(&logs)
	c.JSON(http.StatusOK, logs)
}

func AdminStopLog(c *gin.Context) {
	adminID := c.MustGet("userID").(uint)
	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)

	var log models.TimeLog
	if err := utils.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found for this user"})
		return
	}

	now := time.Now()
	log.EndTime = &now
	log.Duration = now.Sub(log.StartTime).Hours()
	log.ClosedByAdmin = true
	log.AdminID = &adminID

	if err := utils.DB.Save(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not stop session"})
		return
	}

	c.JSON(http.StatusOK, log)
}

func AdminAddLog(c *gin.Context) {
	adminID := c.MustGet("userID").(uint)
	var input struct {
		UserID    uint      `json:"user_id" binding:"required"`
		StartTime time.Time `json:"start_time" binding:"required"`
		EndTime   time.Time `json:"end_time" binding:"required"`
		Comment   string    `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	duration := input.EndTime.Sub(input.StartTime).Hours()
	log := models.TimeLog{
		UserID:        input.UserID,
		StartTime:     input.StartTime,
		EndTime:       &input.EndTime,
		Duration:      duration,
		Comment:       input.Comment,
		ClosedByAdmin: true,
		AdminID:       &adminID,
	}

	if err := utils.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create log"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

type MonthlySummary struct {
	Month string  `json:"month"`
	Hours float64 `json:"hours"`
}

func GetUserMonthlyReport(c *gin.Context) {
	userID := c.Param("id")
	var summaries []MonthlySummary

	// Simple query to group by month (Postgres syntax)
	utils.DB.Raw(`
		SELECT TO_CHAR(start_time, 'YYYY-MM') as month, SUM(duration_hours) as hours
		FROM time_logs
		WHERE user_id = ? AND end_time IS NOT NULL
		GROUP BY month
		ORDER BY month DESC
	`, userID).Scan(&summaries)

	c.JSON(http.StatusOK, summaries)
}
