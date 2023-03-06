package execute

import (
	"encoding/json"
	"time"

	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

type OwlResponse struct {
	//请求命令
	Cmd CommandType
	//返回状态
	Status ResStatus
	//返回结果
	Results string
	//key
	Key string
	//返回内容
	Data []byte
	//程序响应IP
	ResponseHost string
	//内容的创建时间
	KeyCreateTime time.Time
}

//函数:将服务之间返回的字符串转换回结构体
func JsonStrToOwlResponse(json_string string) OwlResponse {

	var tmp OwlResponse
	err := json.Unmarshal([]byte(json_string), &tmp)
	if err != nil {
		loger.Lg.Info("JsonStrToOwlResponse error:", zap.String("error", err.Error())) //日志记录
	}
	return tmp

}
