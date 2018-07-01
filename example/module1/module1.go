package module1

import (
	"github.com/tamago-cn/cfg"

	log "github.com/sirupsen/logrus"
)

func init() {
	conf = &Conf{
		Host: "127.0.0.1",
		Port: 3306,
	}
	cfg.RegistSection("module1", conf, Reload, Destroy)
}

var conf *Conf

// Conf 模块配置
type Conf struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

// Reload 模块初始化入口
func Reload() error {
	log.Info("Reload module1")
	return nil
}

// Destroy 模块析构入口
func Destroy() error {
	log.Info("Destroy module1")
	return nil
}
