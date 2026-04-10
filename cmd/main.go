package main

import (
	"LiveChat/config"
	"LiveChat/db"
	"LiveChat/internal/ai"
	"LiveChat/internal/handler"
	"LiveChat/internal/repository"
	"LiveChat/internal/service"
	"LiveChat/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	userRep := repository.NewRepository(dbConn.GetDB())
	userSvc := service.NewService(userRep)
	userHandler := handler.NewHandler(userSvc)

	convRepo := repository.NewConversationRepository(dbConn.GetDB())
	msgRepo := repository.NewMessageRepository(dbConn.GetDB())
	aiSettingsRepo := repository.NewAISettingsRepository(dbConn.GetDB())
	convSvc := service.NewConversationService(convRepo, msgRepo, aiSettingsRepo)
	msgSvc := service.NewMessageService(msgRepo, convRepo)
	convHandler := handler.NewConversationHandler(convSvc, msgSvc)

	aiClient := ai.NewClient(cfg.OpenRouterAPIKey, cfg.OpenRouterBaseURL)
	aiSvc := service.NewAIService(aiClient, msgRepo, aiSettingsRepo, convRepo)
	aiHandler := handler.NewAIHandler(aiClient, aiSvc)

	router.InitRouter(userHandler, convHandler, aiHandler, cfg.CORSOrigins)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Start("0.0.0.0:" + port)
}
