package main

import (
	"github.com/astaxie/beego"
	"github.com/ljt000/openvpn-manage/lib"
	_ "github.com/ljt000/openvpn-manage/routers"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
