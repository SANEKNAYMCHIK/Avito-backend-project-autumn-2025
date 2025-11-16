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
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer database.Close()

	repo := repositories.NewRepository(database.DB)
	_ = repo

	r := gin.Default()

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }
	port := "8080"

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
