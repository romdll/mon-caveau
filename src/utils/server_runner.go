package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func RunWithQuitNotification(serverEngine *gin.Engine) {
	srv := &http.Server{
		Addr:    ":80",
		Handler: serverEngine.Handler(),
	}

	logger.Println("Starting server at :80")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error running the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Println("Server gracefully stopped")
}
