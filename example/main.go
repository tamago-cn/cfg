package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// 推荐使用此日志库，扩展性非常好, 在业务模块前引用logger即可实现对此日志库的相关配置
	//log "github.com/sirupsen/logrus"
	_ "github.com/tamago-cn/logger"

	// 使用golang的模块加载特性，用_引用需要的模块即可灵活选择要使用的模块
	_ "github.com/tamago-cn/cfg/example/module1"
	_ "github.com/tamago-cn/cfg/example/module2"

	"github.com/tamago-cn/cfg"
)

var conf = flag.String("conf", "conf/app.ini", "全局配置")

func main() {

	flag.Parse()

	cfg.Load(*conf, false)
	// TODO: Load 或Reload 之后将main线程阻塞，以达到启动服务的效果
	cfg.Reload(true)

	cfg.Run()

	defer cfg.Destroy()

	cfg.Save()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("received signal:", <-ch)
}
