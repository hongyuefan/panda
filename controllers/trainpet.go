package controllers

import (
	"panda/types"
	"strconv"

	"github.com/astaxie/beego"
)

// CatchController operations for Catch
type TrainController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *TrainController) HandlerTrainPet() {
	var (
		userId                  int64
		err                     error
		spetId, sAmount, txhash string
		petId                   int64
		result                  types.RspTrain
	)
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	spetId = c.Ctx.Request.FormValue("petId")
	sAmount = c.Ctx.Request.FormValue("amount")

	if petId, err = strconv.ParseInt(spetId, 10, 64); err != nil {
		goto errDeal
	}

	if txhash, err = c.trans.Transactions(types.Trans_Type_Train, userId, petId, 0, sAmount); err != nil {
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
