package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafif/healy-backend/internal/delivery/websocket"
	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/usecase"
	"github.com/rafif/healy-backend/pkg/config"
)

// SetupRouter creates and configures the Gin engine with all the routes.
func SetupRouter(cfg *config.Config, hub *websocket.Hub, telemetryUsecase usecase.TelemetryUsecase, authUsecase usecase.AuthUsecase) *gin.Engine {
	r := gin.Default()

	// Configure CORS (basic example, adjust in production)
	r.Use(func(c *gin.Context) {
		// Example using cfg.CORSAllowedOrigins. For a proper implementation, use github.com/gin-contrib/cors
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, device_id")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				var req domain.LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				resp, err := authUsecase.Login(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, resp)
			})
		}

		// Group that requires JWT validation
		// TODO: Add JWT middleware here: protected.Use(middleware.JWTAuth())
		protected := api.Group("/")
		{
			telemetry := protected.Group("/telemetry")
			{
				telemetry.GET("/history", func(c *gin.Context) {
					// TODO: Implement history retrieval via telemetryUsecase
					c.JSON(http.StatusOK, gin.H{"message": "History endpoint"})
				})
				telemetry.GET("/latest", func(c *gin.Context) {
					// TODO: Implement latest retrieval via telemetryUsecase
					c.JSON(http.StatusOK, gin.H{"message": "Latest endpoint"})
				})
			}

			settings := protected.Group("/settings")
			{
				settings.PUT("/threshold", func(c *gin.Context) {
					// TODO: Implement threshold update
					c.JSON(http.StatusOK, gin.H{"message": "Threshold updated"})
				})
				settings.GET("/threshold", func(c *gin.Context) {
					// TODO: Implement threshold retrieval
					c.JSON(http.StatusOK, gin.H{"message": "Threshold info"})
				})
			}

			device := protected.Group("/device")
			{
				device.GET("/status", func(c *gin.Context) {
					// TODO: Implement device status check
					c.JSON(http.StatusOK, gin.H{"status": "CONNECTED"})
				})
			}
		}
	}

	// WebSocket endpoints
	ws := r.Group("/ws")
	{
		ws.GET("/telemetry", func(c *gin.Context) {
			// JWT should ideally be validated here using token from query param
			// e.g. token := c.Query("token")
			websocket.ServeViewerWs(hub, c.Writer, c.Request)
		})

		ws.GET("/device", func(c *gin.Context) {
			// device_id is expected in header
			websocket.ServeDeviceWs(hub, telemetryUsecase, c.Writer, c.Request)
		})
	}

	return r
}
