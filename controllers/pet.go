package controllers

import (
	"fmt"
	"panda/models"
	"strconv"

	"github.com/astaxie/beego"
)

// PetController operations for Pet
type PetController struct {
	beego.Controller
}

func (t *PetController) HandlerGetPets() {
	var (
		err                             error
		spage, sperpage, ssort, sstatus string
		sPid, sUid, sorder              string
		page, perpage                   int64
		query                           map[string]string
		ml                              []interface{}
	)
	if err = t.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	spage = t.Ctx.Request.FormValue("page")
	sperpage = t.Ctx.Request.FormValue("perpage")
	sorder = t.Ctx.Request.FormValue("order")
	ssort = t.Ctx.Request.FormValue("sort")
	if ssort == "" {
		ssort = "id"
	}
	sUid = t.Ctx.Request.FormValue("memberId")
	sPid = t.Ctx.Request.FormValue("petId")
	sstatus = t.Ctx.Request.FormValue("status")

	if page, err = strconv.ParseInt(spage, 10, 64); err != nil {
		goto errDeal
	}
	if perpage, err = strconv.ParseInt(sperpage, 10, 64); err != nil {
		goto errDeal
	}

	query = make(map[string]string, 0)

	if sUid != "" {
		query["uid"] = sUid
	}
	if sPid != "" {
		query["pid"] = sPid
	}
	if sstatus != "" {
		query["status"] = sstatus
	}

	if ml, err = models.GetAllPet(query, []string{}, []string{ssort}, []string{sorder}, page*perpage, perpage); err != nil {
		goto errDeal
	}
	fmt.Println(ml)
	return
errDeal:
	return
}
