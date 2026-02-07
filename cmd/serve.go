package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dawgdevv/apitestercli/internal/api"
	"github.com/dawgdevv/apitestercli/internal/storage"
	"github.com/spf13/cobra"
)

var (
	port    int
	dataDir string
)

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8443, "Port to run the server on")
	serveCmd.Flags().StringVarP(&dataDir, "data-dir", "d", "", "Directory to store data (default: ~/.apitester)")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API Tester web server",
	Long:  `Starts an HTTPS server with REST API and embedded web UI for managing and running API tests`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set default data directory
		if dataDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				os.Exit(1)
			}
			dataDir = filepath.Join(home, ".apitester")
		}

		// Initialize storage
		store, err := storage.NewStore(dataDir)
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer store.Close()

		// Create API handler
		handler := api.NewHandler(store)

		// Setup router
		router := api.SetupRouter(handler)

		// Generate self-signed certificate
		cert, err := generateSelfSignedCert()
		if err != nil {
			fmt.Println("Error generating certificate:", err)
			os.Exit(1)
		}

		// Create HTTPS server
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		// Start server in goroutine
		go func() {
			fmt.Printf("\n🚀 API Tester Server running on https://localhost:%d\n", port)
			fmt.Printf("📁 Data directory: %s\n", dataDir)
			fmt.Printf("\n Press Ctrl+C to stop\n\n")

			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				fmt.Println("Error starting server:", err)
				os.Exit(1)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("\n🛑 Shutting down server...")

		// Graceful shutdown with 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("Server forced to shutdown:", err)
		}

		fmt.Println("✅ Server stopped")
	},
}

// generateSelfSignedCert creates a self-signed TLS certificate for HTTPS
func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"API Tester CLI"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return tls.X509KeyPair(certPEM, keyPEM)
}
