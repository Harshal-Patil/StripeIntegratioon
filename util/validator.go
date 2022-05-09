package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func Validate(ctx *gin.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	v := validator.New()
	return v.StructCtx(ctx.Request.Context(), request)
}
