package controllers

import (
	"fmt"
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

	query := make(map[string]string)

	result, err := models.GetAllAgree(query, []string{}, []string{"id"}, []string{"asc"}, 0, 100)
	if err != nil {
		goto errDeal
	}
	fmt.Println(result)

	for _, v := range result {
		sub.Text = v.(models.Agree).Title
		sub.Standard = v.(models.Agree).Content

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
