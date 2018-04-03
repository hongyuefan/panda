package controllers

import (
	"panda/types"

	"github.com/astaxie/beego"
)

// CatchController operations for Catch
type WithDrawalController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *WithDrawalController) HandlerWithDrawal() {
	var (
		userId          int64
		err             error
		sAmount, txhash string
		result          types.RspTrain
	)
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	sAmount = c.Ctx.Request.FormValue("amount")

	if txhash, err = c.trans.Transactions(types.Trans_Type_WithDrawal, userId, 0, 0, sAmount); err != nil {
		goto errDeal
	}
	result = types.RspTrain{
		Txhash: txhash,
	}
	SuccessHandler(c.Ctx, result)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
