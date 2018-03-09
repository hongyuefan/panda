package models

import (
	"github.com/astaxie/beego/orm"
)

type BonusStatus struct {
	Id       int64  `orm:"column(id);auto" description:"自增"`
	Name     string `orm:"column(name);size(64)" `
	PetId    int64  `orm:"column(pet_id)" description:"用户ID"`
	OverFlag int    `orm:"column(over_flag)"`
}

func (t *BonusStatus) TableName() string {
	return "bonus_status"
}

func init() {
	orm.RegisterModel(new(BonusStatus))
}

func UpdateBonusProcess(petId int64) (err error) {
	o := orm.NewOrm()
	v := &BonusStatus{Id: 1, PetId: petId}
	_, err = o.Update(v, "PetId")
	return
}

func SetBonusOver() (err error) {
	o := orm.NewOrm()
	v := &BonusStatus{Id: 1, OverFlag: 1}
	_, err = o.Update(v, "OverFlag")
	return
}

func ResetBonus() (err error) {
	o := orm.NewOrm()
	v := &BonusStatus{Id: 1, OverFlag: 0, PetId: 0}
	_, err = o.Update(v, "OverFlag", "PetId")
	return
}

func IsBonusOver() bool {
	o := orm.NewOrm()
	v := &BonusStatus{Id: 1}
	if err := o.Read(v); err != nil {
		return false
	}
	if v.OverFlag == 1 {
		return true
	}
	return false
}
