package flow

import (
	"net/http"

	"allaccessone/blockchains-support/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type controller struct {
	service FlowService
}

func NewController(service FlowService) *controller {
	return &controller{service: service}
}

func (c *controller) CreateFlowAccount(ctx *gin.Context) {
	var input CreateFlowAccountRequest
	ctx.ShouldBindJSON(&input)
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.APIResponse(ctx, validationErrors.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}
	account, err := c.service.CreateFlowAccount(input)
	if err != nil {
		utils.APIResponse(ctx, err.Error(), http.StatusBadRequest, http.MethodPost, nil)
		return
	}
	utils.APIResponse(ctx, "Register new account successfully", http.StatusCreated, http.MethodPost, account)
	return
}
