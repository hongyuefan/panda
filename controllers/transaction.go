package controllers

import (
	"fmt"
	"panda/arithmetic"
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

func (t *TransactionContoller) Transactions(ntype int64, uid, pid int64, amount string) (txhash string, err error) {

	var (
		mPlay *models.Player
		conf  models.Config
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
			return
		}
		if result, err = compareAmount(amount, balance); err != nil {
			return
		}
		if result > 0 {
			err = types.Error_Trans_AmountOver
			return
		}
		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, mPlay.Pubkey, amount); err != nil {
			return
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		_, err = t.InsertTransQ(uid, pid, types.Trans_Type_WithDrawal, amount,
			conf.GetMapType()[types.Trans_Type_WithDrawal].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_WithDrawal].Name)
		return

	case types.Trans_Type_Catch: //捕捉
		var (
			balance  string
			mPet     *models.Pet
			result   int
			coldTime float64
			attr     *models.Attrvalue
		)

		if mPet, err = models.GetPetById(pid); err != nil {
			return
		}
		if attr, err = models.GetAttrvalue(pid, types.Attr_Type_Zhili); err != nil {
			return
		}

		if coldTime, err = arithmetic.CatchCold(conf.CatchRation, float64(conf.CatchTimeIntervel), mPet.CatchTimes, float64(mPet.Years), attr.Value); err != nil {
			return
		}

		coldIntervel := int64(coldTime) - (time.Now().Unix() - mPet.LastCatchTime)

		if coldIntervel > 0 {
			return "", fmt.Errorf("距离下次捕捉时间还有 %v 分钟", coldIntervel/60)
		}

		if balance, err = trans.GetBalance(mPlay.PubPublic); err != nil {
			return
		}
		if result, err = compareAmount(conf.GetMapType()[types.Trans_Type_Catch].Amount, balance); err != nil {
			return
		}
		if result > 0 {
			err = types.Error_Trans_AmountOver
			return
		}
		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, conf.OwnerPub, conf.GetMapType()[types.Trans_Type_Catch].Amount); err != nil {
			return
		}

		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		_, err = t.InsertTransQ(uid, pid, types.Trans_Type_Catch, conf.GetMapType()[types.Trans_Type_Catch].Amount,
			conf.GetMapType()[types.Trans_Type_Catch].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Catch].Name)
		return

	case types.Trans_Type_Train: //训练
		var (
			mPet   *models.Pet
			result int
		)
		if mPet, err = models.GetPetById(pid); err != nil {
			return
		}
		if result, err = compareAmount(mPet.TrainTotle, fmt.Sprintf("%v", conf.TrainLimit)); err != nil {
			return
		}
		if result > 0 {
			return "", types.Error_Train_AmountOver
		}
		if txhash, err = trans.DoTransaction(mPlay.PubPrivkey, conf.OwnerPub, amount); err != nil {
			return
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		_, err = t.InsertTransQ(uid, pid, types.Trans_Type_Train, amount,
			conf.GetMapType()[types.Trans_Type_Train].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Train].Name)
		return

	case types.Trans_Type_Bonus: //分红
		if txhash, err = trans.DoTransaction("SB7XB25ZIZB3XCNUXFY64NNIY5WNWPVCMBOLKPPMWWADHKUBM4HZOXVO", mPlay.PubPublic, amount); err != nil {
			return
		}
		//TODO:建立消息组件，保证数据落地存储，防止数据库与区块链数据不一致
		_, err = t.InsertTransQ(uid, pid, types.Trans_Type_Bonus, amount,
			conf.GetMapType()[types.Trans_Type_Bonus].Fee, txhash,
			conf.GetMapType()[types.Trans_Type_Bonus].Name)
		return

	}
	return "", types.Error_Trans_MisType

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
