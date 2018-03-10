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

	switch ntype {

	case types.Trans_Type_WithDrawal: //提现
		var (
			result  int
			balance string
		)
		if mPlay, err = models.GetPlayerById(uid); err != nil {
			return
		}
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
		if mPlay, err = models.GetPlayerById(uid); err != nil {
			return
		}
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

		h, m, s := arithmetic.ParseSecond(coldIntervel)
		if coldIntervel > 0 {
			return "", fmt.Errorf("距离下次捕捉时间还有 %v 小时 %v 分钟 %v 秒", h, m, s)
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
		if mPlay, err = models.GetPlayerById(uid); err != nil {
			return
		}
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
		go func() {
			if err := t.Bonus(conf); err != nil {
				beego.BeeLogger.Error("Bonus Error: %v", err)
			}
		}()
		return
	case types.Trans_Type_Pet: //pet交易
		return
	}
	return "", types.Error_Trans_MisType

}

func (t *TransactionContoller) CountAllIntrest() (totle int64, err error) {

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	query := make(map[string]string)

	for {
		ml, err := models.GetAllPet(query, []string{"Intrest"}, []string{"id"}, []string{"asc"}, offset, limit)
		if err != nil {
			return 0, err
		}
		if len(ml) <= 0 {
			break
		}
		for _, v := range ml {
			totle += v.(map[string]interface{})["Intrest"].(int64)
		}
		offset = +int64(len(ml))
	}
	if totle <= 0 {
		err = fmt.Errorf("CountAllIntrest Totle is Zero")
	}
	return

}

func (t *TransactionContoller) Bonus(conf models.Config) (err error) {
	var (
		offset      int64 = 0
		limit       int64 = 100
		totle       int64
		sbalance    string
		pub_balance float64
		txhash      string
		count       int
	)

	if sbalance, err = trans.GetBalance(conf.OwnerPub); err != nil {
		return
	}
	if pub_balance, err = strconv.ParseFloat(sbalance, 64); err != nil {
		return
	}
	if totle, err = t.CountAllIntrest(); err != nil {
		return
	}

	query := make(map[string]string)

	query["IsBonus"] = "0"

	for {
		ml, err := models.GetAllPet(query, []string{"Id", "Uid", "Intrest"}, []string{"id"}, []string{"asc"}, offset, limit)
		if err != nil {
			continue
		}
		if len(ml) <= 0 {
			break
		}
		for _, v := range ml {

			player, err := models.GetPlayerById(v.(map[string]interface{})["Uid"].(int64))
			if err != nil {
				continue
			}

			bunos := float64(v.(map[string]interface{})["Intrest"].(int64)) / float64(totle) * pub_balance * conf.BonusRatio
			if bunos <= 0 {
				continue
			}
			sBunos := arithmetic.ParseFloat(bunos)

			for {

				if txhash, err = trans.DoTransaction(conf.OwnerPrv, player.PubPublic, sBunos); err == nil {
					break
				}
				count++
				beego.BeeLogger.Error("Bonus Error %v ,Uid:%v,Pid:%v,Amount:%v,Account:%v", err, v.(map[string]interface{})["Uid"].(int64), v.(map[string]interface{})["Id"].(int64), sBunos, player.PubPublic)
				if count >= 3 {
					break
				}
			}

			_, err = t.InsertTransQ(v.(map[string]interface{})["Uid"].(int64), v.(map[string]interface{})["Id"].(int64), types.Trans_Type_Bonus, sBunos,
				conf.GetMapType()[types.Trans_Type_Bonus].Fee, txhash,
				conf.GetMapType()[types.Trans_Type_Bonus].Name)

			models.UpdateBonusProcess(v.(map[string]interface{})["Id"].(int64))
		}
		offset = +int64(len(ml))
	}

	models.SetBonusOver()

	return
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
