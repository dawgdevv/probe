package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dawgdevv/probe/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all routes and middleware
func SetupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		projects := api.Group("/projects")
		{
			projects.GET("", handler.ListProjects)
			projects.POST("", handler.CreateProject)
			projects.GET("/:id", handler.GetProject)
			projects.GET("/:id/suites", handler.ListSuitesForProject)
		}

		suites := api.Group("/suites")
		{
			suites.POST("", handler.CreateTestSuite)
			suites.GET("/:id", handler.GetTestSuite)
			suites.POST("/:id/run", handler.RunTestSuite)
			suites.GET("/:id/runs", handler.ListTestRuns)
		}

		runs := api.Group("/runs")
		{
			runs.GET("/:id", handler.GetTestRun)
			runs.GET("/:id/results", handler.GetTestResults)
		}
	}

	// Serve embedded React UI with SPA fallback
	staticFS, err := web.StaticFS()
	if err != nil {
		fmt.Println("Warning: could not load embedded UI:", err)
		router.GET("/", handler.Landing)
	} else {
		indexHTML, _ := web.IndexHTML()

		// Serve static assets (js, css, etc.)
		router.GET("/assets/*filepath", func(c *gin.Context) {
			http.FileServer(staticFS).ServeHTTP(c.Writer, c.Request)
		})

		// Serve other static files (probe.svg, etc.)
		router.GET("/probe.svg", func(c *gin.Context) {
			http.FileServer(staticFS).ServeHTTP(c.Writer, c.Request)
		})

		// SPA fallback: serve index.html for all non-API, non-asset routes
		router.NoRoute(func(c *gin.Context) {
			// Don't intercept API routes
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}

			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
		})
	}

	return router
}
