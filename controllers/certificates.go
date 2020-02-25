package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"openvpn-manage/lib"
	"openvpn-manage/lib/client/config"
	"openvpn-manage/models"
)

type NewCertParams struct {
	Name string `form:"Name" valid:"Required;"`
}

type CertificatesController struct {
	BaseController
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "证书",
	}
}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "application/zip")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	if cfgPath, err := saveClientConfig(name); err == nil {
		c.Ctx.Output.Download(cfgPath, filename)
	}

}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := lib.ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	lib.Dump(certs)
	c.Data["certificates"] = &certs
}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if lib.CreateCertificate(cParams.Name) {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}
func (c *CertificatesController) Del() {
	r := NewJSONResponse()
	r.Data = "error"
	c.Data["json"] = r
	cParams := NewCertParams{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cParams); err != nil {
		c.ServeJSON()
		return
	}
	if vMap := validateCertParams(cParams); vMap != nil {
		c.Data["json"] = vMap
		c.ServeJSON()
		return
	} else {
		if lib.DelCertificate(cParams.Name) {
			c.ServeJSON()
			return
		}
	}
	r.Data = "success"
	c.Data["json"] = r
	c.ServeJSON()
}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		beego.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func saveClientConfig(name string) (string, error) {
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/"
	cfg := config.New()
	cert, _ := lib.ReadLine(path + "issued/" + name + ".crt")
	cfg.Cert = cert
	key, _ := lib.ReadLine(path + "private/" + name + ".key")
	cfg.Key = key
	ca, _ := lib.ReadLine(path + "ca.crt")
	cfg.Ca = ca
	tlsCrypt, _ := lib.ReadLine(models.GlobalCfg.OVConfigPath + "tc.key")
	cfg.TlsCrypt = tlsCrypt
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher

	destPath := models.GlobalCfg.OVConfigPath + "ovpn/" + name + ".ovpn"
	if err := config.SaveToFile("conf/openvpn-client-config.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}
	return destPath, nil
}
