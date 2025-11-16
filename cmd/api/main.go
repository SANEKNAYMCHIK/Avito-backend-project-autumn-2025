package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/db"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/handlers"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/repositories"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer database.Close()

	repo := repositories.NewRepository(database.DB)
	reviewService := services.NewReviewService(repo)
	handler := handlers.NewHandler(reviewService)

	r := gin.Default()

	r.Use(handlers.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := database.DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "healthy"})
	})

	r.POST("/team/add", handler.CreateTeam)
	r.GET("/team/get", handler.GetTeam)
	r.POST("/users/setIsActive", handler.SetUserActive)
	r.POST("/pullRequest/create", handler.CreatePR)
	r.POST("/pullRequest/merge", handler.MergePR)
	r.POST("/pullRequest/reassign", handler.ReassignReviewer)
	r.GET("/users/getReview", handler.GetUserReviews)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
