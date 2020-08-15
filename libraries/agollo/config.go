package agollo

import (
	"fmt"

	beego_config "github.com/astaxie/beego/config"
	"gin-frame/libraries/agollo/env"
	library_config "gin-frame/libraries/config"
	util_err "gin-frame/libraries/util/error"
	"github.com/why444216978/http-agollo/agollo"
)

var AppConfig beego_config.Configer

func init() {
	section := library_config.GetConfig("app", "app")
	appId := section.Key("app_id").String()

	sysEnvConfig := &env.SysEnvConfig{AppID: appId}
	agollo.InitCustomConfig(sysEnvConfig.LoadSysConfig)
	if err := agollo.Start(); err != nil {
		panic(err)
	}
}

func Test() {
	base := agollo.GetStringValue("BASE", "")
	if base == "" {
		panic("apollo 配置为空")
	}

	var err error
	AppConfig, err = beego_config.NewConfigData("ini", []byte(base))
	util_err.Must(err)

	section, err := AppConfig.GetSection("db")
	util_err.Must(err)
	fmt.Println(section)
	fmt.Println(section["user_dynamic"])
}
