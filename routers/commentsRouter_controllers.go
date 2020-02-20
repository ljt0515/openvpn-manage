package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["openvpn-manage/controllers:APISessionController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:APISessionController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:APISessionController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:APISessionController"],
		beego.ControllerComments{
			Method: "Kill",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:APISignalController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:APISignalController"],
		beego.ControllerComments{
			Method: "Send",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:APISysloadController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:APISysloadController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method: "Download",
			Router: `/certificates/:key`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/certificates`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/certificates`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
