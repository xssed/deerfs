package main

import (
	"encoding/json"
	"net/http"
	"time"
)

//request command type
type CommandType string

//response status type
type ResStatus int

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
func JsonByteToOwlResponse(j []byte) OwlResponse {

	var tmp OwlResponse
	_ = json.Unmarshal(j, &tmp)
	// if err != nil {
	// 	fmt.Println(tools.JoinString("JsonStrToOwlResponse error:", err.Error()))
	// }
	return tmp

}

//定义响应类
type Response struct {
	*http.Response
}
