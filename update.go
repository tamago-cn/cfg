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
		ct := reflect.TypeOf(cf)
		cv := reflect.ValueOf(cf)
		fmt.Println(ct)
		fmt.Println(cv)

		inputs = append(inputs, &Input{
			Name: "test",
			Type: "text",
		})
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
