package handlers

import (
	"fmt"
	"net/http"

	"github.com/bianquiviri/control-horario/models"
	"github.com/bianquiviri/control-horario/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GenerateExcelReport(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var logs []models.TimeLog
	utils.DB.Where("user_id = ?", userID).Order("start_time desc").Find(&logs)

	f := excelize.NewFile()
	sheet := "Reporte de Horas"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	headers := []string{"ID", "Fecha Inicio", "Fecha Fin", "Duración (Horas)", "Comentario"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Data
	for i, log := range logs {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), log.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), log.StartTime.Format("2006-01-02 15:04:05"))
		if log.EndTime != nil {
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), log.EndTime.Format("2006-01-02 15:04:05"))
		}
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), log.Duration)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), log.Comment)
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=reporte_horas.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate Excel report"})
	}
}
