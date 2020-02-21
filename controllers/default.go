package controllers

import (
	"github.com/astaxie/beego"
	"openvpn-manage/lib"
	"openvpn-manage/lib/server/mi"
	"openvpn-manage/models"
)

type MainController struct {
	BaseController
}

func (c *MainController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Status",
	}
}

func (c *MainController) Get() {
	c.Data["sysInfo"] = lib.GetSystemInfo()
	lib.Dump(lib.GetSystemInfo())
	client := mi.NewClient(models.GlobalCfg.MINetwork, models.GlobalCfg.MIAddress)
	status, err := client.GetStatus()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovStatus"] = status
		c.Data["ovVersion"] = status.Title
	}
	lib.Dump(status)

	pid, err := client.GetPid()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovPid"] = pid
	}
	lib.Dump(pid)

	loadStats, err := client.GetLoadStats()
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["ovStats"] = loadStats
	}
	lib.Dump(loadStats)

	c.TplName = "index.html"
}
