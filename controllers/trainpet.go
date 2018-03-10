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
	c.HandlerResult(true, nil, txhash)
	return
errDeal:
	c.HandlerResult(false, err, "")
	return
}

func (c *TrainController) HandlerResult(success bool, err error, hash string) {

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
