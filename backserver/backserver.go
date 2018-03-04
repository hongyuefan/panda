package backserver

import (
	"fmt"
	"panda/arithmetic"
	"panda/models"
	t "panda/transaction"
	"panda/types"
	"time"

	"github.com/astaxie/beego"
)

type BackServer struct {
	com       *models.Common
	chanExit  chan bool
	sleepTime time.Duration
}

func NewBackServer(sleepTime int) *BackServer {
	return &BackServer{
		com:       models.NewCommon(),
		chanExit:  make(chan bool, 0),
		sleepTime: time.Duration(sleepTime) * time.Second,
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

			s.JudgeResult(v.Type, v.TxHash, v.Id, v.UID, v.PID)
		}
	}
	return
}

func (s *BackServer) JudgeResult(itype int64, txhash string, txid, uid, pid int64) {
	switch itype {
	case types.Trans_Type_Catch:
		s.CatchResult(txhash, txid, uid, pid)
	}
}

func (s *BackServer) TrainResult(txhash string, txid, uid, pid int64) {

}

func (s *BackServer) CatchResult(txhash string, txid, uid, pid int64) {
	var (
		err    error
		result int
	)
	attr := &models.Attrvalue{
		Uid: uid,
		Pid: pid,
		Aid: 1,
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
