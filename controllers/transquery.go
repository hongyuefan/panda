package controllers

import (
	"fmt"
	"panda/models"
	"strconv"

	"panda/types"

	"github.com/astaxie/beego"
)

type TransQContoller struct {
	beego.Controller
}

func (t *TransQContoller) GetTransQ() {
	var (
		err                            error
		data                           types.TransQData
		datas                          []types.TransQData
		userId                         int64
		spage, sperpage, sorder, stype string
		page, perpage                  int64
		query                          map[string]string
		ml                             []interface{}
		count                          int
		conf                           models.Config
	)

	if err = t.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	if userId, err = ParseAndValidToken(t.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	query = make(map[string]string, 0)
	query["uid"] = fmt.Sprintf("%v", userId)
	spage = t.Ctx.Request.FormValue("page")
	sperpage = t.Ctx.Request.FormValue("perpage")
	sorder = t.Ctx.Request.FormValue("order")
	stype = t.Ctx.Request.FormValue("type")
	if stype != "" && stype != "0" {
		query["type"] = fmt.Sprintf("%v", stype)
	}

	if page, err = strconv.ParseInt(spage, 10, 64); err != nil {
		goto errDeal
	}
	if perpage, err = strconv.ParseInt(sperpage, 10, 64); err != nil {
		goto errDeal
	}

	switch sorder {
	case "0":
		if ml, err = models.GetTrans(query, []string{}, []string{"id"}, []string{"desc"}, page*perpage, perpage); err != nil {
			goto errDeal
		}
	case "1":
		if ml, err = models.GetTrans(query, []string{}, []string{"id"}, []string{"asc"}, page*perpage, perpage); err != nil {
			goto errDeal
		}
	case "2":
		if ml, err = models.GetTrans(query, []string{}, []string{"amount"}, []string{"asc"}, page*perpage, perpage); err != nil {
			goto errDeal
		}
	case "3":
		if ml, err = models.GetTrans(query, []string{}, []string{"amount"}, []string{"desc"}, page*perpage, perpage); err != nil {
			goto errDeal
		}
	}

	conf = GetConfigData()

	for _, v := range ml {
		data.TxId = fmt.Sprintf("%v", v.(models.TransQ).Id)
		data.Amount = fmt.Sprintf("%v", v.(models.TransQ).Amount)
		data.Name = v.(models.TransQ).Name
		data.TxHash = v.(models.TransQ).TxHash
		data.Time = fmt.Sprintf("%v", v.(models.TransQ).Time)
		data.Type = fmt.Sprintf("%v", v.(models.TransQ).Type)
		data.Status = conf.GetMapType()[v.(models.TransQ).Type].Name + TranstateString(v.(models.TransQ).Status)
		datas = append(datas, data)
		count++
	}
	t.HandlerResult(count, datas)
	return
errDeal:
	ErrorHandler(t.Ctx, err)
	return
}

func (t *TransQContoller) HandlerResult(total int, datas []types.TransQData) {
	t.Ctx.Output.JSON(
		types.RspTransQ{
			Total: total,
			Data:  datas,
		}, false, false)
}
