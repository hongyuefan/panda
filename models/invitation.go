package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Invitation struct {
	Id    int64  `orm:column(id)`
	Uid   int64  `orm:"column(uid)"`
	Code  string `orm:"column(code)"`
	Count int    `orm:"column(count)"`
	Flag  int    `orm:"column(flag)"`
}

func (t *Invitation) TableName() string {
	return "invitation"
}

func init() {
	orm.RegisterModel(new(Invitation))
}

func SetInvitationFlag(id int64, flag int) (err error) {
	o := orm.NewOrm()
	v := &Invitation{
		Id:   id,
		Flag: flag,
	}
	_, err = o.Update(v, "flag")
	return
}

func AddInvitation(uid int64, code string) (err error) {
	o := orm.NewOrm()

	inv := &Invitation{
		Uid:   uid,
		Code:  code,
		Count: 0,
		Flag:  0,
	}
	_, err = o.Insert(inv)
	return
}

func UpdateInvitationCount(code string) (err error) {
	o := orm.NewOrm()
	inv := &Invitation{
		Code: code,
	}
	if err = o.Read(inv, "code"); err != nil {
		return
	}
	inv.Count += 1
	_, err = o.Update(inv, "count")
	return
}

func GetInvitationByUid(uid int64) (code string, count, flag int, err error) {
	o := orm.NewOrm()
	v := &Invitation{
		Uid: uid,
	}
	if err = o.Read(v, "uid"); err != nil {
		return
	}
	code = v.Code
	count = v.Count
	flag = v.Flag

	return
}

func GetAllInvitations(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Invitation))
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

	var l []Invitation
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
