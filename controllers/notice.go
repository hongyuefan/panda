package controllers

import (
	"panda/models"

	"github.com/astaxie/beego"
)

type NoticeController struct {
	beego.Controller
}

type Notice struct {
	Text string `json:"text"`
}

func (c *NoticeController) HandlerNotice() {
	var (
		err    error
		notice string
	)
	if notice, err = models.GetNoticPub(); err != nil {
		goto errDeal
	}
	c.Ctx.Output.JSON(&Notice{
		Text: notice,
	}, false, false)
	return
errDeal:
	c.Ctx.Output.JSON(&Notice{
		Text: err.Error(),
	}, false, false)
	return
}
