package controllers

import (
	"fmt"
	"panda/models"
	"panda/types"
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
		count                           int
		arryPets                        []types.GetPet
		onePet                          types.GetPet
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
		query["id"] = sPid
	}
	if sstatus != "" {
		query["status"] = sstatus
	}

	if ml, err = models.GetAllPet(query, []string{}, []string{ssort}, []string{sorder}, page*perpage, perpage); err != nil {
		goto errDeal
	}
	for _, v := range ml {
		onePet.Uid = fmt.Sprintf("%v", v.(models.Pet).Uid)
		onePet.Pid = fmt.Sprintf("%v", v.(models.Pet).Id)
		onePet.Cid = fmt.Sprintf("%v", v.(models.Pet).Cid)
		onePet.CreateTime = v.(models.Pet).CreatTime
		onePet.Fid = fmt.Sprintf("%v", v.(models.Pet).Fid)
		onePet.Imag = fmt.Sprintf("%v", v.(models.Pet).SvgPath)
		onePet.PetName = fmt.Sprintf("%v", v.(models.Pet).Petname)
		onePet.Status = v.(models.Pet).Status
		onePet.Years = v.(models.Pet).Years

		arryPets = append(arryPets, onePet)

		count++
	}
	t.HandlerResult(count, arryPets)
	return
errDeal:
	ErrorHandler(t.Ctx, err)
	return
}

func (t *PetController) HandlerResult(total int, datas []types.GetPet) {
	t.Ctx.Output.JSON(
		types.RspGetPets{
			Total: total,
			Pets:  datas,
		}, false, false)
}
