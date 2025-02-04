package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Adjust the import path as needed.
)

// createMyRenderer creates an HTML renderer for Gin.
func createMyRenderer() render.HTMLRender {
	tmpl := template.Must(template.New("status").Parse(`
	<div class="space-y-2">
		<p class="text-green-600">
			Online Players: {{ .OnlinePlayers }}
		</p>
		<p class="text-gray-500 text-sm">
			Last updated: {{ .LastUpdated }}
		</p>
	</div>
	`))
	return render.HTMLProduction{Template: tmpl}
}

// setupRouter configures the Gin router with routes and middleware.
func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.HTMLRender = createMyRenderer()

	// Store the underlying *sql.DB in Gin's context.
	router.Use(func(c *gin.Context) {
		if sqlDB, err := db.DB(); err == nil {
			c.Set("db", sqlDB)
		}
		c.Next()
	})

	// Health check endpoint.
	router.GET("/health", func(c *gin.Context) {
		val, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"error":   "database connection missing from context",
				"version": "0.1.0",
			})
			return
		}

		sqlDB, ok := val.(*sql.DB)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "unhealthy",
				"error":   "invalid database connection type",
				"version": "0.1.0",
			})
			return
		}

		// Check database health with a timeout.
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"error":   err.Error(),
				"version": "0.1.0",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	// Example API endpoint.
	router.GET("/api/guild-status", func(c *gin.Context) {
		c.Header("HX-Reswap", "innerHTML")
		c.Header("HX-Retarget", "#status-container")
		c.HTML(http.StatusOK, "status", gin.H{
			"OnlinePlayers": 42,
			"LastUpdated":   time.Now().Format("15:04:05"),
		})
	})

	return router
}

func connectDBWithRetries(dsn string, maxRetries int, delay time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}

		// Retrieve the underlying *sql.DB
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error retrieving underlying *sql.DB on attempt %d/%d: %v", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}

		// Ping the database to ensure the connection is ready.
		if pingErr := sqlDB.Ping(); pingErr != nil {
			log.Printf("Database ping attempt %d/%d failed: %v", i+1, maxRetries, pingErr)
			time.Sleep(delay)
			continue
		}

		log.Printf("Successfully connected to the database on attempt %d", i+1)
		return db, nil
	}

	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}

func main() {
	// Build the DSN string for PostgreSQL.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)

	// Get the maximum number of retries from the environment (default to 10 if not set).
	maxRetries := 10
	if s := os.Getenv("MIGRATION_MAX_RETRIES"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			maxRetries = v
		}
	}

	// Use a fixed delay between retries.
	delay := 5 * time.Second

	// Attempt to connect to the database with retries.
	db, err := connectDBWithRetries(dsn, maxRetries, delay)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize and start the HTTP server.
	router := setupRouter(db)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
