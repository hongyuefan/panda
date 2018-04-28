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
	orm.RegisterDataBase("default", "mysql", "root:350999@tcp(47.92.67.93:3306)/panda")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.SetStaticPath(types.Svg_File_Path, beego.AppConfig.String("svg_path"))
	beego.SetStaticPath(types.Pic_File_Path, beego.AppConfig.String("pic_path"))

	beego.Run()
}
