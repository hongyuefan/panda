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
		err                     error
		data                    types.TransQData
		datas                   []types.TransQData
		userId                  int64
		spage, sorder, stype    string
		ssort, sperpage, txhash string
		status, spetId          string
		page, perpage           int64
		query                   map[string]string
		ml                      []interface{}
		conf                    models.Config
		transQ                  types.RspTransQ
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
	ssort = t.Ctx.Request.FormValue("sort")
	status = t.Ctx.Request.FormValue("status")
	stype = t.Ctx.Request.FormValue("type")
	txhash = t.Ctx.Request.FormValue("txhash")
	spetId = t.Ctx.Request.FormValue("petId")

	if spetId != "" {
		query["pid"] = spetId
	}
	if stype != "" {
		query["type"] = stype
	}
	if txhash != "" {
		query["txhash"] = txhash
	}
	if ssort == "" {
		ssort = "id"
	}
	if status != "" {
		query["status"] = status
	}

	if page, err = strconv.ParseInt(spage, 10, 64); err != nil {
		goto errDeal
	}
	if perpage, err = strconv.ParseInt(sperpage, 10, 64); err != nil {
		goto errDeal
	}

	if ml, err = models.GetTrans(query, []string{}, []string{ssort}, []string{sorder}, page*perpage, perpage); err != nil {
		goto errDeal
	}

	conf = GetConfigData()

	for _, v := range ml {
		data.TxId = fmt.Sprintf("%v", v.(models.TransQ).Id)
		data.Amount = fmt.Sprintf("%v", v.(models.TransQ).Amount)
		data.Name = v.(models.TransQ).Name
		data.TxHash = v.(models.TransQ).TxHash
		data.Time = fmt.Sprintf("%v", v.(models.TransQ).Time)
		data.Type = fmt.Sprintf("%v", v.(models.TransQ).Type)
		data.PetId = fmt.Sprintf("%v", v.(models.TransQ).PID)
		data.Status = conf.GetMapType()[v.(models.TransQ).Type].Name + TranstateString(v.(models.TransQ).Status)
		datas = append(datas, data)
	}
	transQ = types.RspTransQ{
		Total: len(datas),
		Data:  datas,
	}
	SuccessHandler(t.Ctx, transQ)
	return
errDeal:
	ErrorHandler(t.Ctx, err)
	return
}
