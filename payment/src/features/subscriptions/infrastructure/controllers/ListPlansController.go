package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/features/subscriptions/application"
)

type ListPlansController struct {
	useCase *application.ListPlansUseCase
}

func NewListPlansController(useCase *application.ListPlansUseCase) *ListPlansController {
	return &ListPlansController{useCase: useCase}
}

func (ctrl *ListPlansController) Handle(c *gin.Context) {
	plans, err := ctrl.useCase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}
