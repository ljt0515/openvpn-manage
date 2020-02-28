package routers

import (
	"github.com/astaxie/beego"
	"github.com/ljt000/openvpn-manage/controllers"
)

func init() {
	beego.SetStaticPath("/swagger", "swagger")
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{}, "get,post:Login")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("/profile", &controllers.ProfileController{})
	beego.Router("/settings", &controllers.SettingsController{})
	beego.Router("/ov/config", &controllers.OVConfigController{})
	beego.Router("/logs", &controllers.LogsController{})

	beego.Include(&controllers.CertificatesController{})

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/session",
			beego.NSInclude(
				&controllers.APISessionController{},
			),
		),
		beego.NSNamespace("/sysload",
			beego.NSInclude(
				&controllers.APISysloadController{},
			),
		),
		beego.NSNamespace("/signal",
			beego.NSInclude(
				&controllers.APISignalController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
