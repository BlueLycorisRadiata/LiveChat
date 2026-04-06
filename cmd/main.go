package main

import (
	"LiveChat/config"
	"LiveChat/db"
	"LiveChat/internal/ai"
	"LiveChat/internal/handler"
	"LiveChat/internal/repository"
	"LiveChat/internal/service"
	"LiveChat/router"
	"LiveChat/ws"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env (ignore error if file not present — env vars may be set externally)
	_ = godotenv.Load()

	cfg := config.Load()

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	userRep := repository.NewRepository(dbConn.GetDB())
	userSvc := service.NewService(userRep)
	userHandler := handler.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	convRepo := repository.NewConversationRepository(dbConn.GetDB())
	msgRepo := repository.NewMessageRepository(dbConn.GetDB())
	convSvc := service.NewConversationService(convRepo, msgRepo)
	msgSvc := service.NewMessageService(msgRepo, convRepo)
	convHandler := handler.NewConversationHandler(convSvc, msgSvc)

	// AI wiring
	aiClient := ai.NewClient(cfg.OpenRouterAPIKey, cfg.OpenRouterBaseURL)
	aiSettingsRepo := repository.NewAISettingsRepository(dbConn.GetDB())
	aiSvc := service.NewAIService(aiClient, msgRepo, aiSettingsRepo, convRepo)
	aiHandler := handler.NewAIHandler(aiClient, aiSvc)

	router.InitRouter(userHandler, wsHandler, convHandler, aiHandler)
	router.Start("0.0.0.0:8080")
}
