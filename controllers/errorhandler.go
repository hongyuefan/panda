package controllers

import (
	"encoding/json"
	"panda/types"

	"github.com/astaxie/beego/context"
)

func ErrorHandler(ctx *context.Context, err error) {

	msg := types.RspCommon{
		Success: false,
		Message: err.Error(),
	}
	ctx.Output.JSON(msg, false, false)
	return
}

func SuccessHandler(ctx *context.Context, out interface{}) {

	jsonData, err := json.Marshal(out)
	if err != nil {
		ErrorHandler(ctx, err)
		return
	}

	msg := types.RspCommon{
		Success: true,
		Message: "success",
		Data:    jsonData,
	}
	ctx.Output.JSON(msg, false, false)
	return
}
