package main

import (
	"LiveChat/db"
	"LiveChat/internal/handler"
	"LiveChat/internal/repository"
	"LiveChat/internal/service"
	"LiveChat/router"
	"LiveChat/ws"
	"log"
)

func main() {
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

	router.InitRouter(userHandler, wsHandler, convHandler)
	router.Start("0.0.0.0:8080")
}
