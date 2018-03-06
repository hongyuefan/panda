package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Attrvalue struct {
	Id    int64  `orm:"column(id);auto"`
	Pid   int64  `orm:"column(pid)"`
	Uid   int64  `orm:"column(uid)"`
	Aid   int64  `orm:"column(aid)"`
	Years int    `orm:"column(years)"`
	Value string `orm:"column(value);size(64);null"`
	Multi string `orm:"column(multi);size(32)"`
}

func (t *Attrvalue) TableName() string {
	return "attrvalue"
}

func init() {
	orm.RegisterModel(new(Attrvalue))
}

func NewAttrvalue(pid, uid, aid int64, years int, value string) *Attrvalue {
	return &Attrvalue{
		Pid:   pid,
		Uid:   uid,
		Aid:   aid,
		Years: years,
		Value: value,
	}
}

// AddAttrvalue insert a new Attrvalue into database and returns
// last inserted Id on success.
func AddAttrvalue(m *Attrvalue) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAttrvalueById retrieves Attrvalue by Id. Returns error if
// Id doesn't exist
func GetAttrvalueById(id int64) (v *Attrvalue, err error) {
	o := orm.NewOrm()
	v = &Attrvalue{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAttrvalue retrieves all Attrvalue matches certain condition. Returns empty list if
// no records exist
func GetAllAttrvalue(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Attrvalue))
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

	var l []Attrvalue
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

// UpdateAttrvalue updates Attrvalue by Id and returns error if
// the record to be updated doesn't exist
func UpdateAttrvalueById(m *Attrvalue) (err error) {
	o := orm.NewOrm()
	v := Attrvalue{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAttrvalue deletes Attrvalue by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAttrvalue(id int64) (err error) {
	o := orm.NewOrm()
	v := Attrvalue{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Attrvalue{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
