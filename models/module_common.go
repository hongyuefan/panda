package models

import (
	"github.com/astaxie/beego/orm"
)

type Common struct {
	o orm.Ormer
}

func NewCommon() *Common {
	return &Common{
		o: orm.NewOrm(),
	}
}
func (c *Common) CommonGetOne(v interface{}, col ...string) (err error) {
	if err = c.o.Read(v, col...); err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
	}
	return
}

func (c *Common) CommonInsert(v interface{}) (id int64, err error) {
	if id, err = c.o.Insert(v); err != nil {
		return 0, err
	}
	return
}

func (c *Common) CommonUpdate(v interface{}, col ...string) (int64, error) {
	return c.o.Update(v, col...)
}

func (c *Common) BeginTx() error {
	return c.o.Begin()
}
func (c *Common) Rollback() error {
	return c.o.Rollback()
}
func (c *Common) Commit() error {
	return c.o.Commit()
}
