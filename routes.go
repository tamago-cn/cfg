package cfg

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Input 单个配置节点
type Input struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Value   interface{} `json:"value"`
	Comment string      `json:"comment"`
}

// Form 配置表单
type Form struct {
	Title  string   `json:"title"`
	Inputs []*Input `json:"inputs"`
}

func getInputType(f reflect.StructField) string {
	switch f.Type.Kind() {
	case reflect.Bool:
		return "checkbox"
	default:
		return "text"
	}
}

func getValue(v reflect.Value) interface{} {
	return v.Interface()
}

func indexHandler(c *gin.Context) {
	forms := make([]*Form, 0, len(sections))
	for _, s := range sections {
		inputs := []*Input{}
		cf := cfgMap[s]
		switch cs := cf.(type) {
		default:
			ct := reflect.TypeOf(cs)
			cv := reflect.ValueOf(cs)
			if ct.Kind() == reflect.Ptr {
				ct = ct.Elem()
				cv = cv.Elem()
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"error_msg": "non-pointer struct",
				})
				return
			}

			for i := 0; i < ct.NumField(); i++ {
				f := ct.Field(i)
				v := cv.Field(i)
				name := f.Tag.Get("form")
				if name == "" {
					// 不绑定form的不展示在网页上
					continue
				}
				inputs = append(inputs, &Input{
					Name:    name,
					Type:    getInputType(f),
					Value:   getValue(v),
					Comment: f.Tag.Get("comment"),
				})
			}
		}
		if len(inputs) == 0 {
			continue
		}
		form := &Form{
			Title:  s,
			Inputs: inputs,
		}
		forms = append(forms, form)
	}
	c.HTML(http.StatusOK, "index.tmpl", forms)
}

func updateHandler(c *gin.Context) {
	s := c.Param("section")
	fmt.Println(s)
	if cf, ok := cfgMap[s]; ok {
		switch x := cf.(type) {
		default:
			err := c.BindQuery(x)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error_msg": err.Error(),
				})
				return
			}
			err = Save()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error_msg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				s: x,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"error_msg": "unsupported section",
	})
}
