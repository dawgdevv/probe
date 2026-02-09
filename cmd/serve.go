package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dawgdevv/probe/internal/api"
	"github.com/dawgdevv/probe/internal/storage"
	"github.com/spf13/cobra"
)

var (
	port    int
	dataDir string
)

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	serveCmd.Flags().StringVarP(&dataDir, "data-dir", "d", "", "Directory to store data (default: ~/.probe)")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Probe web server",
	Long:  `Starts an HTTP server with REST API and embedded web UI for managing and running API tests`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set default data directory
		if dataDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				os.Exit(1)
			}
			dataDir = filepath.Join(home, ".probe")
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

		// Create HTTP server
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		// Start server in goroutine
		go func() {
			fmt.Printf("\n🚀 Probe running on http://localhost:%d\n", port)
			fmt.Printf("📁 Data directory: %s\n", dataDir)
			fmt.Printf("\n   Press Ctrl+C to stop\n\n")

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
