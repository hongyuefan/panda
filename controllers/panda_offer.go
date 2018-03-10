package controllers

import (
	"fmt"
	"panda/models"
	"panda/types"
	"strconv"

	"github.com/astaxie/beego"
)

type OfferController struct {
	beego.Controller
	trans *TransactionContoller
}

func (c *OfferController) HandlerDeleteOffer() {
	var (
		err             error
		sOfferId        string
		offerId, userId int64
		rspUpdate       types.RspUpdateOffer
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	sOfferId = c.Ctx.Request.FormValue("offerId")

	if offerId, err = strconv.ParseInt(sOfferId, 10, 64); err != nil {
		goto errDeal
	}
	if err = models.DeleteOffer(offerId, userId); err != nil {
		goto errDeal
	}
	rspUpdate = types.RspUpdateOffer{
		Success: true,
	}
	c.Ctx.Output.JSON(rspUpdate, false, false)
	return
errDeal:
	rspUpdate = types.RspUpdateOffer{
		Success: false,
		Message: err.Error(),
	}
	c.Ctx.Output.JSON(rspUpdate, false, false)
	return
}

func (c *OfferController) HandlerUpdatePrice() {
	var (
		sOfferId        string
		sPrice          string
		err             error
		userId, offerId int64
		rspUpdate       types.RspUpdateOffer
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	if userId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	sOfferId = c.Ctx.Request.FormValue("offerId")
	sPrice = c.Ctx.Request.FormValue("price")

	if offerId, err = strconv.ParseInt(sOfferId, 10, 64); err != nil {
		goto errDeal
	}
	if _, err = strconv.ParseFloat(sPrice, 64); err != nil {
		goto errDeal
	}

	if err = models.UpdateOfferPrice(offerId, userId, sPrice); err != nil {
		goto errDeal
	}

	rspUpdate = types.RspUpdateOffer{
		Success: true,
	}
	c.Ctx.Output.JSON(rspUpdate, false, false)
	return
errDeal:
	rspUpdate = types.RspUpdateOffer{
		Success: false,
		Message: err.Error(),
	}
	c.Ctx.Output.JSON(rspUpdate, false, false)
	return
}

func (c *OfferController) HandlerGetOffer() {
	var (
		sPetId, sUid, sStatus          string
		spage, sperpage, sorder, ssort string
		err                            error
		page, perpage                  int64
		query                          map[string]string
		arryOffers                     []types.GetOffers
		offer                          types.GetOffers
		ml                             []interface{}
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	sPetId = c.Ctx.Request.FormValue("petId")
	sUid = c.Ctx.Request.FormValue("memberId")
	sStatus = c.Ctx.Request.FormValue("status")
	spage = c.Ctx.Request.FormValue("page")
	sperpage = c.Ctx.Request.FormValue("perpage")
	sorder = c.Ctx.Request.FormValue("order")
	ssort = c.Ctx.Request.FormValue("sort")
	if ssort == "" {
		ssort = "id"
	}
	query = make(map[string]string, 0)

	if sPetId != "" {
		query["pid"] = sPetId
	}
	if sUid != "" {
		query["uid"] = sUid
	}
	if sStatus != "" {
		query["status"] = sStatus
	}

	if page, err = strconv.ParseInt(spage, 10, 64); err != nil {
		goto errDeal
	}
	if perpage, err = strconv.ParseInt(sperpage, 10, 64); err != nil {
		goto errDeal
	}

	if ml, err = models.GetAllOffer(query, []string{}, []string{ssort}, []string{sorder}, page*perpage, perpage); err != nil {
		goto errDeal
	}

	for _, v := range ml {
		offer.Id = v.(models.PetOffer).Id
		offer.CreateTime = v.(models.PetOffer).CreateTime
		offer.Pid = fmt.Sprintf("%v", v.(models.PetOffer).Pid)
		offer.Price = v.(models.PetOffer).Price
		offer.Status = v.(models.PetOffer).Status
		offer.Uid = fmt.Sprintf("%v", v.(models.PetOffer).Uid)
		offer.UpdateTime = v.(models.PetOffer).UpdateTime
		offer.Years = v.(models.PetOffer).Years

		arryOffers = append(arryOffers, offer)
	}

	c.HandlerGetOfferResult(len(ml), arryOffers)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *OfferController) HandlerGetOfferResult(total int, offers []types.GetOffers) {
	respOffers := &types.RspGetOffers{
		Total:  total,
		Offers: offers,
	}
	c.Ctx.Output.JSON(respOffers, false, false)
}

func (c *OfferController) HandlerBuyPet() {
	var (
		sofferId, txhash string
		buyerId, offerId int64
		err              error
		ok               bool
		rspTransPet      types.RspTransPet
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}
	if buyerId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	sofferId = c.Ctx.Request.FormValue("offerId")

	if offerId, err = strconv.ParseInt(sofferId, 10, 64); err != nil {
		goto errDeal
	}

	if ok, err = models.IsOfferNormal(offerId); err != nil {
		goto errDeal
	}
	if !ok {
		err = fmt.Errorf("宠物交易处理中，暂时不能买卖")
		goto errDeal
	}
	if txhash, err = c.trans.Transactions(types.Trans_Type_Offer, buyerId, 0, offerId, "0"); err != nil {
		goto errDeal
	}
	rspTransPet = types.RspTransPet{
		Success: true,
		TxHash:  txhash,
	}
	c.Ctx.Output.JSON(rspTransPet, false, false)
	return
errDeal:
	rspTransPet = types.RspTransPet{
		Success: false,
		Message: err.Error(),
	}
	c.Ctx.Output.JSON(rspTransPet, false, false)
	return
}

func (c *OfferController) HandlerDoOffer() {
	var (
		petId, userId, offerId int64
		spetId, sprice         string
		years                  int
		err                    error
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
	sprice = c.Ctx.Request.FormValue("price")
	if sprice == "" {
		err = fmt.Errorf("sprice is null")
		goto errDeal
	}

	if petId, err = strconv.ParseInt(spetId, 10, 64); err != nil {
		goto errDeal
	}

	if years, err = IsExistPanda(userId, petId); err != nil {
		goto errDeal
	}

	if offerId, err = models.AddOffer(petId, userId, years, sprice); err != nil {
		goto errDeal
	}

	c.HandlerDoOfferResult(true, nil, offerId)

	return
errDeal:
	c.HandlerDoOfferResult(false, err, 0)
	return
}

func (c *OfferController) HandlerDoOfferResult(success bool, err error, id int64) {

	var msg string

	if err != nil {
		msg = err.Error()
	}
	result := &types.RspAddOffer{
		Success: success,
		Message: msg,
		OfferId: fmt.Sprintf("%v", id),
	}
	c.Ctx.Output.JSON(result, false, false)
	return
}
