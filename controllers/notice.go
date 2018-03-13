package controllers

import (
	"panda/models"
	"strconv"

	"github.com/astaxie/beego"
)

type NoticeController struct {
	beego.Controller
}

type Notice struct {
	Text string `json:"text"`
}

type RspNews struct {
	Total int    `json:"total"`
	Data  []News `json:"data"`
}

type News struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Creatime string `json:"create_time"`
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

func (c *NoticeController) HandlerNews() {
	var (
		err             error
		news            []News
		one             News
		rspNews         RspNews
		spage, sperpage string
		sorder          string
		page, perpage   int64
		query           map[string]string
		ml              []interface{}
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	spage = c.Ctx.Request.FormValue("page")
	sperpage = c.Ctx.Request.FormValue("perpage")
	sorder = c.Ctx.Request.FormValue("order")

	if page, err = strconv.ParseInt(spage, 10, 64); err != nil {
		goto errDeal
	}
	if perpage, err = strconv.ParseInt(sperpage, 10, 64); err != nil {
		goto errDeal
	}

	query = make(map[string]string)

	query["flag"] = "0"

	if ml, err = models.GetNews(query, []string{}, []string{"id"}, []string{sorder}, page*perpage, perpage); err != nil {
		goto errDeal
	}

	for _, v := range ml {
		one.Id = v.(models.Notice).Id
		one.Creatime = v.(models.Notice).CreateTime
		one.Text = v.(models.Notice).Text
		one.Title = v.(models.Notice).Title

		news = append(news, one)
	}
	rspNews = RspNews{
		Total: len(ml),
		Data:  news,
	}
	c.Ctx.Output.JSON(rspNews, false, false)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
