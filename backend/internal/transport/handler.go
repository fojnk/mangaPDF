package transport

import (
	_ "github.com/fojnk/Task-Test-devBack/docs"

	"github.com/fojnk/Task-Test-devBack/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api/v1", h.userIdentity)
	{
		api.GET("/account", h.getAccountInfo)

		manga := api.Group("/manga")
		{
			manga.GET("/list", h.getManga)
			manga.GET("/chapters", h.getMangaChapters)
			manga.POST("/download", h.downloadMangaChapters)
			manga.GET("/status", h.downloadStatus)
			manga.GET("/result", h.downloadResult)
		}
	}

	return router
}
