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
		err    error
		result types.RspTrain
	)

	if _, err = c.trans.Transactions(types.Trans_Type_Bonus, 0, 0, 0, ""); err != nil {
		goto errDeal
	}
	result = types.RspTrain{
		Txhash: "",
	}
	SuccessHandler(c.Ctx, result)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
