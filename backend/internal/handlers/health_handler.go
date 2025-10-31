package handlers

import (
	"net/http"
	"runtime"
	"time"

	"filmfolk/internal/db"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks"`
}

type HealthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

var startTime = time.Now()

// HealthCheck returns basic health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "filmfolk-api",
	})
}

// DetailedHealthCheck returns comprehensive health information
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	checks := make(map[string]HealthCheck)

	// Database health check
	dbStatus := h.checkDatabase()
	checks["database"] = dbStatus

	// Memory health check
	memStatus := h.checkMemory()
	checks["memory"] = memStatus

	// Determine overall status
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status != "healthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	// Calculate uptime
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "filmfolk-api",
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// checkDatabase verifies database connectivity
func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()

	sqlDB, err := db.DB.DB()
	if err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "Failed to get database instance: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "Database ping failed: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	return HealthCheck{
		Status:  "healthy",
		Message: "Database connection is healthy",
		Latency: time.Since(start).String(),
	}
}

// checkMemory checks memory usage
func (h *HealthHandler) checkMemory() HealthCheck {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert bytes to MB
	allocMB := m.Alloc / 1024 / 1024
	sysMB := m.Sys / 1024 / 1024

	status := "healthy"
	message := "Memory usage is normal"

	// Alert if memory usage is high (>500MB allocated)
	if allocMB > 500 {
		status = "warning"
		message = "High memory usage detected"
	}

	return HealthCheck{
		Status:  status,
		Message: message + " (Alloc: " + string(rune(allocMB)) + "MB, Sys: " + string(rune(sysMB)) + "MB)",
	}
}

// ReadinessCheck returns whether the service is ready to accept traffic
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check critical dependencies
	dbCheck := h.checkDatabase()

	if dbCheck.Status != "healthy" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":   false,
			"message": "Service not ready: " + dbCheck.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready":   true,
		"message": "Service is ready to accept traffic",
	})
}

// LivenessCheck returns whether the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive": true,
	})
}
