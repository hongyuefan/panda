package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	Offer_Doing = 0x02
)

type PetOffer struct {
	Id         int64  `orm:"column(Id);auto"`
	Pid        int64  `orm:"column(pid);"`
	Price      string `orm:"column(price);size(128)"`
	Status     int    `orm:"column(status)"`
	Uid        int64  `orm:"column(uid);"`
	Years      int    `orm:"column(years)"`
	CreateTime int64  `orm:"column(create_time)"`
	UpdateTime int64  `orm:"column(update_time)"`
}

func (t *PetOffer) TableName() string {
	return "pet_offer"
}

func init() {
	orm.RegisterModel(new(PetOffer))
}

func AddOffer(pid int64, uid int64, years int, price string) (id int64, err error) {
	o := orm.NewOrm()
	v := &PetOffer{
		Pid:        pid,
		Uid:        uid,
		Price:      price,
		Years:      years,
		Status:     0,
		CreateTime: time.Now().Unix(),
	}
	return o.Insert(v)
}

func UpdateOfferPrice(id int64, price string) (err error) {
	o := orm.NewOrm()
	v := &PetOffer{
		Id:         id,
		Price:      price,
		UpdateTime: time.Now().Unix(),
	}
	_, err = o.Update(v, "Price", "UpdateTime")
	return
}

func DeleteOffer(id int64) (err error) {
	o := orm.NewOrm()
	v := &PetOffer{
		Id: id,
	}
	if err = o.Read(v); err != nil {
		return
	}
	if v.Status == Offer_Doing {
		return fmt.Errorf("交易进行中，不能撤销")
	}
	_, err = o.Delete(&PetOffer{Id: id})
	return
}

func GetAllOffer(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PetOffer))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PetOffer
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}