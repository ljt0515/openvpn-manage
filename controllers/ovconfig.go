package controllers

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"openvpn-manage/lib"
	"openvpn-manage/lib/server/config"
	mi "openvpn-manage/lib/server/mi"
	"openvpn-manage/models"
)

type OVConfigController struct {
	BaseController
}

func (c *OVConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "OpenVPN 配置",
	}
}

func (c *OVConfigController) Get() {
	c.TplName = "ovconfig.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	cfg := models.OVConfig{Profile: "default"}
	cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *OVConfigController) Post() {
	c.TplName = "ovconfig.html"
	flash := beego.NewFlash()
	cfg := models.OVConfig{Profile: "default"}
	cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		beego.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPath := models.GlobalCfg.OVConfigPath + "/server.conf"
	err := config.SaveToFile("conf/openvpn-server-config.tpl", cfg.Config, destPath)
	if err != nil {
		beego.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("配置已更新")
		client := mi.NewClient(models.GlobalCfg.MINetwork, models.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("配置已更新，但未重新加载OpenVPN服务器： " + err.Error())
		}
	}
	flash.Store(&c.Controller)
}
