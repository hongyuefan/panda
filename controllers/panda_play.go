package controllers

import (
	"fmt"
	"panda/models"
	"panda/types"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

func AddBalance(uid int64, amount string) (err error) {
	var (
		player  *models.Player
		balance float64
		famount float64
	)
	if balance, err = strconv.ParseFloat(player.Balance, 10); err != nil {
		return
	}
	if famount, err = strconv.ParseFloat(amount, 10); err != nil {
		return
	}
	if famount == 0 {
		return
	}
	if player, err = models.GetPlayerById(uid); err != nil {
		return
	}
	balance += famount

	player.Balance = fmt.Sprintf("%v", balance)

	if err = models.UpdatePlayerById(player, "Balance"); err != nil {
		beego.BeeLogger.Info("AddBalance:", uid, amount)
	}
	return
}

func SubBalance(uid int64, amount string) (err error) {
	var (
		player  *models.Player
		balance float64
		famount float64
	)
	if balance, err = strconv.ParseFloat(player.Balance, 10); err != nil {
		return
	}
	if famount, err = strconv.ParseFloat(amount, 10); err != nil {
		return
	}
	if famount > balance {
		return types.Error_Player_Balance
	}
	if player, err = models.GetPlayerById(uid); err != nil {
		return
	}
	balance -= famount

	player.Balance = fmt.Sprintf("%v", balance)

	if err = models.UpdatePlayerById(player, "Balance"); err != nil {
		beego.BeeLogger.Info("SubBalance :", uid, amount)
	}
	return
}

func UpCatchTime(uid int64) (err error) {
	var (
		player *models.Player
	)
	if player, err = models.GetPlayerById(uid); err != nil {
		return
	}
	player.LastCatchTime = time.Now().Unix()

	if err = models.UpdatePlayerById(player, "LastCatchTime"); err != nil {
		beego.BeeLogger.Info("UpCatchTime:", uid, player.LastCatchTime)
	}
	return
}
