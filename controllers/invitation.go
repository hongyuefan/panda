package controllers

import (
	"panda/arithmetic"
	"panda/models"

	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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
	Image         string `json:"image"`
}

func (c *InvitationController) HandlerGenerateInvitationCode() {
	var (
		err                      error
		userId                   int64
		surplus                  int
		count                    int
		code, invitationUrl, img string
		flag                     int
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
		} else {
			goto errDeal
		}
	}
	if flag == 1 || count >= conf.InvitationLimit {
		surplus = 0
	} else {
		surplus = conf.InvitationLimit - count
	}

	invitationUrl = conf.HostUrl + "/v1/tsxm/regist?code=" + code

	if img, err = GenQRCode(invitationUrl); err != nil {
		goto errDeal
	}
	c.Ctx.Output.JSON(RspInvitationCode{
		Success:       true,
		InvitationUrl: invitationUrl,
		Surplus:       surplus,
		IsReward:      flag,
		Image:         img,
	}, false, false)
	return
errDeal:
	c.Ctx.Output.JSON(&RspInvitationCode{
		Success: false,
		Message: err.Error(),
	}, false, false)
	return
}

func GenQRCode(code string) (encode string, err error) {

	var barcode barcode.Barcode

	prepng := "data:image/png;base64,"

	if barcode, err = qr.Encode(code, qr.L, qr.Unicode); err != nil {
		return
	}
	c := new(bytes.Buffer)

	if err = png.Encode(c, barcode); err != nil {
		return
	}

	encode = base64.StdEncoding.EncodeToString(c.Bytes())

	encode = prepng + encode

	return
}
