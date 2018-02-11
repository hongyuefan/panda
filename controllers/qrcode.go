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
	Success bool   `json:"success"`
	Message string `json:"message"`
	PngCode string `json:"qrcode"`
}

func (q *QRCodeController) GenCode() {
	var (
		err     error
		barcode barcode.Barcode
		c       *bytes.Buffer
		code    string
		prepng  string
		userId  int64
		mUser   *models.Player
	)

	if userId, err = TokenValidate(q.Ctx.Input.Header("token")); err != nil {
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

	q.HandlerResult(true, prepng+code, "")

	return

errDeal:
	q.HandlerResult(false, "", err.Error())
	return
}

func (q *QRCodeController) HandlerResult(success bool, code string, message string) {

	rspQRCode := RspQRCode{
		Success: success,
		Message: message,
		PngCode: code,
	}

	q.Ctx.Output.JSON(rspQRCode, false, false)

	return
}
