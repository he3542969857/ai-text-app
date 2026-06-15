package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-text-app/backend/internal/model"
)

// functionsHandler 返回所有可用功能(中译英/英译中/文本总结)。
func functionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.Functions())
}
