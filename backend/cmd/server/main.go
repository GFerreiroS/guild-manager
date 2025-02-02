package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	_ "github.com/lib/pq"
)

func createMyRenderer() render.HTMLRender {
	// Create template from string
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

	// Create HTML production renderer
	return render.HTMLProduction{Template: tmpl}
}

func setupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.HTMLRender = createMyRenderer()

	// Add database connection to context
	router.Use(func(c *gin.Context) {
		if db != nil {
			c.Set("db", db)
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Get database connection from context safely
		val, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"error":   "database connection missing from context",
				"version": "0.1.0",
			})
			return
		}

		// Type assertion with proper error handling
		db, ok := val.(*sql.DB)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "unhealthy",
				"error":   "invalid database connection type",
				"version": "0.1.0",
			})
			return
		}

		// Verify database connection
		ctx, cancel := context.WithTimeout(c, 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
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

func main() {
	// Database connection with retries
	db, err := connectDBWithRetries(5, 5*time.Second)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Initialize HTTP server
	router := setupRouter(db)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func connectDBWithRetries(maxRetries int, interval time.Duration) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Database initialization failed: %v", err)
			time.Sleep(interval)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = db.PingContext(ctx); err == nil {
			log.Println("Successfully connected to PostgreSQL")
			return db, nil
		}

		log.Printf("Connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
}
