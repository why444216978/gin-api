package config

import (
	"github.com/spf13/viper"
)

type Viper struct {
	viper *viper.Viper
}

func InitConfig(path string, typ string) *Viper {
	config := viper.New()
	config.AddConfigPath(path)

	return &Viper{
		viper: config,
	}
}

func (v *Viper) ReadConfig(file, typ string, data interface{}) error {
	v.viper.SetConfigName(file)
	v.viper.SetConfigType(typ)
	v.viper.ReadInConfig()

	return v.viper.Unmarshal(&data)
}

func (v *Viper) GetString(key string) string {
	return v.viper.GetString(key)
}
