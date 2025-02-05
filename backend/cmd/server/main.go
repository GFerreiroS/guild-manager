package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	// Import your internal packages – adjust the import paths if necessary.
	"github.com/GFerreiroS/guild-manager/backend/internal/api"
	"github.com/GFerreiroS/guild-manager/backend/internal/config"
	"github.com/GFerreiroS/guild-manager/backend/internal/database"
	"github.com/GFerreiroS/guild-manager/backend/internal/middleware"
	"github.com/GFerreiroS/guild-manager/backend/pkg/redis"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

// setupRouter configures the Gin router, registers routes, and applies middleware.
func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.HTMLRender = createMyRenderer()

	// (Optional) You can add additional middleware here.
	// For example, a custom middleware to store the underlying *sql.DB for health checks:
	router.Use(func(c *gin.Context) {
		if sqlDB, err := db.DB(); err == nil {
			c.Set("db", sqlDB)
		}
		c.Next()
	})

	// Register your API endpoints.
	api.RegisterRoutes(router, db)

	return router
}

// connectDBWithRetries attempts to connect to PostgreSQL with retries.
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
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error retrieving underlying *sql.DB on attempt %d/%d: %v", i+1, maxRetries, err)
			time.Sleep(delay)
			continue
		}
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
	// Load configuration using Viper.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Build the DSN string from the configuration.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Get maximum retry count from an environment variable or use a default.
	maxRetries := 10
	if s := os.Getenv("MIGRATION_MAX_RETRIES"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			maxRetries = v
		}
	}
	delay := 5 * time.Second

	// Connect to the database with retry logic.
	db, err := connectDBWithRetries(dsn, maxRetries, delay)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// (Optional) Run migrations. If you are managing your schema via SQL migrations,
	// call the following function (ensure it’s exported from internal/database/migrations.go):
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Database migrations failed:", err)
	}

	// Initialize Redis client.
	redisClient := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.Timeout)

	// Create a new Gin router.
	router := setupRouter(db)

	// Apply rate-limiting middleware using Redis.
	router.Use(middleware.RateLimitMiddleware(redisClient.Conn, cfg.RateLimit.RequestsPerMinute))

	// Start the HTTP server.
	addr := ":8080"
	log.Printf("🚀 Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
