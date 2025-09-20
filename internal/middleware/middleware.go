package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or get user thread ID
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = uuid.New().String()
		}
		c.Set("userID", userID)
		c.Set("userThreadID", userID)

		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for admin token/header
		adminToken := c.GetHeader("X-Admin-Token")
		if adminToken == "" {
			// For development, allow all requests
			// In production, validate against configured admin token
			c.Next()
			return
		}

		// TODO: Implement proper admin authentication
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}