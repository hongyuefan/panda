package controllers

import (
	"panda/models"
	"strconv"

	"github.com/astaxie/beego"
)

type BalanceConroller struct {
	beego.Controller
}

type RspBalance struct {
	Balance float64 `json:"balance"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

func (b *BalanceConroller) GetBalance() {
	var (
		token   string
		userId  int64
		err     error
		balance float64
		orm     *models.Common
		mUser   models.Player
	)
	token = b.Ctx.Input.Header("Authorization")

	if token, err = ParseToken(token); err != nil {
		goto errDeal
	}
	if userId, err = TokenValidate(token); err != nil {
		goto errDeal
	}

	mUser.Id = userId

	orm = models.NewCommon()

	if err = orm.CommonGetOne(&mUser, "id"); err != nil {
		goto errDeal
	}

	if balance, err = strconv.ParseFloat(mUser.Balance, 10); err != nil {
		goto errDeal
	}

	b.HandlerResult(balance, true, "")
	return
errDeal:
	b.HandlerResult(0, false, err.Error())
	return
}

func (b *BalanceConroller) HandlerResult(balance float64, success bool, message string) {

	rspBalance := RspBalance{
		Balance: balance,
		Success: success,
		Message: message,
	}

	b.Ctx.Output.JSON(rspBalance, false, false)

	return
}
