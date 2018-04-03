package controllers

import (
	"bytes"
	"encoding/base64"
	//	"io"
	"panda/models"

	"image/png"

	"github.com/astaxie/beego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QRCodeController struct {
	beego.Controller
}

type RspQRCode struct {
	PngCode  string `json:"qrcode"`
	WordCode string `json:"wcode"`
}

func (q *QRCodeController) GenCode() {
	var (
		err       error
		barcode   barcode.Barcode
		c         *bytes.Buffer
		code      string
		prepng    string
		userId    int64
		mUser     *models.Player
		token     string
		rspQRCode RspQRCode
	)

	if token, err = ParseToken(q.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if userId, err = TokenValidate(token); err != nil {
		goto errDeal
	}

	if mUser, err = models.GetPlayerById(userId); err != nil {
		goto errDeal
	}

	prepng = "data:image/png;base64,"

	if barcode, err = qr.Encode(mUser.PubPublic, qr.L, qr.Unicode); err != nil {
		goto errDeal
	}
	c = new(bytes.Buffer)

	if err = png.Encode(c, barcode); err != nil {
		goto errDeal
	}

	code = base64.StdEncoding.EncodeToString(c.Bytes())

	rspQRCode = RspQRCode{
		PngCode:  prepng + code,
		WordCode: mUser.PubPublic,
	}
	SuccessHandler(q.Ctx, rspQRCode)
	return
errDeal:
	ErrorHandler(q.Ctx, err)
	return
}
