package main

import (
	"allaccessone/flow-service/routes"
	"allaccessone/flow-service/utils"
	"log"
	"net/http"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	/**
	@description Setup Server
	*/
	router := SetupRouter()
	/**
	@description Run Server
	*/
	log.Fatal(router.Run(":" + utils.GodotEnv("GO_PORT")))
}

func SetupRouter() *gin.Engine {
	/**
	@description Setup Database Connection
	*/
	// db := config.Connection()
	/**
	@description Init Router
	*/
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		utils.APIResponse(c, "Not found route", http.StatusNotFound, c.Request.Method, nil)
	})
	/**
	@description Setup Mode Application
	*/
	if utils.GodotEnv("GO_ENV") != "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	if utils.GodotEnv("GO_ENV") != "development" {
		gin.SetMode(gin.DebugMode)
	}
	if utils.GodotEnv("GO_ENV") != "test" {
		gin.SetMode(gin.TestMode)
	}
	/**
	@description Setup Middleware
	*/
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))
	/**
	@description Init All Route
	*/
	routes.InitFlowRoutes(router)

	return router
}
