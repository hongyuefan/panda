package models

import (
	"github.com/astaxie/beego/orm"
)

type Config struct {
	Id                int64  `orm:"column(id);auto"`
	BaseFee           string `orm:"column(base_fee);size(64)"`
	OwnerPub          string `orm:"column(owner_address);size(256)"`
	CatchTimeIntervel int64  `orm:"column(catch_time_intervel)"`
	TrainLimit        int64  `orm:"column(train_limit)"`
	JudgeTime         int    `orm:"column(judge_time_sleep)"`
	ownerPrv          string
	maptsType         map[int64]*TransType
}

func (t *Config) TableName() string {
	return "config"
}

func (t *Config) SetMapType(m map[int64]*TransType) {
	t.maptsType = m
}

func (t *Config) GetMapType() map[int64]*TransType {
	return t.maptsType
}

func init() {
	orm.RegisterModel(new(Config))
}

func GetConfigById(id int64) (v *Config, err error) {

	o := orm.NewOrm()

	v = &Config{Id: id}

	if err = o.Read(v); err == nil {
		return v, nil
	}

	return nil, err
}
