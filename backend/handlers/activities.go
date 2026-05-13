package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bianquiviri/control-horario/models"
	"github.com/bianquiviri/control-horario/utils"
	"github.com/gin-gonic/gin"
)

func GetActivities(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var activities []models.Activity
	
	err := utils.DB.Scopes(models.FilterByUser(userID.(uint), role.(string))).
		Preload("Files").
		Preload("AsignadoTo").
		Preload("CreadoBy").
		Order("created_at desc").
		Find(&activities).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func CreateActivity(c *gin.Context) {
	userID, _ := c.Get("userID")

	name := c.PostForm("name")
	description := c.PostForm("description")
	asignadoToIDStr := c.PostForm("asignado_to_id")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	var asignadoToID uint
	if asignadoToIDStr != "" {
		id, err := strconv.ParseUint(asignadoToIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asignado_to_id"})
			return
		}
		asignadoToID = uint(id)
	} else {
		asignadoToID = userID.(uint)
	}

	activity := models.Activity{
		Name:         name,
		Description:  description,
		IndicadoPor:  "", // Ya no lo guardamos como string, usamos CreadoBy
		Estado:       "pendiente",
		AsignadoToID: asignadoToID,
		CreadoByID:   userID.(uint),
	}

	tx := utils.DB.Begin()

	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating activity"})
		return
	}

	// Handle files
	form, _ := c.MultipartForm()
	if form != nil {
		files := form.File["files"]
		for _, file := range files {
			// Validate file type
			ext := strings.ToLower(filepath.Ext(file.Filename))
			var fileType string
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				fileType = "image"
			} else if ext == ".pdf" {
				fileType = "pdf"
			} else {
				continue // Skip unsupported files
			}

			// Ensure uploads directory exists
			uploadDir := "./uploads"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.Mkdir(uploadDir, 0755)
			}

			// Save file
			filename := fmt.Sprintf("%d_%d%s", activity.ID, time.Now().UnixNano(), ext)
			filepathStr := filepath.Join(uploadDir, filename)
			if err := c.SaveUploadedFile(file, filepathStr); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
				return
			}

			activityFile := models.ActivityFile{
				ActivityID: activity.ID,
				FilePath:   fmt.Sprintf("/uploads/%s", filename),
				FileType:   fileType,
			}
			if err := tx.Create(&activityFile).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file record"})
				return
			}
		}
	}

	tx.Commit()

	// Fetch created activity with preloads to return
	utils.DB.Preload("Files").Preload("AsignadoTo").Preload("CreadoBy").First(&activity, activity.ID)
	c.JSON(http.StatusCreated, activity)
}

func UpdateActivityStatus(c *gin.Context) {
	activityIDStr := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	estado := c.PostForm("estado")
	resumen := c.PostForm("resumen")

	var activity models.Activity
	if err := utils.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	// Only assigned user, creator, or admin can update
	if role != "admin" && activity.AsignadoToID != userID.(uint) && activity.CreadoByID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this activity"})
		return
	}

	tx := utils.DB.Begin()

	if estado != "" {
		activity.Estado = estado
	}
	if resumen != "" {
		activity.Resumen = resumen
	}

	if err := tx.Save(&activity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating activity"})
		return
	}

	// Handle additional files if provided
	form, _ := c.MultipartForm()
	if form != nil && form.File != nil {
		files := form.File["files"]
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			var fileType string
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				fileType = "image"
			} else if ext == ".pdf" {
				fileType = "pdf"
			} else {
				continue // Skip unsupported files
			}

			uploadDir := "./uploads"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.Mkdir(uploadDir, 0755)
			}

			filename := fmt.Sprintf("%d_%d%s", activity.ID, time.Now().UnixNano(), ext)
			filepathStr := filepath.Join(uploadDir, filename)
			if err := c.SaveUploadedFile(file, filepathStr); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
				return
			}

			activityFile := models.ActivityFile{
				ActivityID: activity.ID,
				FilePath:   fmt.Sprintf("/uploads/%s", filename),
				FileType:   fileType,
			}
			if err := tx.Create(&activityFile).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file record"})
				return
			}
		}
	}

	tx.Commit()

	// Load the updated files and relations to return them
	utils.DB.Preload("Files").Preload("AsignadoTo").Preload("CreadoBy").First(&activity, activity.ID)
	c.JSON(http.StatusOK, activity)
}
