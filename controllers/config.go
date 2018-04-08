package controllers

import (
	"panda/backserver"
	"panda/models"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
)

var (
	configData models.Config
	lock       sync.RWMutex
	backServer *backserver.BackServer
)

func GetConfigData() models.Config {
	lock.RLock()
	defer lock.RUnlock()

	return configData
}

func SetConfigData(conf *models.Config, mtx map[int64]*models.TransType, mat map[int64]*models.Attribute) {
	lock.Lock()
	defer lock.Unlock()

	configData.Id = conf.Id
	configData.BaseFee = conf.BaseFee
	configData.OwnerPub = conf.OwnerPub
	configData.OwnerPrv = "SDETYYST6UIKAC3MCDU33677AUIWWVCR4MSNKKDFGFJYAJYWALLTTGQ2"
	configData.JudgeTime = conf.JudgeTime
	configData.CatchTimeIntervel = conf.CatchTimeIntervel
	configData.TrainLimit = conf.TrainLimit
	configData.CatchRation = conf.CatchRation
	configData.RareAttribute = conf.RareAttribute
	configData.HostUrl = conf.HostUrl
	configData.BonusRatio = conf.BonusRatio
	configData.IsInvitation = conf.IsInvitation
	configData.InvitationLimit = conf.InvitationLimit
	configData.InvitationYears = conf.InvitationYears
	configData.AppId = conf.AppId
	configData.AppKey = conf.AppKey
	configData.TplId = conf.TplId
	configData.SetMapType(mtx)
	configData.SetMapAttr(mat)
}

type ConfigDataController struct {
	beego.Controller
}

func (c *ConfigDataController) LoadConfig() {
	var (
		id       int64
		sid      string
		err      error
		conf     *models.Config
		arryTx   []*models.TransType
		arryAttr []*models.Attribute
		mTx      map[int64]*models.TransType
		mAt      map[int64]*models.Attribute
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	sid = c.Ctx.Request.FormValue("id")

	if id, err = strconv.ParseInt(sid, 10, 64); err != nil {
		goto errDeal
	}

	if conf, err = models.GetConfigById(id); err != nil {
		goto errDeal
	}

	if _, err = models.GetAllTransType(&arryTx); err != nil {
		goto errDeal
	}

	mTx = make(map[int64]*models.TransType, 0)

	for _, v := range arryTx {
		mTx[v.Id] = v
	}

	if _, err = models.GetAllAttribute(&arryAttr); err != nil {
		goto errDeal
	}

	mAt = make(map[int64]*models.Attribute, 0)

	for _, v := range arryAttr {
		mAt[v.Id] = v
	}

	conf.HostUrl += c.Ctx.Request.Host

	SetConfigData(conf, mTx, mAt)

	if backServer != nil {
		backServer.OnStop()
		backServer = nil
	}

	backServer = backserver.NewBackServer(GetConfigData())
	backServer.OnStart()

	c.Ctx.Output.JSON(configData, false, false)

	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
