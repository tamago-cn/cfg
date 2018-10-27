package cfg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"

	log "github.com/sirupsen/logrus"
)

//go:generate go-bindata -o=asset/asset.go -pkg=asset static/... templates/...

func init() {
	iniFile = "conf/app.ini"
	cfgMap = make(map[string]interface{})
	rm = make(map[string]func() error)
	dm = make(map[string]func() error)

	conf = &Conf{
		Addr: "0.0.0.0:8086",
	}
	RegistSection("cfg", conf, nil, nil)
}

var (
	iniFile  string
	cfgMap   map[string]interface{}
	rm       map[string]func() error
	dm       map[string]func() error
	sections []string
	conf     *Conf
	server   *http.Server
)

// Conf 配置
type Conf struct {
	Addr         string        `ini:"addr" json:"addr"`
	ReadTimeout  time.Duration `ini:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `ini:"write_timeout" json:"write_timeout"`
}

//// Conf 配置
//type Conf struct {
//	Addr         string        `ini:"addr" form:"addr" json:"addr"`
//	ReadTimeout  time.Duration `ini:"read_timeout" form:"read_timeout" json:"read_timeout"`
//	WriteTimeout time.Duration `ini:"write_timeout" form:"write_timeout" json:"write_timeout"`
//}

// RegistSection 注册配置节
func RegistSection(section string, conf interface{}, reload func() error, destroy func() error) {
	if _, ok := cfgMap[section]; ok {
		panic(fmt.Errorf("conflict section: %s", section))
	}
	// 配置映射
	cfgMap[section] = conf
	// 重载函数
	rm[section] = reload
	// 析构函数
	dm[section] = destroy
	sections = append(sections, section)
}

// Load 加载配置
func Load(file string, init bool) error {
	//if file != "" {
	iniFile = file
	//}
	return Reload(init)
}

// Save 保存配置
func Save() error {
	if iniFile != "" {
		conf := ini.Empty()
		for sec, s := range cfgMap {
			section := conf.Section(sec)
			err := section.ReflectFrom(s)
			if err != nil {
				return err
			}
		}
		return conf.SaveTo(iniFile)
	}
	log.Warn("iniFile is null")
	return nil
}

// Reload 重载配置, 若iniFile为空,表示加载默认配置
func Reload(init bool) error {
	if iniFile != "" {
		conf, err := ini.Load(iniFile)
		if err != nil {
			panic(err)
		}
		for _, sec := range sections {
			s := cfgMap[sec]
			switch x := s.(type) {
			default:
				section := conf.Section(sec)
				err = section.MapTo(x)
				if err != nil {
					return err
				}
			}
		}
	}
	if init {
		for _, sec := range sections {
			reload := rm[sec]
			if reload != nil {
				err := reload()
				if err != nil {
					//return err
					panic(fmt.Errorf("error in reload %s: %s", sec, err.Error()))
				}
			}
			log.Infof("%s reload success", sec)
		}
	}
	return nil
}

// Destroy 析构函数
func Destroy() error {
	for i := len(sections) - 1; i >= 0; i-- {
		fname := sections[i]
		destroy := dm[fname]
		if destroy != nil {
			err := destroy()
			if err != nil {
				log.Errorf("destroy %s error: %s", fname, err.Error())
			}
		}
		log.Infof("%s destroyed success", fname)
	}
	return nil
}

// Run 将配置暴露在网页上使其可修改
func Run() error {
	router := gin.Default()

	if pathExists("static") {
		router.Static("/cfg/static", "static")
	} else {
		router.Use(static.Serve("/cfg/static", binaryFileSystem("static")))
	}

	if pathExists("templates") {
		router.LoadHTMLGlob("templates/*")
	} else {
		tmpl, err := loadTemplate()
		if err != nil {
			return err
		}
		router.SetHTMLTemplate(tmpl)
	}

	router.Use(webCommomHandler())

	router.GET("/cfg", indexHandler)
	router.GET("/cfg/index", indexHandler)
	router.GET("/cfg/list", confListHandler)
	router.GET("/cfg/update/:section", updateHandler)

	server = &http.Server{
		Addr:           conf.Addr,
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("http server start error:", err)
		}
	}()
	log.Infof("http server listen on %s", conf.Addr)
	return nil
}

// Stop 停止HTTP服务
func Stop() error {
	if server != nil {
		server.Close()
	}
	return nil
}
