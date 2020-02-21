package controllers

import (
	"openvpn-manage/lib/server/mi"
	"strings"

	"github.com/astaxie/beego"
	"openvpn-manage/models"
)

type LogsController struct {
	BaseController
}

func (c *LogsController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
}

func (c *LogsController) Get() {
	c.TplName = "logs.html"
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "日志",
	}

	settings := models.Settings{Profile: "default"}
	settings.Read("Profile")

	if err := settings.Read("OVConfigPath"); err != nil {
		beego.Error(err)
		return
	}

	client := mi.NewClient(models.GlobalCfg.MINetwork, models.GlobalCfg.MIAddress)
	getLogs, _ := client.GetLogs()
	var logs = strings.Split(getLogs, "\n")

	start := len(logs) - 200
	if start < 0 {
		start = 0
	}
	c.Data["logs"] = logs
}

func reverse(lines []string) []string {
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - i - 1
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}
