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
	)

	query := make(map[string]string, 0)

	result, err := models.GetAllAgree(query, []string{"title", "content"}, []string{"id"}, []string{"asc"}, 0, 100)
	if err != nil {
		goto errDeal
	}

	for _, v := range result {
		sub.Text = v.(map[string]interface{})["title"].(string)
		sub.Standard = v.(map[string]interface{})["content"].(string)

		subContents = append(subContents, sub)
	}

	content.Text = "协议规则"
	content.Subtages = subContents

	a.Ctx.Output.JSON(content, false, false)

	return

errDeal:
	ErrorHandler(a.Ctx, err)
	return
}
