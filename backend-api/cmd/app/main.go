package backend_api

import (
	handler "x-ui/backend-api/internal/handlers"
	"x-ui/backend-api/internal/repository"
	"x-ui/backend-api/internal/service"
	"x-ui/database"
)

func Start() {
	db := database.GetDB()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	handlers.InitRoutes("8080")
}
