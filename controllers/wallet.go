package controllers

import (
	"panda/models"
	t "panda/transaction"
	"strings"

	"github.com/astaxie/beego"
)

type WalletController struct {
	beego.Controller
}

type RspWallet struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (w *WalletController) SetWalletAddress() {
	var (
		err    error
		mUser  models.Player
		orm    *models.Common
		wallet string
	)

	if mUser.Id, err = ParseAndValidToken(w.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if err = w.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	wallet = w.Ctx.Request.FormValue("wallet")

	wallet = strings.TrimPrefix(strings.ToUpper(wallet), "0X")

	if err = t.ValidatePublicKey(wallet); err != nil {
		goto errDeal
	}

	mUser.Pubkey = wallet

	orm = models.NewCommon()

	if _, err = orm.CommonUpdateById(&mUser, "pubkey"); err != nil {
		goto errDeal
	}
	w.HandlerResult(true, "设置钱包成功")
	return
errDeal:
	w.HandlerResult(false, err.Error())
	return
}

func (w *WalletController) HandlerResult(success bool, message string) {
	rspWallet := RspWallet{
		Success: success,
		Message: message,
	}
	w.Ctx.Output.JSON(rspWallet, false, false)
	return
}