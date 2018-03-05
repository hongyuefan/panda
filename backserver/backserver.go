package backserver

import (
	"fmt"
	"panda/arithmetic"
	"panda/models"
	t "panda/transaction"
	"panda/types"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type BackServer struct {
	conf      *models.Config
	com       *models.Common
	chanExit  chan bool
	sleepTime time.Duration
}

func NewBackServer(conf *models.Config) *BackServer {
	return &BackServer{
		com:       models.NewCommon(),
		chanExit:  make(chan bool, 0),
		conf:      conf,
		sleepTime: time.Duration(conf.JudgeTime) * time.Second,
	}
}

func (s *BackServer) OnStart() {
	go s.HandlerTransaction()
}

func (s *BackServer) OnStop() {
	close(s.chanExit)
}

func (s *BackServer) HandlerTransaction() {
	ticker := time.NewTicker(s.sleepTime)
	for {
		select {
		case <-ticker.C:
			s.JudgeTransaction()
		case <-s.chanExit:
			ticker.Stop()
			return
		}
	}
}

func (s *BackServer) JudgeTransaction() (err error) {
	var (
		mFilter    map[string]interface{}
		arryTransQ []*models.TransQ
		count      int64
		result     int
	)
	mFilter = make(map[string]interface{}, 0)

	mFilter["Status"] = types.Trans_Status_Waiting

	if count, err = s.com.CommonGetFilterAll("transaction", mFilter, &arryTransQ); err != nil {
		beego.BeeLogger.Error("JudgeTransaction getTransQ error:%v", err)
		return
	}
	beego.BeeLogger.Debug("JudgeTransaction Count:%v", count)

	for _, v := range arryTransQ {
		result, err = t.QueryTransaction(v.TxHash)
		if result == t.Trans_Retry {
			beego.BeeLogger.Error("JudgeTransaction QueryTransaction Need Retry Txhash:%v,Error:%v", v.TxHash, err)
		} else if result == t.Trans_Failed {
			beego.BeeLogger.Debug("JudgeTransaction QueryTransaction Failed Txhash:%v,Error:%v", v.TxHash, err)
			v.Status = types.Trans_Status_Failed
			models.UpdateTransById(v, "Status")
		} else if result == t.Trans_Success {
			beego.BeeLogger.Debug("JudgeTransaction QueryTransaction Success Txhash:%v", v.TxHash)
			v.Status = types.Trans_Status_Success
			models.UpdateTransById(v, "Status")

			s.JudgeResult(v.Type, v.TxHash, v.Amount, v.Id, v.UID, v.PID)
		}
	}
	return
}

func (s *BackServer) JudgeResult(itype int64, txhash, amount string, txid, uid, pid int64) {
	switch itype {
	case types.Trans_Type_Catch:
		s.CatchResult(txhash, txid, uid, pid)
	case types.Trans_Type_Train:
		s.TrainResult(txhash, amount, txid, uid, pid)
	}
}

func (s *BackServer) TrainResult(txhash, amount string, txid, uid, pid int64) {
	var (
		err error
	)
	st1, err := arithmetic.TrainResult(txhash, 1)
	if err != nil {
		beego.BeeLogger.Error("TrainResult error %v,txhash %v N %v", err, txhash, 1)
		return
	}
	st2, err := arithmetic.TrainResult(txhash, 2)
	if err != nil {
		beego.BeeLogger.Error("TrainResult error %v,txhash %v N %v", err, txhash, 2)
		return
	}
	st3, err := arithmetic.TrainResult(txhash, 3)
	if err != nil {
		beego.BeeLogger.Error("TrainResult error %v,txhash %v N %v", err, txhash, 3)
		return
	}

	mjAttr, err := s.getAttr(uid, pid, types.Attr_Type_Minjie)
	if err != nil {
		return
	}
	llAttr, err := s.getAttr(uid, pid, types.Attr_Type_Liliang)
	if err != nil {
		return
	}
	zlAttr, err := s.getAttr(uid, pid, types.Attr_Type_Zhili)
	if err != nil {
		return
	}

	mjAttr.Value, err = s.trainArithmetic(types.Attr_Type_Minjie, st1, mjAttr.Value, amount)
	if err != nil {
		return
	}
	llAttr.Value, err = s.trainArithmetic(types.Attr_Type_Liliang, st2, llAttr.Value, amount)
	if err != nil {
		return
	}
	zlAttr.Value, err = s.trainArithmetic(types.Attr_Type_Zhili, st3, zlAttr.Value, amount)
	if err != nil {
		return
	}
	s.upDateAttrValue(mjAttr)
	s.upDateAttrValue(llAttr)
	s.upDateAttrValue(zlAttr)
	return
}

func (s *BackServer) upDateAttrValue(attr *models.Attrvalue) (err error) {
	if _, err = s.com.CommonUpdateById(attr, "value"); err != nil {
		beego.BeeLogger.Error("upDateAttrValue error:%v, attrId:%v value:%v", err, attr.Id, attr.Value)
	}
	return
}

func (s *BackServer) getAttr(uid, pid, aid int64) (attr *models.Attrvalue, err error) {
	attr = &models.Attrvalue{
		Uid: uid,
		Pid: pid,
		Aid: aid,
	}
	if err = s.com.CommonGetOne(attr, "Uid", "Pid", "Aid"); err != nil {
		beego.BeeLogger.Error(" getAttr Error %v ,Uid:%v,Pid:%v,Aid:%v", err, uid, pid, aid)
		return
	}
	return
}

func (s *BackServer) trainArithmetic(ntype int64, state int, balance, amount string) (result string, err error) {
	var (
		strMulti string = "1"
		nMulti   float64
		nBalance float64
		nResult  float64
		nAmount  float64
	)

	if state == 1 {
		strMulti = s.conf.GetMapAttr()[ntype].Normal
	} else if state == 2 {
		strMulti = s.conf.GetMapAttr()[ntype].Special
	}

	if nMulti, err = strconv.ParseFloat(strMulti, 64); err != nil {
		beego.BeeLogger.Error("trainArithmetic parse Multi error %v", err)
		return
	}
	if nBalance, err = strconv.ParseFloat(balance, 64); err != nil {
		beego.BeeLogger.Error("trainArithmetic parse Balance error %v", err)
		return
	}
	if nAmount, err = strconv.ParseFloat(amount, 64); err != nil {
		beego.BeeLogger.Error("trainArithmetic pars Amount error %v", err)
		return
	}

	nResult = nBalance + nBalance*nMulti*nAmount/10.0

	result = fmt.Sprintf("%v", nResult)

	return
}

func (s *BackServer) CatchResult(txhash string, txid, uid, pid int64) {
	var (
		err    error
		result int
	)
	attr := &models.Attrvalue{
		Uid: uid,
		Pid: pid,
		Aid: types.Attr_Type_Minjie,
	}

	if err = s.com.CommonGetOne(attr, "Uid", "Pid", "Aid"); err != nil {
		beego.BeeLogger.Error("CatchResult GetAttrValue Error %v ,Uid:%v,Pid:%v", err, uid, pid)
		return
	}
	if result, err = arithmetic.CatchResult(txhash, attr.Value); err != nil {
		beego.BeeLogger.Error("CatchResult CatchResult Error %v ,txhash:%v,value:%v", err, txhash, attr.Value)
		return
	}

	catch := &models.Catch{
		Uid:        uid,
		Createtime: time.Now().Unix(),
		Txid:       txid,
		Result:     result,
		Pid:        pid,
	}

	cid, err := models.AddCatch(catch)
	if err != nil {
		beego.BeeLogger.Error("CatchResult AddCatch Error %v ,uid,txid,result,pid:%v,%v,%v,%v", err, uid, txid, result, pid)
		return
	}
	if result != 0 {
		//TODO:生成svg文件，返回路径
		fmt.Println(result)

		pet := &models.Pet{
			Uid:     uid,
			Cid:     cid,
			Fid:     pid,
			Petname: txhash,
			Years:   attr.Years + 1,
			SvgPath: "svgpath",
			Status:  0,
		}

		if _, err = models.AddPet(pet); err != nil {
			beego.BeeLogger.Info("CatchResult AddPet Error %v need operation manual:Uid,Cid,Fid,Petnam,years,svgpath",
				err, uid, cid, pid, txhash, attr.Years+1, "svgpath")
			return
		}
	}
	return
}
