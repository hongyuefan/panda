package backserver

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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

func NewBackServer(conf models.Config) *BackServer {
	return &BackServer{
		com:       models.NewCommon(),
		chanExit:  make(chan bool, 0),
		conf:      &conf,
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
		}
		s.JudgeResult(result, v.Type, v.TxHash, v.Amount, v.Id, v.UID, v.PID)
	}
	return
}

func (s *BackServer) JudgeResult(result int, itype int64, txhash, amount string, txid, uid, pid int64) {
	switch itype {
	case types.Trans_Type_Catch:
		if result == t.Trans_Success {
			if err := models.UpCatchTime(pid); err != nil {
				beego.BeeLogger.Error("UpCatchTime error %v,txhash %v Uid %v ", err, txhash, uid)
			}
			s.CatchResult(txhash, txid, uid, pid)
		}

	case types.Trans_Type_Train:
		if result == t.Trans_Success {
			s.TrainResult(txhash, amount, txid, uid, pid)
		}
	case types.Trans_Type_Bonus:
		if result == t.Trans_Success {
			if err := models.BonusOver(pid, 1); err != nil {
				beego.BeeLogger.Error("BonusOver error %v, Pid %v, Txhash %v", err, pid, txhash)
			}
		} else if result == t.Trans_Failed {
			if err := models.BonusOver(pid, 2); err != nil {
				beego.BeeLogger.Error("BonusOver error %v, Pid %v, Txhash %v", err, pid, txhash)
			}
		}
		if models.IsBonusOver() {
			models.BonusReset()
		}
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

	mjAttr.Value, mjAttr.Multi, err = s.trainArithmetic(types.Attr_Type_Minjie, st1, mjAttr.Value, amount)
	if err != nil {
		return
	}
	llAttr.Value, llAttr.Multi, err = s.trainArithmetic(types.Attr_Type_Liliang, st2, llAttr.Value, amount)
	if err != nil {
		return
	}
	zlAttr.Value, zlAttr.Multi, err = s.trainArithmetic(types.Attr_Type_Zhili, st3, zlAttr.Value, amount)
	if err != nil {
		return
	}
	s.upDateAttrValue(mjAttr)
	s.upDateAttrValue(llAttr)
	s.upDateAttrValue(zlAttr)

	if err := models.UpTrainTotleAndIntrest(pid, amount, zlAttr.Value, llAttr.Value, s.conf.RareAttribute); err != nil {
		beego.BeeLogger.Error("UpTrainTotle error %v,txhash %v Uid %v Amount %v", err, txhash, uid, amount)
	}
	return
}

func (s *BackServer) upDateAttrValue(attr *models.Attrvalue) (err error) {
	if _, err = s.com.CommonUpdateById(attr, "value", "multi"); err != nil {
		beego.BeeLogger.Error("upDateAttrValue error:%v, attrId:%v value:%v multi:%v", err, attr.Id, attr.Value, attr.Multi)
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

func (s *BackServer) trainArithmetic(ntype int64, state int, balance, amount string) (result, strMulti string, err error) {
	var (
		nMulti   float64
		nBalance float64
		nResult  float64
		nAmount  float64
	)
	strMulti = "1"

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
		err                       error
		result                    int
		mjValue, llValue, zlValue float64
		petId                     int64
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
		svgPath := s.Generate_svg(result, beego.AppConfig.String("svg_path"), fmt.Sprintf("%v", pid))

		pet := &models.Pet{
			Uid:           uid,
			Cid:           cid,
			Fid:           pid,
			Petname:       txhash,
			Years:         attr.Years + 1,
			SvgPath:       s.conf.HostUrl + types.Svg_File_Path + "/" + svgPath,
			Status:        0,
			TrainTotle:    "0",
			LastCatchTime: time.Now().Unix(),
			CreatTime:     time.Now().Unix(),
			CatchTimes:    0,
			IsRare:        result,
			IsBonus:       0,
		}

		if petId, err = models.AddPet(pet); err != nil {
			beego.BeeLogger.Info("CatchResult AddPet Error %v need operation manual:Uid %v,Cid %v,Fid %v,Petnam %v,years %v,svgpath %v",
				err, uid, cid, pid, txhash, attr.Years+1, "svgpath")
			return
		}

		if mjValue, err = s.RandValue(s.conf.GetMapAttr()[types.Attr_Type_Minjie].Limit); err != nil {
			beego.BeeLogger.Error("CatchResult RandValue Error %v", err)
			return
		}

		if _, err = models.AddAttrvalue(models.NewAttrvalue(petId, uid, types.Attr_Type_Minjie, attr.Years+1, fmt.Sprintf("%v", mjValue))); err != nil {
			beego.BeeLogger.Info("CatchResult AddAttrvalue Error %v need operation manual:Aid %v Uid %v,Cid %v,Fid %v,Petnam %v,years %v",
				err, types.Attr_Type_Minjie, uid, cid, petId, txhash, attr.Years+1)
			return
		}
		if llValue, err = s.RandValue(s.conf.GetMapAttr()[types.Attr_Type_Liliang].Limit); err != nil {
			beego.BeeLogger.Error("CatchResult RandValue Error %v", err)
			return
		}
		if _, err = models.AddAttrvalue(models.NewAttrvalue(petId, uid, types.Attr_Type_Liliang, attr.Years+1, fmt.Sprintf("%v", llValue))); err != nil {
			beego.BeeLogger.Info("CatchResult AddAttrvalue Error %v need operation manual:Aid %v Uid %v,Cid %v,Fid %v,Petnam %v,years %v",
				err, types.Attr_Type_Liliang, uid, cid, petId, txhash, attr.Years+1)
			return
		}
		if zlValue, err = s.RandValue(s.conf.GetMapAttr()[types.Attr_Type_Zhili].Limit); err != nil {
			beego.BeeLogger.Error("CatchResult RandValue Error %v", err)
			return
		}
		if _, err = models.AddAttrvalue(models.NewAttrvalue(petId, uid, types.Attr_Type_Zhili, attr.Years+1, fmt.Sprintf("%v", zlValue))); err != nil {
			beego.BeeLogger.Info("CatchResult AddAttrvalue Error %v need operation manual:Aid %v Uid %v,Cid %v,Fid %v,Petnam %v,years %v",
				err, types.Attr_Type_Zhili, uid, cid, pid, txhash, attr.Years+1)
			return
		}

	}
	return
}

func (s *BackServer) RandValue(attrLimit string) (value float64, err error) {

	min, max, err := arithmetic.ParseLimit(attrLimit)
	if err != nil {
		return 0, err
	}
	value = arithmetic.GetRand(min, max)
	return
}

func (s *BackServer) Generate_svg(flag int, basePath string, petID string) (svgPath string) {

	//svg head
	svg := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 800\">"

	//random for panda base color
	color := generate_rand(7)

	//memset
	var selectflag_times = make([]int, 4, 10)
	var selectflag_check = make([]int, 4, 10)

	for i, _ := range selectflag_check {
		if i > 1 {
			count, _ := models.GetCountBySelectId(int64(i))
			selectflag_check[i] = int(generate_rand(count))
		}
	}

	//query structController
	query := make(map[string]string, 0)

	//Get svg_catatory order by rank
	ml, err := models.GetAllSvgcata(query, []string{}, []string{"rank"}, []string{"asc"}, 0, 20)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range ml {
		catagory_id := v.(models.Svg_catagory).Id
		select_flag := v.(models.Svg_catagory).Select_flag
		bodycolor_flag := v.(models.Svg_catagory).Bodycolor_flag
		percent := v.(models.Svg_catagory).Probability

		switch select_flag {
		case 0:
			//use random color
			if bodycolor_flag == 1 {

				rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
				if rand == -1 {
					break
				}

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["base_color"] = fmt.Sprintf("%v", color)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			} else {
				// bodyline
				count := models.GetCountByCatagoryId(catagory_id)

				if count == 0 {
					break
				}
				rand := generate_rand(count)

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			}

		case 1:
			//use random element
			//Be or not be
			rand := generate_rand(int64(100 / percent))

			if (rand != 0 && percent != 1) || (flag == 1 && percent == 1) {
				//get count of this item
				count := models.GetCountByCatagoryId(catagory_id)
				if count == 0 {
					break
				}
				rand := generate_rand(count)

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			}

		default:
			//choose one; drop another
			//How many items where select_flag = this, calculator only first time.
			if selectflag_times[select_flag] == selectflag_check[select_flag] {
				if bodycolor_flag == 1 {
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
					if rand == -1 {
						break
					}

					query := make(map[string]string, 0)
					query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
					query["base_color"] = fmt.Sprintf("%v", color)
					query["p_id"] = fmt.Sprintf("%v", 0)
					resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

					for _, v := range resultl {
						svg += v.(models.Svg_info).Svg_dtl
						if v.(models.Svg_info).N_id != 0 {
							// link the next svg to be strcat
							query := make(map[string]string, 0)
							query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
							query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
							models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
							for _, v := range resultl {
								svg += v.(models.Svg_info).Svg_dtl
							}
						}
					}
				} else {
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id))
					if rand == -1 {
						break
					}

					query := make(map[string]string, 0)
					query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
					query["p_id"] = fmt.Sprintf("%v", 0)
					resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)
					for _, v := range resultl {
						svg += v.(models.Svg_info).Svg_dtl
						// link the next svg to be strcat
						if v.(models.Svg_info).N_id != 0 {
							query := make(map[string]string, 0)
							query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
							query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
							models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
							for _, v := range resultl {
								svg += v.(models.Svg_info).Svg_dtl
							}
						}
					}
				}
			} else {
				//do nothing
			}
			selectflag_times[select_flag]++
		}
	}

	svg += "</svg>"

	// Create svg file and write inside
	//	if basePath[0:1] != "\\" {
	//		basePath += "\\"
	//	}
	fileName := fmt.Sprintf("%s%v.svg", time.Now().Format("20060102150405"), petID)
	strFile := basePath + fileName
	f, err := os.Create(strFile)
	w := bufio.NewWriter(f)
	if _, err = w.WriteString(svg); err != nil {
		//err deal
		fmt.Println(err)
	}
	w.Flush()

	f.Close()

	return fileName
}

func getSvgDetail(catagory_id int64, color_flag int, color int64, index int64) (svg_detail string) {
	query := make(map[string]string, 0)
	if color_flag == 1 {
		query["base_color"] = fmt.Sprintf("%v", color)
	}
	query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
	query["p_id"] = fmt.Sprintf("%v", 0)
	resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, index, 1)

	for _, v := range resultl {
		svg_detail = v.(models.Svg_info).Svg_dtl
		// link the next svg to be strcat
		if v.(models.Svg_info).N_id != 0 {
			query := make(map[string]string, 0)
			query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
			query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
			resultlin, _ := models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
			for _, vin := range resultlin {
				svg_detail += vin.(models.Svg_info).Svg_dtl
			}
		}
	}
	return svg_detail
}

/*get random number*/
func generate_rand(number int64) (nRand int64) {
	if number <= 0 {
		return -1
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	return int64(seed.Intn(int(number)))
}
