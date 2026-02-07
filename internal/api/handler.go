package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dawgdevv/apitestercli/internal/loader"
	"github.com/dawgdevv/apitestercli/internal/service"
	"github.com/dawgdevv/apitestercli/internal/storage"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Handler contains all HTTP handlers for the API
type Handler struct {
	store *storage.Store
}

// NewHandler creates a new API handler
func NewHandler(store *storage.Store) *Handler {
	return &Handler{store: store}
}

// HealthCheck returns server health status
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Landing serves a simple landing page
func (h *Handler) Landing(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>API Tester CLI</title>
    <style>
        body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 API Tester CLI Server</h1>
    <p>Server is running! The web UI will be available here once the frontend is built.</p>
    
    <h2>Available API Endpoints:</h2>
    <div class="endpoint"><strong>GET</strong> <code>/api/health</code> - Server health check</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/projects</code> - List all projects</div>
    <div class="endpoint"><strong>POST</strong> <code>/api/projects</code> - Create a new project</div>
    <div class="endpoint"><strong>POST</strong> <code>/api/suites</code> - Create a test suite</div>
    <div class="endpoint"><strong>POST</strong> <code>/api/suites/:id/run</code> - Run tests</div>
    
    <p style="margin-top: 30px; color: #666;">
        API Documentation: <a href="/api/health">/api/health</a>
    </p>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// --- Project Handlers ---

// CreateProject creates a new project
func (h *Handler) CreateProject(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.store.CreateProject(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProject retrieves a project by ID
func (h *Handler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	project, err := h.store.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// ListProjects retrieves all projects
func (h *Handler) ListProjects(c *gin.Context) {
	projects, err := h.store.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// --- Test Suite Handlers ---

// CreateTestSuite creates a new test suite from YAML
func (h *Handler) CreateTestSuite(c *gin.Context) {
	var req struct {
		ProjectID   int64  `json:"project_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		YAMLContent string `json:"yaml_content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate YAML by attempting to parse it
	var testSuite interface{}
	if err := yaml.Unmarshal([]byte(req.YAMLContent), &testSuite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid YAML: " + err.Error()})
		return
	}

	suite, err := h.store.CreateTestSuite(req.ProjectID, req.Name, req.YAMLContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, suite)
}

// GetTestSuite retrieves a test suite by ID
func (h *Handler) GetTestSuite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid suite ID"})
		return
	}

	suite, err := h.store.GetTestSuite(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "test suite not found"})
		return
	}

	c.JSON(http.StatusOK, suite)
}

// ListSuitesForProject lists all test suites for a project
func (h *Handler) ListSuitesForProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	suites, err := h.store.ListTestSuites(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suites": suites})
}

// RunTestSuite executes a test suite and stores results
func (h *Handler) RunTestSuite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid suite ID"})
		return
	}

	// Get suite from database
	storedSuite, err := h.store.GetTestSuite(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "test suite not found"})
		return
	}

	// Parse YAML content
	suite, err := loader.LoadSuiteFromString(storedSuite.YAMLContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse test suite: " + err.Error()})
		return
	}

	// Create test run record
	testRun, err := h.store.CreateTestRun(id, len(suite.Tests))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create test run: " + err.Error()})
		return
	}

	// Execute tests using service layer
	runner := service.NewRunner(service.RunOptions{
		MaxConcurrent: 10,
	})

	results, err := runner.RunSuite(suite)
	if err != nil {
		// Mark run as error
		h.store.CompleteTestRun(testRun.ID, "error", 0, 0)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save individual test results
	passed := 0
	failed := 0
	for _, result := range results {
		errorMsg := ""
		if result.Error != nil {
			errorMsg = result.Error.Error()
			failed++
		} else {
			passed++
		}

		durationMs := result.Duration.Milliseconds()
		if err := h.store.SaveTestResult(testRun.ID, result.Name, result.Passed, result.StatusCode, errorMsg, durationMs); err != nil {
			fmt.Printf("Warning: failed to save test result: %v\n", err)
		}
	}

	// Update test run with final status
	status := "passed"
	if failed > 0 {
		status = "failed"
	}
	if err := h.store.CompleteTestRun(testRun.ID, status, passed, failed); err != nil {
		fmt.Printf("Warning: failed to complete test run: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"run_id":       testRun.ID,
		"status":       status,
		"total_tests":  len(suite.Tests),
		"passed_tests": passed,
		"failed_tests": failed,
		"results":      results,
	})
}

// --- Test Run Handlers ---

// GetTestRun retrieves a test run by ID
func (h *Handler) GetTestRun(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run ID"})
		return
	}

	run, err := h.store.GetTestRun(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "test run not found"})
		return
	}

	c.JSON(http.StatusOK, run)
}

// ListTestRuns lists test runs for a suite
func (h *Handler) ListTestRuns(c *gin.Context) {
	suiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid suite ID"})
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	runs, err := h.store.ListTestRuns(suiteID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"runs": runs})
}

// GetTestResults retrieves all results for a test run
func (h *Handler) GetTestResults(c *gin.Context) {
	runID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run ID"})
		return
	}

	results, err := h.store.GetTestResults(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
