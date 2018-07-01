package cfg

import (
	"fmt"

	"github.com/go-ini/ini"

	// 推荐使用此日志库，扩展性非常好
	log "github.com/sirupsen/logrus"
)

func init() {
	iniFile = "conf/app.ini"
	c = make(map[string]interface{})
	rm = make(map[string]func() error)
	dm = make(map[string]func() error)
}

var (
	iniFile  string
	c        map[string]interface{}
	rm       map[string]func() error
	dm       map[string]func() error
	sections []string
)

// RegistSection 注册配置节
func RegistSection(section string, conf interface{}, reload func() error, destroy func() error) {
	if _, ok := c[section]; ok {
		panic(fmt.Errorf("conflict section: %s", section))
	}
	// 配置映射
	c[section] = conf
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
		for sec, s := range c {
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
			s := c[sec]
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
