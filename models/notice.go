package models

import (
	"github.com/astaxie/beego/orm"
)

type Notice struct {
	Id    int64  `orm:"column(id);auto"`
	Title string `orm:"column(title)"`
	Text  string `orm:"column(text)"`
	Flag  int    `orm:"column(flag)"`
}

func (t *Notice) TableName() string {
	return "Notice"
}

func init() {
	orm.RegisterModel(new(Notice))
}

func GetNoticPub() (notice string, err error) {
	o := orm.NewOrm()
	v := &Notice{
		Flag: 1,
	}
	if err = o.Read(v, "flag"); err != nil {
		return
	}
	notice = v.Text

	return
}
