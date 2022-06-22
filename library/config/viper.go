package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Viper struct {
	*viper.Viper
	path string
}

var env = flag.String("env", "dev", "config path")

var Env = map[string]struct{}{
	"dev":      {},
	"liantiao": {},
	"qa":       {},
	"online":   {},
}

var defaultConf = New(*env)

func init() {
	flag.Parse()
}

func ReadConfig(file, typ string, data interface{}) (err error) {
	return defaultConf.ReadConfig(file, typ, data)
}

func Path() string {
	return defaultConf.Path()
}

func Dir() (string, error) {
	return filepath.Abs(Path())
}

func Config() *Viper {
	return defaultConf
}

func New(env string) *Viper {
	log.Println("The environment is :" + env)

	if _, ok := Env[env]; !ok {
		panic(env + " error")
	}

	path := "conf/" + env
	if _, err := os.Stat(path); err != nil {
		panic(err)
	}

	config := viper.New()
	config.AddConfigPath(path)

	return &Viper{
		Viper: config,
		path:  path,
	}
}

func (v *Viper) ReadConfig(file, typ string, data interface{}) (err error) {
	v.SetConfigName(file)
	v.SetConfigType(typ)
	if err = v.ReadInConfig(); err != nil {
		return
	}

	return v.Unmarshal(&data)
}

func (v *Viper) Path() string {
	return v.path
}
