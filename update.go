package cfg

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Input 单个配置节点
type Input struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

// Form 配置表单
type Form struct {
	Title  string   `json:"title"`
	Inputs []*Input `json:"inputs"`
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
			fmt.Println(ct)
			fmt.Println(cv)

			for i := 0; i < ct.NumField(); i++ {
				f := ct.Field(i)
				v := cv.Field(i)
				inputs = append(inputs, &Input{
					Name:    f.Tag.Get("form"),
					Type:    f.Type.String(),
					Value:   v.String(),
					Comment: f.Tag.Get("comment"),
				})
			}
		}
		form := &Form{
			Title:  s,
			Inputs: inputs,
		}
		forms = append(forms, form)
	}
	c.HTML(http.StatusOK, "index.tmpl", forms)
}

type updateReq struct {
	Section string                 `json:"section"`
	Conf    map[string]interface{} `json:"conf"`
}

func updateHandler(c *gin.Context) {
	s := c.Param("section")
	fmt.Println(s)
	//req := &updateReq{}
	//err := c.BindJSON(req)
	//req := map[string]interface{}{}
	//err := c.BindQuery(&req)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error_msg": err.Error(),
	//	})
	//	return
	//}
	//fmt.Println(req)
	c.JSON(http.StatusOK, gin.H{})
	//if s, ok := cfgMap[req.Section]; ok {
	//	// TODO: 更新配置
	//	//for k, v := range req.Conf {
	//	//	fmt.Println(s, k, v)
	//	//}
	//	c.JSON(http.StatusOK, gin.H{})
	//	return
	//}
	//c.JSON(http.StatusBadRequest, gin.H{
	//	"error_msg": fmt.Sprintf("unsupported section: '%s'", req.Section),
	//})
}
