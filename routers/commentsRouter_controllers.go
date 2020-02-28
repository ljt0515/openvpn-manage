package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISessionController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISessionController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISessionController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISessionController"],
		beego.ControllerComments{
			Method:           "Kill",
			Router:           `/`,
			AllowHTTPMethods: []string{"delete"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISignalController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISignalController"],
		beego.ControllerComments{
			Method:           "Send",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISysloadController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:APISysloadController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method:           "Download",
			Router:           `/certificates/:key`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/certificates`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/certificates`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})
	beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"] = append(beego.GlobalControllerRouter["github.com/ljt000/openvpn-manage/controllers:CertificatesController"],
		beego.ControllerComments{
			Method:           "Del",
			Router:           `/certificates`,
			AllowHTTPMethods: []string{"delete"},
			Params:           nil})

}
