package utils

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

var (
	isServerStartedWithHttp = false
)

func RunWithQuitNotification(serverEngine *gin.Engine) {
	logServerStartDetails(serverEngine)

	useTLS := os.Getenv("USE_TLS") == "true"
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	domainName := os.Getenv("DOMAIN_NAME")

	addr := ":80"
	if useTLS {
		addr = ":443"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           serverEngine.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	if useTLS {
		isServerStartedWithHttp = true
		if certFile != "" && keyFile != "" {
			logger.Info("Using custom TLS certificates")
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{
					loadCertificate(certFile, keyFile),
				},
				MinVersion: tls.VersionTLS12,
			}
		} else {
			logger.Info("Using Let's Encrypt auto-certification")
			manager := autocert.Manager{
				Cache:      autocert.DirCache(".certs"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(domainName),
			}
			srv.TLSConfig = &tls.Config{
				GetCertificate: manager.GetCertificate,
				MinVersion:     tls.VersionTLS12,
			}
			srv.Handler = manager.HTTPHandler(serverEngine.Handler())
		}
	}

	logger.Infof("Starting server at %s", addr)

	go func() {
		var err error
		if useTLS {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error running the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down all registered shutdown handlers...")
	GracefulShutdownRegistry.TriggerShutdowns()
	logger.Info("Successfully triggered all registered shutdown channels.")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Info("Server gracefully stopped")
}

func loadCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logger.Fatalf("Failed to load certificates: %v", err)
	}
	return cert
}
