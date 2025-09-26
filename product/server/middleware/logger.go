package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware creates a custom logger middleware with IP, method, path, status, and timing
func LoggerMiddleware() gin.HandlerFunc {
	// Create or open log file
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		// Fall back to stdout if log file creation fails
		logFile = os.Stdout
	}

	logger := log.New(logFile, "", log.LstdFlags)

	return gin.HandlerFunc(func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get request method
		method := c.Request.Method

		// Get status code
		statusCode := c.Writer.Status()

		// Get response size
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log format: [TIMESTAMP] IP | STATUS | LATENCY | METHOD | PATH | SIZE | USER-AGENT
		logMsg := fmt.Sprintf("%s | %d | %v | %s | %s | %d bytes | %s",
			clientIP,
			statusCode,
			latency,
			method,
			path,
			bodySize,
			userAgent,
		)

		// Color coding based on status code (for console output)
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = "\033[32m" // Green
		case statusCode >= 300 && statusCode < 400:
			statusColor = "\033[33m" // Yellow
		case statusCode >= 400 && statusCode < 500:
			statusColor = "\033[31m" // Red
		case statusCode >= 500:
			statusColor = "\033[35m" // Magenta
		default:
			statusColor = "\033[0m" // Reset
		}
		resetColor := "\033[0m"

		// Log to file (without colors)
		logger.Printf("[REQUEST] %s", logMsg)

		// Also log to console with colors (optional)
		if gin.Mode() == gin.DebugMode {
			fmt.Printf("[GIN] %s[%d]%s %s | %v | %s | %s\n",
				statusColor, statusCode, resetColor,
				clientIP, latency, method, path)
		}
	})
}

// ErrorLoggerMiddleware logs errors with more detail
func ErrorLoggerMiddleware() gin.HandlerFunc {
	// Create error log file
	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open error log file: %v", err)
		errorFile = os.Stderr
	}

	errorLogger := log.New(errorFile, "[ERROR] ", log.LstdFlags)

	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				errorLogger.Printf("IP: %s | Method: %s | Path: %s | Error: %s",
					c.ClientIP(),
					c.Request.Method,
					c.Request.URL.Path,
					ginErr.Error())
			}
		}
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	})
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// InitLogDirectories creates log directories if they don't exist
func InitLogDirectories() error {
	return os.MkdirAll("logs", 0755)
}