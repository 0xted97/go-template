package routes

import (
	"github.com/gin-gonic/gin"

	"allaccessone/blockchains-support/modules/flow"
)

func InitFlowRoutes(route *gin.Engine) {
	/**
	@description Flow
	*/
	flowService := flow.NewFlowService()
	flowController := flow.NewController(*flowService)

	groupRoute := route.Group("/api/v1/flow")
	groupRoute.POST("/create-account", flowController.CreateFlowAccount)
}
