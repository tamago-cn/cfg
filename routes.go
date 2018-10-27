package cfg

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func getInputType(f reflect.StructField) string {
	switch f.Type.Kind() {
	case reflect.Bool:
		return "switch"
	case reflect.Slice:
		return "list"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "text"
	}
}

func getValue(v reflect.Value) interface{} {
	switch v.Type().Kind() {
	case reflect.Slice:
		vs := []*ValueItem{
			&ValueItem{
				Type:  "text",
				Value: "xxxx",
			},
		}
		return vs
	case reflect.Ptr:
		// 指针
		return v.Interface()
	default:
		return v.Interface()

	}
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// Field 字段
type Field struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Value   interface{} `json:"value"`
	Comment string      `json:"comment"`
}

// ValueItem 当配置项为列表时使用单个Value使用此结构
type ValueItem struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Section 配置节
type Section struct {
	Title  string   `json:"title"`
	Fields []*Field `json:"fields"`
}

func confListHandler(c *gin.Context) {
	confs := []*Section{}
	for _, s := range sections {
		fields := []*Field{}
		cf := cfgMap[s]
		switch cs := cf.(type) {
		default:
			ct := reflect.TypeOf(cs)
			cv := reflect.ValueOf(cs)
			if ct.Kind() == reflect.Ptr {
				ct = ct.Elem()
				cv = cv.Elem()
			} else {
				RenderErrMsg(c, "non-pointer struct")
				return
			}

			for i := 0; i < ct.NumField(); i++ {
				f := ct.Field(i)
				v := cv.Field(i)
				name := f.Tag.Get("json")
				if name == "" {
					// 没有json标签的不支持在网页修改
					continue
				}
				fields = append(fields, &Field{
					Name:    name,
					Type:    getInputType(f),
					Value:   getValue(v),
					Comment: f.Tag.Get("comment"),
				})
			}
		}
		if len(fields) == 0 {
			continue
		}
		sec := &Section{
			Title:  s,
			Fields: fields,
		}
		confs = append(confs, sec)
	}
	RenderSuccess(c, confs)
}

func updateHandler(c *gin.Context) {
	s := c.Param("section")
	if cf, ok := cfgMap[s]; ok {
		switch x := cf.(type) {
		default:
			fmt.Println(x)
			err := c.BindJSON(x)
			if err != nil {
				RenderErrMsg(c, err.Error())
				return
			}
			err = Save()
			if err != nil {
				RenderErrMsg(c, err.Error())
				return
			}
			RenderSuccess(c, gin.H{
				s: x,
			})
			return
		}
	}
	RenderErrMsg(c, "unsupported section")
}
