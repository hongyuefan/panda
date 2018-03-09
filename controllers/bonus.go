package controllers

import (
	"panda/types"

	"github.com/astaxie/beego"
)

type BonusController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *BonusController) HandlerBonus() {
	var (
		err error
	)

	if _, err = c.trans.Transactions(types.Trans_Type_Bonus, 0, 0, ""); err != nil {
		goto errDeal
	}
	c.HandlerResult(true, nil, "启动分红成功")
	return
errDeal:
	c.HandlerResult(false, err, "")
	return
}

func (c *BonusController) HandlerResult(success bool, err error, hash string) {

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
