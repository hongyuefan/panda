package main

import (
	_ "panda/routers"
	_ "panda/transaction"

	"panda/types"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root:350999@tcp(localhost:3306)/panda")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.SetStaticPath(types.Svg_File_Path, beego.AppConfig.String("svg_path"))

	beego.Run(beego.AppConfig.String("httpport"))
}
