package module2

import (
	"github.com/tamago-cn/cfg"

	log "github.com/sirupsen/logrus"
)

func init() {
	conf = &Conf{
		WorkDir: "/tmp/module2",
	}
	cfg.RegistSection("module2", conf, Reload, Destroy)
}

var conf *Conf

// Conf 模块配置
type Conf struct {
	WorkDir string `ini:"work_dir" form:"work_dir" comment:"工作目录"`
}

// Reload 模块初始化入口
func Reload() error {
	log.Info("Reload module2")
	return nil
}

// Destroy 模块析构入口
func Destroy() error {
	log.Info("Destroy module2")
	return nil
}
