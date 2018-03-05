package controllers

import (
	"panda/models"
	t "panda/transaction"

	"github.com/astaxie/beego"
)

type BalanceConroller struct {
	beego.Controller
}

type RspBalance struct {
	Balance string `json:"balance"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (b *BalanceConroller) GetBalance() {
	var (
		userId  int64
		err     error
		balance string
		orm     *models.Common
		mUser   models.Player
	)

	if userId, err = ParseAndValidToken(b.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	mUser.Id = userId

	orm = models.NewCommon()

	if err = orm.CommonGetOne(&mUser, "id"); err != nil {
		goto errDeal
	}

	if balance, err = t.GetBalance(mUser.PubPublic); err != nil {
		goto errDeal
	}

	b.HandlerResult(balance, true, "")
	return
errDeal:
	b.HandlerResult("0", false, err.Error())
	return
}

func (b *BalanceConroller) HandlerResult(balance string, success bool, message string) {

	rspBalance := RspBalance{
		Balance: balance,
		Success: success,
		Message: message,
	}

	b.Ctx.Output.JSON(rspBalance, false, false)

	return
}
