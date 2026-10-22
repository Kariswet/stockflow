package util

import (
	"net/http"
	"stockflow/model"
	"time"

	"github.com/gin-gonic/gin"
)

func SetMetadataResponse(ctx *gin.Context, statusCode int, startTime time.Time, resp *model.Response) {
	resp.Metadata.Status = statusCode >= 200 && statusCode < 300

	if resp.Metadata.Message == "" {
		if resp.Metadata.Status {
			resp.Metadata.Message = "OK"
		} else {
			resp.Metadata.Message = http.StatusText(statusCode)
		}
	}

	resp.Metadata.TimeExecution = time.Since(startTime).String()

	ctx.JSON(statusCode, resp)
}
