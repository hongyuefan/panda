package models

import (
	"github.com/astaxie/beego/orm"
)

type Config struct {
	Id                int64   `orm:"column(id);auto"`
	BaseFee           string  `orm:"column(base_fee);size(64)"`
	OwnerPub          string  `orm:"column(owner_address);size(256)"`
	CatchTimeIntervel int64   `orm:"column(catch_time_intervel)"`
	TrainLimit        int64   `orm:"column(train_limit)"`
	JudgeTime         int     `orm:"column(judge_time_sleep)"`
	CatchRation       float64 `orm:"column(catch_ration)"`
	RareAttribute     float64 `orm:"column(rare_attribute)"`
	HostUrl           string  `orm:"column(host_url)"`
	BonusRatio        float64 `orm:"column(bonus_ratio)"`
	OwnerPrv          string  `orm:"column(owner_priv)"`
	InvitationLimit   int     `orm:"column(invitation_limit)"`
	IsInvitation      int     `orm:"column(is_invitation)"`
	InvitationYears   int     `orm:"column(invitation_years)"`
	maptsType         map[int64]*TransType
	mapattrType       map[int64]*Attribute
}

func (t *Config) TableName() string {
	return "config"
}

func (t *Config) SetMapAttr(m map[int64]*Attribute) {
	t.mapattrType = m
}

func (t *Config) GetMapAttr() map[int64]*Attribute {
	return t.mapattrType
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
