package controllers

import (
	"panda/models"

	"github.com/astaxie/beego"
)

type AgreeContoller struct {
	beego.Controller
}

type Content struct {
	Text     string       `json:"text"`
	Subtages []SubContent `json:"subtags"`
}

type SubContent struct {
	Text     string `json:"text"`
	Standard string `json:"subtags"`
}

func (a *AgreeContoller) GetAgreement() {

	var (
		content     Content
		subContents []SubContent
		sub         SubContent
		result      []interface{}
		types       string
		err         error
	)

	query := make(map[string]string)

	if err = a.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	types = a.Ctx.Request.FormValue("type")

	query["stype"] = types

	if result, err = models.GetAllAgree(query, []string{}, []string{"id"}, []string{"asc"}, 0, 1000); err != nil {
		goto errDeal
	}

	for _, v := range result {
		sub.Text = v.(models.Agree).Title

		sub.Standard = v.(models.Agree).Content

		subContents = append(subContents, sub)
	}

	switch types {
	case "help":
		content.Text = "帮助中心"
	case "luckdraw":
		content.Text = "抽奖协议"
	case "regist":
		content.Text = "注册协议"
	}
	content.Subtages = subContents

	SuccessHandler(a.Ctx, content)

	return
errDeal:
	ErrorHandler(a.Ctx, err)
	return
}
