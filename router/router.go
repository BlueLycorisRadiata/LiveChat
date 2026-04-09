package router

import (
	"LiveChat/internal/handler"
	"LiveChat/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *handler.Handler, convHandler *handler.ConversationHandler, aiHandler *handler.AIHandler, corsOrigins []string) {
	r = gin.Default()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	if len(corsOrigins) == 0 {
		corsOrigins = []string{"http://localhost:5173"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", userHandler.Logout)

	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.POST("/conversations", convHandler.CreateConversation)
		auth.GET("/conversations", convHandler.GetConversations)
		auth.GET("/conversations/:id", convHandler.GetConversation)
		auth.PATCH("/conversations/:id", convHandler.UpdateConversation)
		auth.DELETE("/conversations/:id", convHandler.DeleteConversation)
		auth.POST("/conversations/:id/leave", convHandler.LeaveConversation)
		auth.GET("/conversations/:id/messages", convHandler.GetMessages)
		auth.POST("/conversations/:id/messages", convHandler.SendMessage)
		auth.DELETE("/conversations/:id/messages/:messageId", convHandler.DeleteMessage)

		// AI routes
		auth.GET("/ai/models", aiHandler.ListModels)
		auth.POST("/conversations/:id/ai/stream", aiHandler.StreamMessage)
		auth.GET("/conversations/:id/ai-settings", aiHandler.GetSettings)
		auth.PATCH("/conversations/:id/ai-settings", aiHandler.UpdateSettings)
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
