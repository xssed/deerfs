package controller

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

// PageSize 每页的数据显示数量
var PageSize = 10

func PageNum(c *gin.Context) int {
	count := 0
	// c.Query("page") 取回 URL 中的参数，然后再转换成 Int
	// GET /path?id=1234&name=Manu&value=
	page := com.StrTo(c.Query("page")).MustInt()
	// page <= 1 时，count 为 0
	if page > 0 {
		// page = 2 时，count = 10
		count = (page - 1) * 10
	}
	return count
}

// OwlResponse 统一返回格式
type OwlResponse struct {
	//请求命令
	Cmd string `json:"Cmd"`
	//返回状态
	Status int `json:"Status"`
	//返回结果
	Results string `json:"Results"`
	//key
	Key string `json:"Key"`
	//返回内容
	Data interface{} `json:"Data"`
	//程序响应IP
	ResponseHost string `json:"ResponseHost"`
	//内容的创建时间
	KeyCreateTime time.Time `json:"KeyCreateTime"`
}

// Response 根据数据返回响应
func Response(c *gin.Context, httpCode, resCode int, data interface{}) {
	c.JSON(httpCode, OwlResponse{Status: resCode, Results: errno.Msg[resCode], Data: data, ResponseHost: config.SystemTagUriAddress()})
	return
}

// Response 根据文件数据返回响应
func FileResponse(c *gin.Context, httpCode int, file_md5 string, resCode int, data interface{}) {
	c.JSON(httpCode, OwlResponse{Status: resCode, Results: errno.Msg[resCode], Key: file_md5, Data: data, ResponseHost: config.SystemTagUriAddress()})
	return
}

// BindAndValid 绑定并验证表单
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	// c.Bind(form) 会根据 Content-Type 选择 binding
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, errno.InvalidParams
	}
	valid := validation.Validation{}
	// 验证该表单，必须是结构体或结构体指针
	ok, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, errno.Error
	}
	// 验证失败
	if !ok {
		LogErrors(valid.Errors)
		// for _, err := range valid.Errors {
		// 	log.Print(err.Key, err.Message)
		// }
		return http.StatusBadRequest, errno.InvalidParams
	}
	return http.StatusOK, errno.Success
}

// LogErrors 把验证错误输出到日志
func LogErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Print(err.Key, err.Message)
		loger.Lg.Error("ERROR: ", zap.String("Key", err.Key), zap.String("Message", err.Message))
	}
}

//高效将错误类型切片拼接字符串
func ErrorsSliceJoinToString(errors []*validation.Error) string {

	var args_buffer bytes.Buffer
	for _, err := range errors {
		loger.Lg.Error("ERROR: ", zap.String("Key", err.Key), zap.String("Message", err.Message))
		args_buffer.WriteString(err.Error())
	}
	return args_buffer.String()

}
