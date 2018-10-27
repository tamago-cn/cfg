package cfg

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func webCommomHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		//c.Header("Content-Type", "application/json")
	}
}

// ParserObjectFromReader 定义了解析参数到结构体的方法
func ParserObjectFromReader(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&obj)
}

// RenderSuccess 定义了返回成功的方法
func RenderSuccess(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, data)
}

// RenderErrMsg 定义了返回错误信息的方法
func RenderErrMsg(c *gin.Context, errMsg string) {
	result := gin.H{
		"success":   false,
		"error_msg": errMsg,
	}
	//c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, result)
}

// RenderErrMsgWithErrCode 定义了返回错误信息的方法
func RenderErrMsgWithErrCode(c *gin.Context, errMsg string, errCode int32) {
	result := gin.H{
		"success":    false,
		"error_msg":  errMsg,
		"error_code": errCode,
	}
	//c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}
