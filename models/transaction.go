package models

type Transaction struct {
	Id_RENAME   int    `orm:"column(id);null"`
	Txhash      string `orm:"column(txhash);size(128);null"`
	Transtypeid int    `orm:"column(transtypeid);null"`
	Status      int    `orm:"column(status);null"`
}
