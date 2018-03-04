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

func SetConfigData(conf *models.Config, mtx map[int64]*models.TransType) {
	lock.Lock()
	defer lock.Unlock()

	configData.Id = conf.Id
	configData.BaseFee = conf.BaseFee
	configData.OwnerPub = conf.OwnerPub
	configData.JudgeTime = conf.JudgeTime
	configData.CatchTimeIntervel = conf.CatchTimeIntervel
	configData.TrainLimit = conf.TrainLimit
	configData.SetMapType(mtx)
}

type ConfigDataController struct {
	beego.Controller
}

func (c *ConfigDataController) LoadConfig() {
	var (
		id     int64
		sid    string
		err    error
		conf   *models.Config
		arryTx []*models.TransType
		mTx    map[int64]*models.TransType
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

	SetConfigData(conf, mTx)

	if backServer != nil {
		backServer.OnStop()
		backServer = nil
	}

	backServer = backserver.NewBackServer(conf.JudgeTime)
	backServer.OnStart()

	c.Ctx.Output.JSON(configData, false, false)

	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}
