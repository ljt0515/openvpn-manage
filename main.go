package main

import (
	"github.com/astaxie/beego"
	"openvpn-manage/lib"
	_ "openvpn-manage/routers"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
