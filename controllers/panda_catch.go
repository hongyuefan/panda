package controllers

import (
	"fmt"

	"panda/models"
	"panda/types"
	"strconv"

	"github.com/astaxie/beego"
)

type PandaCatchController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *PandaCatchController) HandlerGetPandaCatch() {
	var (
		txId   int64
		stxId  string
		userId int64
		petId  int64
		spetId string
		catch  *models.Catch
		com    *models.Common
		result *types.RspCatchResult
		err    error
	)
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	spetId = c.Ctx.Request.FormValue("petId")
	stxId = c.Ctx.Request.FormValue("txId")

	if petId, err = strconv.ParseInt(spetId, 10, 64); err != nil {
		goto errDeal
	}
	if txId, err = strconv.ParseInt(stxId, 10, 64); err != nil {
		goto errDeal
	}

	catch = &models.Catch{
		Uid:  userId,
		Pid:  petId,
		Txid: txId,
	}

	com = models.NewCommon()
	if err = com.CommonGetOne(catch, "Uid", "Pid", "Txid"); err != nil {
		goto errDeal
	}

	result = &types.RspCatchResult{
		Success: true,
		CTime:   catch.Createtime,
		Result:  catch.Result,
	}

	c.Ctx.Output.JSON(result, false, false)
	return
errDeal:
	result = &types.RspCatchResult{
		Success: false,
		Message: err.Error(),
	}
	c.Ctx.Output.JSON(result, false, false)
	return
}

func (c *PandaCatchController) HandlerPandaCatch() {
	var (
		petId  int64
		spetId string
		err    error
		txhash string
		userId int64
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	spetId = c.Ctx.Request.FormValue("petId")

	if spetId == "" {
		err = fmt.Errorf("petId is null")
		goto errDeal
	}

	if petId, err = strconv.ParseInt(spetId, 10, 64); err != nil {
		goto errDeal
	}

	if _, err = IsExistPanda(userId, petId); err != nil {
		goto errDeal
	}

	if txhash, err = c.trans.Transactions(types.Trans_Type_Catch, userId, petId, ""); err != nil {
		goto errDeal
	}

	c.HandlerResult(true, nil, txhash)

	return
errDeal:
	c.HandlerResult(false, err, "")
	return
}

func (c *PandaCatchController) HandlerResult(success bool, err error, hash string) {

	var msg string

	if err != nil {
		msg = err.Error()
	}
	result := &types.RspCatch{
		Success: success,
		Message: msg,
		Txhash:  hash,
	}
	c.Ctx.Output.JSON(result, false, false)
	return
}

func IsExistPanda(uid, pid int64) (year int, err error) {

	if uid <= 0 || pid <= 0 {
		return 0, fmt.Errorf("request param error")
	}

	pet := &models.Pet{
		Id:  pid,
		Uid: uid,
	}

	orm := models.NewCommon()

	if err = orm.CommonGetOne(pet, "Id", "Uid"); err != nil {
		return 0, err
	}

	return pet.Years, nil
}
