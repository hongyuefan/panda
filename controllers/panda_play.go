package controllers

import (
	"panda/models"

	"time"

	"github.com/astaxie/beego"
)

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
