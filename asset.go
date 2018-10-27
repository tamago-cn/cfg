package cfg

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/tamago-cn/cfg/asset"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

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
