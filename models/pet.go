package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Pet struct {
	Id            int64  `orm:"column(id);auto" description:"自增"`
	Petname       string `orm:"column(petname);size(64);null" description:"宠物名称"`
	Years         int    `orm:"column(years)" description:"第几代"`
	Uid           int64  `orm:"column(uid)" description:"用户ID"`
	Cid           int64  `orm:"column(cid)"`
	Fid           int64  `orm:"column(fid)"`
	Status        int    `orm:"column(status)"`
	SvgPath       string `orm:"column(svg_path);size(256)"`
	TrainTotle    string `orm:"column(train_totle);size(128)"`
	LastCatchTime int64  `orm:"column(lastcatchtime)"`
	CreatTime     int64  `orm:"column(createtime)"`
	CatchTimes    int    `orm:"column(catch_times)"`
}

func (t *Pet) TableName() string {
	return "pet"
}

func init() {
	orm.RegisterModel(new(Pet))
}

// AddPet insert a new Pet into database and returns
// last inserted Id on success.
func AddPet(m *Pet) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPetById retrieves Pet by Id. Returns error if
// Id doesn't exist
func GetPetById(id int64) (v *Pet, err error) {
	o := orm.NewOrm()
	v = &Pet{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPet retrieves all Pet matches certain condition. Returns empty list if
// no records exist
func GetAllPet(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Pet))
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

	var l []Pet
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

// UpdatePet updates Pet by Id and returns error if
// the record to be updated doesn't exist
func UpdatePetById(m *Pet, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Pet{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePet deletes Pet by Id and returns error if
// the record to be deleted doesn't exist
func DeletePet(id int64) (err error) {
	o := orm.NewOrm()
	v := Pet{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Pet{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func UpCatchTime(pid int64) (err error) {
	var pet *Pet
	if pet, err = GetPetById(pid); err != nil {
		return
	}
	pet.LastCatchTime = time.Now().Unix()
	pet.CatchTimes += 1

	if err = UpdatePetById(pet, "LastCatchTime", "CatchTimes"); err != nil {
		return
	}
	return
}

func UpTrainTotle(pid int64, addAmount string) (err error) {

	pet, err := GetPetById(pid)
	if err != nil {
		return
	}

	nTotle, err := strconv.ParseFloat(pet.TrainTotle, 64)
	if err != nil {
		return
	}
	nAmount, err := strconv.ParseFloat(addAmount, 64)
	if err != nil {
		return
	}

	Totle := nTotle + nAmount

	pet.TrainTotle = fmt.Sprintf("%v", Totle)

	return UpdatePetById(pet, "TrainTotle")
}
