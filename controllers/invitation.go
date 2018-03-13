package controllers

import (
	"panda/arithmetic"
	"panda/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type InvitationController struct {
	beego.Controller
}
type RspInvitationCode struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	InvitationUrl string `json:"invitation_url"`
	Surplus       int    `json:"surplus"`
	IsReward      int    `json:"isReward"`
}

func (c *InvitationController) HandlerGenerateInvitationCode() {
	var (
		err     error
		userId  int64
		surplus int
		count   int
		code    string
		flag    int
	)
	conf := GetConfigData()

	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	if code, count, flag, err = models.GetInvitationByUid(userId); err != nil {
		if err == orm.ErrNoRows {
			code = arithmetic.GenCode(8)
			if err = models.AddInvitation(userId, code); err != nil {
				goto errDeal
			}
		}
		goto errDeal
	}
	if flag == 1 || count >= conf.InvitationLimit {
		surplus = 0
	} else {
		surplus = conf.InvitationLimit - count
	}
	c.Ctx.Output.JSON(RspInvitationCode{
		Success:       true,
		InvitationUrl: conf.HostUrl + "/v1/tsxm/regist?code=" + code,
		Surplus:       surplus,
		IsReward:      flag,
	}, false, false)
	return
errDeal:
	c.Ctx.Output.JSON(&RspInvitationCode{
		Success: false,
		Message: err.Error(),
	}, false, false)
	return
}
