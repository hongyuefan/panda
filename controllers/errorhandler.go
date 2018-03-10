package controllers

import (
	"github.com/astaxie/beego/context"
)

type Errortail struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"required"`
}
type ErrorMsg struct {
	Message string    `json:"message"`
	Errors  Errortail `json:"errors"`
}

func ErrorHandler(ctx *context.Context, err error) {
	msg := &ErrorMsg{
		Message: err.Error(),
		Errors:  Errortail{},
	}
	ctx.Output.JSON(msg, false, false)
	return
}
