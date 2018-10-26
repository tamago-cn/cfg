package cfg

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/tamago-cn/cfg/asset"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

func webCommomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		//c.Header("Content-Type", "application/json")
	}
}

type binaryFS struct {
	fs http.FileSystem
}

func (b *binaryFS) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFS) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func binaryFileSystem(root string) *binaryFS {
	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
		Prefix:    root,
	}
	return &binaryFS{
		fs,
	}
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range asset.AssetNames() {
		file, err := asset.AssetInfo(name)
		if err != nil {
			return nil, err
		}
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") || !strings.HasPrefix(name, "templates/") {
			continue
		}
		h, err := asset.Asset(name)
		if err != nil {
			return nil, err
		}
		t, err = t.New(strings.Replace(name, "templates/", "", 1)).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//// Reload 实现route的重载方法
//func Reload() error {
//	mutex.Lock()
//	defer mutex.Unlock()
//	Destroy()
//	router = gin.Default()
//
//	router.Delims("{[", "]}")
//
//	if pathExists("static") {
//		router.Static("/static", "static")
//	} else {
//		router.Use(static.Serve("/static", binaryFileSystem("static")))
//	}
//
//	if pathExists("templates") {
//		router.LoadHTMLGlob("templates/*")
//	} else {
//		tmpl, err := loadTemplate()
//		if err != nil {
//			return err
//		}
//		router.SetHTMLTemplate(tmpl)
//	}
//
//	router.Use(webCommomHandler())
//
//	if middlewares != nil && len(middlewares) != 0 {
//		for _, mid := range middlewares {
//			router.Use(mid.middleware...)
//			log.Infof("load middleware <%s>", mid.name)
//		}
//	}
//
//	if moduleList != nil {
//		for _, module := range moduleList {
//			for _, r := range module.routes {
//				for _, m := range r.methods {
//					m.register()
//					log.Infof("load module router: %10s %-8s %s ", "<"+module.name+">", "["+m.name+"]", r.relativePath)
//				}
//			}
//		}
//	}
//	return nil
//}
//
//// Destroy 析构
//func Destroy() error {
//	router = nil
//	return nil
//}
//
//// Router 获取router
//func Router() *gin.Engine {
//	return router
//}
