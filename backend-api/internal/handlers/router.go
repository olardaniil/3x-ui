package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"x-ui/backend-api/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	//
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		inbounds := api.Group("/inbounds")
		{
			//inbounds.GET("/", h.GetInbounds)
			//inbounds.GET("/:id", h.GetInboundById)

			clients := inbounds.Group(":id/clients")
			{
				clients.GET("/", h.GetInboundClients)
				clients.POST("/", h.AddInboundClient)
			}
		}
	}

	//	Запуск сервера
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
