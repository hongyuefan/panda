package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type YearsCount struct {
	Id    int    `orm:column(id)`
	Count int    `orm:"column(count)"`
	Limit int    `orm:"column(limit)"`
	Range string `orm:"column(range)"`
}

func (t *YearsCount) TableName() string {
	return "years_count"
}

func init() {
	orm.RegisterModel(new(YearsCount))
}

func GetAllYearsCount() (yearsCounts []*YearsCount, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(&YearsCount{}).All(yearsCounts)
	return
}

func IsYearsOver(years int) (err error) {
	o := orm.NewOrm()
	v := &YearsCount{Id: years}
	if err = o.Read(v); err != nil {
		return
	}
	if v.Count >= v.Limit {
		err = fmt.Errorf("years %v has full", years)
		return
	}
	return nil
}

func AddYearsCount(years int) (err error) {
	o := orm.NewOrm()
	v := &YearsCount{Id: years}
	if err = o.Read(v); err != nil {
		if err == orm.ErrNoRows {
			v := &YearsCount{
				Id:    years,
				Count: 1,
				Limit: 500,
				Range: "0_0",
			}
			_, err = o.Insert(v)
			return
		}
		return
	}
	v.Count += 1
	_, err = o.Update(v, "count")
	return
}
