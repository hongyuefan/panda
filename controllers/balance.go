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
}

func (b *BalanceConroller) GetBalance() {
	var (
		userId     int64
		err        error
		balance    string
		orm        *models.Common
		mUser      models.Player
		rspBalance RspBalance
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

	rspBalance = RspBalance{
		Balance: balance,
	}

	SuccessHandler(b.Ctx, rspBalance)
	return
errDeal:
	ErrorHandler(b.Ctx, err)
	return
}
