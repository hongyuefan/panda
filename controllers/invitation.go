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
		rspInvitationCode        RspInvitationCode
	)

	conf := GetConfigData()

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

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

	invitationUrl = conf.HostUrl + "/regist?code=" + code

	if img, err = GenQRCode(invitationUrl, 300, 300); err != nil {
		goto errDeal
	}
	rspInvitationCode = RspInvitationCode{
		InvitationUrl: invitationUrl,
		Surplus:       surplus,
		IsReward:      flag,
		Image:         img,
	}

	SuccessHandler(c.Ctx, rspInvitationCode)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func GenQRCode(code string, width, height int) (encode string, err error) {

	var sbarcode barcode.Barcode

	prepng := "data:image/png;base64,"

	if sbarcode, err = qr.Encode(code, qr.L, qr.Unicode); err != nil {
		return
	}

	if sbarcode, err = barcode.Scale(sbarcode, width, height); err != nil {
		return
	}

	c := new(bytes.Buffer)

	if err = png.Encode(c, sbarcode); err != nil {
		return
	}

	encode = base64.StdEncoding.EncodeToString(c.Bytes())

	encode = prepng + encode

	return
}
