package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

func Success(ctx *gin.Context, data gin.H) {
	Response(ctx, http.StatusOK, 200, data, "success")
}

func Fail(ctx *gin.Context, data gin.H) {
	Response(ctx, http.StatusInternalServerError, 500, data, "fail")
}
