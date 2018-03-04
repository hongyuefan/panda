package models

import (
	"github.com/astaxie/beego/orm"
)

type TransType struct {
	Id     int64  `orm:"column(id);null"`
	Name   string `orm:"column(name);size(64);null"`
	Amount string `orm:"column(amount);size(128)"`
	Fee    string `orm:"column(fee);size(128)"`
}

func init() {
	orm.RegisterModel(new(TransType))
}

func (t *TransType) TableName() string {
	return "transtype"
}

func GetAllTransType(arryTrans interface{}) (int64, error) {
	com := NewCommon()
	return com.CommonGetAll("transtype", arryTrans)
}
