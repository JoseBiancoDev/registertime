package handlers

import (
	"net/http"
	"time"

	"github.com/bianquiviri/control-horario/models"
	"github.com/bianquiviri/control-horario/utils"
	"github.com/gin-gonic/gin"
)

func StartLog(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Check if there's an active log
	var activeLog models.TimeLog
	if err := utils.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&activeLog).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is already an active session"})
		return
	}

	log := models.TimeLog{
		UserID:    userID,
		StartTime: time.Now(),
	}

	if err := utils.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start session"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

func StopLog(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var log models.TimeLog
	if err := utils.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	now := time.Now()
	log.EndTime = &now
	duration := now.Sub(log.StartTime).Hours()
	log.Duration = duration

	if err := utils.DB.Save(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not stop session"})
		return
	}

	// Send notification
	var user models.User
	utils.DB.First(&user, userID)
	go utils.SendNotificationEmail(user.Email, log)

	c.JSON(http.StatusOK, log)
}

func GetLogs(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var logs []models.TimeLog
	utils.DB.Where("user_id = ?", userID).Order("start_time desc").Find(&logs)

	c.JSON(http.StatusOK, logs)
}
