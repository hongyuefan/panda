package models

type Transtype struct {
	Id_RENAME int     `orm:"column(id);null"`
	Name      string  `orm:"column(name);size(64);null"`
	Limit     float32 `orm:"column(limit);null"`
}
