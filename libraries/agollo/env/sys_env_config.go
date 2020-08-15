package env

import (
	"os"

	"github.com/why444216978/http-agollo/agollo/env/config"
)

type SysEnvConfig struct {
	AppID         string
	Cluster       string
	NamespaceName string
	Secret        string
	//Token        string
	IP string
}

func (sysEnvConfig *SysEnvConfig) LoadSysConfig() (*config.AppConfig, error) {
	cf := &config.AppConfig{Cluster: "default", NamespaceName: "application"}
	if sysEnvConfig.AppID != "" {
		cf.AppID = sysEnvConfig.AppID
	} else if os.Getenv("CONFIG_CENTER_APPID") != "" {
		cf.AppID = os.Getenv("CONFIG_CENTER_APPID")
	} else {
		panic("appId未设置")
	}
	if sysEnvConfig.Cluster != "" {
		cf.Cluster = sysEnvConfig.Cluster
	} else if os.Getenv("RUNTIME_CLUSTER") != "" {
		cf.Cluster = os.Getenv("RUNTIME_CLUSTER")
	}
	if sysEnvConfig.NamespaceName != "" {
		cf.NamespaceName = sysEnvConfig.NamespaceName
	}
	cf.Secret = os.Getenv("CONFIG_CENTER_TOKEN")
	//cf.Token = os.Getenv("CONFIG_CENTER_TOKEN")
	cf.IP = os.Getenv("CONFIG_CENTER_URL")
	return cf, nil
}
