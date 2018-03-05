package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Player struct {
	Id            int64  `orm:"column(id);auto"`
	Username      string `orm:"column(username);size(32)"`
	Password      string `orm:"column(password);size(64)"`
	Nickname      string `orm:"column(nickname);size(32);null"`
	Createtime    int64  `orm:"column(createtime)"`
	Lastime       int64  `orm:"column(lastime)"`
	Mobile        string `orm:"column(mobile);size(32);null"`
	Email         string `orm:"column(email);size(64);null"`
	Pubkey        string `orm:"column(pubkey);size(128);null"`
	Isdel         int8   `orm:"column(isdel);null"`
	Paypass       string `orm:"column(paypass);size(128);null"`
	UserType      string `orm:"column(usertype);size(32);null"`
	Avatar        string `orm:"column(avatar);size(256);null"`
	PubPublic     string `orm:"column(pub_pubkey);size(128);null"`
	PubPrivkey    string `orm:"column(pub_privkey);size(128);null"`
	LastCatchTime int64  `orm:"column(lastcatchtime)"`
}

func (t *Player) TableName() string {
	return "player"
}

func init() {
	orm.RegisterModel(new(Player))
}

// AddPlayer insert a new Player into database and returns
// last inserted Id on success.
func AddPlayer(m *Player) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPlayerById retrieves Player by Id. Returns error if
// Id doesn't exist
func GetPlayerById(id int64) (v *Player, err error) {
	o := orm.NewOrm()
	v = &Player{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPlayer retrieves all Player matches certain condition. Returns empty list if
// no records exist
func GetAllPlayer(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Player))
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

	var l []Player
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

// UpdatePlayer updates Player by Id and returns error if
// the record to be updated doesn't exist
func UpdatePlayerById(m *Player, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Player{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePlayer deletes Player by Id and returns error if
// the record to be deleted doesn't exist
func DeletePlayer(id int64) (err error) {
	o := orm.NewOrm()
	v := Player{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Player{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func UpCatchTime(uid int64) (err error) {
	var (
		player *Player
	)
	if player, err = GetPlayerById(uid); err != nil {
		return
	}
	player.LastCatchTime = time.Now().Unix()

	if err = UpdatePlayerById(player, "LastCatchTime"); err != nil {
		return
	}
	return
}
