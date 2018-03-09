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
	)
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	sAmount = c.Ctx.Request.FormValue("amount")

	if txhash, err = c.trans.Transactions(types.Trans_Type_WithDrawal, userId, 0, sAmount); err != nil {
		goto errDeal
	}
	c.HandlerResult(true, nil, txhash)
	return
errDeal:
	c.HandlerResult(false, err, "")
	return
}

func (c *WithDrawalController) HandlerResult(success bool, err error, hash string) {

	var msg string

	if err != nil {
		msg = err.Error()
	}
	result := &types.RspTrain{
		Success: success,
		Message: msg,
		Txhash:  hash,
	}
	c.Ctx.Output.JSON(result, false, false)
	return
}
