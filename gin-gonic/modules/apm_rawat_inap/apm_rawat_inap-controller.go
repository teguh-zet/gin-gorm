package apm_rawat_inap

import (
	"gin-gonic/helper"

	"github.com/gin-gonic/gin"
)

type ApmRawatInapController struct {
	service ApmRawatInapService
}

func NewApmRawatInapController(service ApmRawatInapService) *ApmRawatInapController {
	return &ApmRawatInapController{service: service}
}

func (c *ApmRawatInapController) Create(ctx *gin.Context) {
	var input ApmRawatInapCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		helper.ValidationError(ctx, err.Error())
		return
	}
	res, err := c.service.Create(&input)
	if err != nil {
		helper.InternalServerError(ctx, "Failed to create APM Rawat Inap", err.Error())
		return
	}
	helper.CreatedResponse(ctx, "APM Rawat Inap created successfully", res)
}

func (c *ApmRawatInapController) Publish(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.Publish(id)
	if err != nil {
		helper.InternalServerError(ctx, "Failed to publish APM Rawat Inap", err.Error())
		return
	}
	helper.SuccessResponse(ctx, "APM Rawat Inap published successfully", res)
}
