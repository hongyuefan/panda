package controllers

import (
	"panda/models"
	"panda/types"

	"github.com/astaxie/beego"
)

type GamblingController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *GamblingController) HandlerGambling() {
	var (
		err    error
		userId int64
		txhash string
		result *types.RspCatch
	)
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	if err = models.IsHaveGamblingCount(userId); err != nil {
		goto errDeal
	}
	if txhash, err = c.trans.Transactions(types.Trans_Type_Gambling, userId, 0, 0, "0"); err != nil {
		goto errDeal
	}
	result = &types.RspCatch{
		Txhash: txhash,
	}
	SuccessHandler(c.Ctx, result)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
