package models

import (
	"os"

	"openvpn-manage/lib/server/config"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gopkg.in/hlandau/passlib.v1"
)

var GlobalCfg Settings

func init() {
	initDB()
	createDefaultUsers()
	createDefaultSettings()
	createDefaultOVConfig()
}

func initDB() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	dbSource := "file:" + beego.AppConfig.String("dbPath")

	err := orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
	)

	// Database alias.
	name := "default"
	// Drop table and re-create.
	force := false
	// Print log.
	verbose := true

	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Error(err)
		return
	}
}

func createDefaultUsers() {
	hash, err := passlib.Hash("123456")
	if err != nil {
		beego.Error("Unable to hash password", err)
	}
	user := User{
		Id:       1,
		Login:    "admin",
		Name:     "超级管理员",
		Email:    "ljt_0515@163.com",
		Password: hash,
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			beego.Info("Default admin account created")
		} else {
			beego.Debug(user)
		}
	}

}

func createDefaultSettings() {
	s := Settings{
		Profile:       "default",
		MIAddress:     "127.0.0.1:5555",
		MINetwork:     "tcp",
		ServerAddress: "127.0.0.1",
		OVConfigPath:  "/etc/openvpn/server",
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		GlobalCfg = s

		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(s)
		}
	} else {
		beego.Error(err)
	}
}

func createDefaultOVConfig() {
	c := OVConfig{
		Profile: "default",
		Config: config.Config{
			Port:                1194,
			Proto:               "tcp",
			Cipher:              "AES-256-CBC",
			Auth:                "SHA256",
			Dh:                  "dh.pem",
			Keepalive:           "10 120",
			IfconfigPoolPersist: "ipp.txt",
			Management:          "0.0.0.0 5555",
			MaxClients:          100,
			Server:              "10.8.0.0 255.255.255.0",
			Ca:                  "ca.crt",
			Cert:                "server.crt",
			Key:                 "server.key",
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(c)
		}
		path := GlobalCfg.OVConfigPath + "/server.conf"
		if _, err = os.Stat(path); os.IsNotExist(err) {
			destPath := GlobalCfg.OVConfigPath + "/server.conf"
			if err = config.SaveToFile("conf/openvpn-server-config.tpl",
				c.Config, destPath); err != nil {
				beego.Error(err)
			}
		}
	} else {
		beego.Error(err)
	}
}
