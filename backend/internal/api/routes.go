package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers your API endpoints.
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check DB connection
		sqlDB, err := db.DB()
		dbStatus := "up"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"db_status": dbStatus,
		})
	})

	// Add more endpoints as needed.
	// For example:
	router.GET("/api/guild-status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Guild status endpoint",
			"data":    "Sample data here",
		})
	})
}
