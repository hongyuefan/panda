package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Attribute struct {
	Id       int64  `orm:"column(id);auto"`
	Attrname string `orm:"column(attrname);size(32);null"`
	Normal   string `orm:"column(normal);size(32);null"`
	Special  string `orm:"column(special);size(32)"`
	Limit    string `orm:"column(limit);size(64)"`
}

func (t *Attribute) TableName() string {
	return "attribute"
}

func init() {
	orm.RegisterModel(new(Attribute))
}

func GetAllAttribute(arryAttr interface{}) (int64, error) {
	com := NewCommon()
	return com.CommonGetAll("attribute", arryAttr)
}

// AddAttribute insert a new Attribute into database and returns
// last inserted Id on success.
func AddAttribute(m *Attribute) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAttributeById retrieves Attribute by Id. Returns error if
// Id doesn't exist
func GetAttributeById(id int64) (v *Attribute, err error) {
	o := orm.NewOrm()
	v = &Attribute{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateAttribute updates Attribute by Id and returns error if
// the record to be updated doesn't exist
func UpdateAttributeById(m *Attribute) (err error) {
	o := orm.NewOrm()
	v := Attribute{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAttribute deletes Attribute by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAttribute(id int64) (err error) {
	o := orm.NewOrm()
	v := Attribute{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Attribute{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
