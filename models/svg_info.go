package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Svg_info struct {
	S_id        int64  `orm:"column(s_id);auto"`
	S_name      string `orm:"column(s_name);size(32)"`
	Catagory_id int    `orm:"column(catagory_id)"`
	Gender      string `orm:"column(gender);size(1)"`
	Svg_dtl     string `orm:"column(svg_dtl);size(4096)"`
	N_id        int64  `orm:"column(n_id);"`
	Base_color  int    `orm:"column(base_color);"`
	P_id        int64  `orm:"column(p_id);"`
	Link_id     string `orm:"column(link_id);"`
}

func (t *Svg_info) TableName() string {
	return "svg_info"
}

func init() {
	orm.RegisterModel(new(Svg_info))
}

// last inserted Id on success.
func AddSvginfo(m *Svg_info) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSvginfoById retrieves Svg_info by Id. Returns error if
// Id doesn't exist
func GetSvginfoById(id int64) (v *Svg_info, err error) {
	o := orm.NewOrm()
	v = &Svg_info{S_id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetCountByCatagoryId(id int64) (count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Svg_info))

	qs = qs.Filter("Catagory_id", id)
	qs = qs.Filter("N_id", 0)

	count, _ = qs.Count()

	return
}

// GetSvginfoByCatagoryId retrieves Svg_info by Id. Returns error if
// Id doesn't exist
func GetSvginfoByBasecolor(base_color int) (v *Svg_info, err error) {
	o := orm.NewOrm()
	v = &Svg_info{Base_color: base_color}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSvginfo retrieves all Svg_info matches certain condition. Returns empty list if
// no records exist
func GetAllSvginfo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Svg_info))
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

	var l []Svg_info
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

// UpdateSvginfo updates Svg_info by Id and returns error if
// the record to be updated doesn't exist
func UpdateSvginfoById(m *Svg_info) (err error) {
	o := orm.NewOrm()
	v := Svg_info{S_id: m.S_id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSvginfo deletes Svg_info by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSvginfo(id int64) (err error) {
	o := orm.NewOrm()
	v := Svg_info{S_id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Svg_info{S_id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
