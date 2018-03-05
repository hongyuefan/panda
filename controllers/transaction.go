package controllers

import (
	"panda/models"
	trans "panda/transaction"
	"panda/types"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type TransactionContoller struct {
	beego.Controller
}

func (t *TransactionContoller) Transactions(ntype int64, uid, pid int64, amount string) (tId int64, err error) {

	var (
		mPlay  *models.Player
		txhash string
		conf   models.Config
	)

	conf = GetConfigData()

	if mPlay, err = models.GetPlayerById(uid); err != nil {
		return
	}

	switch ntype {

	case types.Trans_Type_WithDrawal: //提现
		var (
			result  int
			balance string
		)
		if balance, err = trans.GetBalance(mPlay.PubPublic); err != nil {
			return 0, err
		}
		if result, err = compareAmount(amount, balance); err != nil {
			return 0, err
		}
		if result > 0 {
			return 0, types.Error_Trans_AmountOver
		}
		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, mPlay.Pubkey, amount); err != nil {
			return
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		return t.InsertTransQ(uid, pid, types.Trans_Type_WithDrawal, amount,
			conf.GetMapType()[types.Trans_Type_WithDrawal].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_WithDrawal].Name)

	case types.Trans_Type_Catch: //捕捉
		if time.Now().Unix()-mPlay.LastCatchTime < conf.CatchTimeIntervel {
			return 0, types.Error_Trans_CatchIntervel
		}
		balance, err := trans.GetBalance(mPlay.PubPublic)
		if err != nil {
			return 0, err
		}
		result, err := compareAmount(conf.GetMapType()[types.Trans_Type_Catch].Amount, balance)
		if err != nil {
			return 0, err
		}
		if result > 0 {
			return 0, types.Error_Trans_AmountOver
		}

		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, conf.OwnerPub, conf.GetMapType()[types.Trans_Type_Catch].Amount); err != nil {
			return 0, err
		}

		UpCatchTime(uid)

		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		return t.InsertTransQ(uid, pid, types.Trans_Type_Catch, conf.GetMapType()[types.Trans_Type_Catch].Amount,
			conf.GetMapType()[types.Trans_Type_Catch].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Catch].Name)

	case types.Trans_Type_Train: //训练
		balance, err := trans.GetBalance(mPlay.PubPublic)
		if err != nil {
			return 0, err
		}
		result, err := compareAmount(conf.GetMapType()[types.Trans_Type_Train].Amount, balance)
		if err != nil {
			return 0, err
		}
		if result > 0 {
			return 0, types.Error_Trans_AmountOver
		}
		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, conf.OwnerPub, conf.GetMapType()[types.Trans_Type_Train].Amount); err != nil {
			return 0, err
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		return t.InsertTransQ(uid, pid, types.Trans_Type_Train, conf.GetMapType()[types.Trans_Type_Train].Amount,
			conf.GetMapType()[types.Trans_Type_Train].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Train].Name)

	case types.Trans_Type_Bonus: //分红
		if txhash, err = trans.DoTransaction("SB7XB25ZIZB3XCNUXFY64NNIY5WNWPVCMBOLKPPMWWADHKUBM4HZOXVO", mPlay.PubPublic, amount); err != nil {
			return
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		return t.InsertTransQ(uid, pid, types.Trans_Type_Bonus, amount,
			conf.GetMapType()[types.Trans_Type_Bonus].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Bonus].Name)

	}
	return 0, types.Error_Trans_MisType

}

func (t *TransactionContoller) InsertTransQ(uid, pid, ntype int64, amount, fee, txhash, stype string) (tid int64, err error) {

	transQ := &models.TransQ{
		TxHash: txhash,
		Name:   stype,
		Type:   ntype,
		Status: types.Trans_Status_Waiting,
		UID:    uid,
		PID:    pid,
		Fee:    fee,
		Amount: amount,
		Time:   time.Now().Unix(),
	}
	if tid, err = models.AddTrans(transQ); err != nil {
		beego.BeeLogger.Info("InsertTransQ Failed ,need mamual operation, txhash:%v,type:%v,uid:%v,fee:%v,amount:%v", txhash, ntype, uid, fee, amount)
	}
	return
}

func compareAmount(amount, balance string) (int, error) {

	famount, err := strconv.ParseFloat(amount, 10)
	if err != nil {
		return 0, err
	}

	fbalance, err := strconv.ParseFloat(balance, 10)
	if err != nil {
		return 0, err
	}

	if famount > fbalance {
		return 1, nil
	} else if famount == fbalance {
		return 0, nil
	}
	return -1, nil
}
