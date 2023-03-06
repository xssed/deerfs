package owlcache_httpclient

import (
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/xssed/deerfs/deerfs_service/core/common"
)

//定义OwlcacheHttpclient
type OwlcacheHttpclient struct {
	sa *gorequest.SuperAgent
}

//初始化gorequest库
func New() *OwlcacheHttpclient {
	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()
	return &OwlcacheHttpclient{sa: grsa}
}

//设置超时
func (ohc *OwlcacheHttpclient) Timeout(owlcache_http_request_timeout time.Duration) {
	ohc.sa.Timeout(owlcache_http_request_timeout * time.Millisecond)
}

//清理资源
func (ohc *OwlcacheHttpclient) ClearSuperAgent() {
	ohc.sa.ClearSuperAgent()
}

func (ohc *OwlcacheHttpclient) Get(addr, key string) (*Response, []error) {

	ohc.sa.Get(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "get")
	ohc.sa.Param("key", key)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Set(addr, key, valuedata, exptime, token string) (*Response, []error) {

	ohc.sa.Post(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "set")
	ohc.sa.Param("key", key)
	ohc.sa.Param("valuedata", valuedata)
	ohc.sa.Param("exptime", exptime)
	ohc.sa.Param("token", token)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Exists(addr, key string) (*Response, []error) {

	ohc.sa.Get(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "exist")
	ohc.sa.Param("key", key)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Pass(addr, pass string) (*Response, []error) {

	ohc.sa.Post(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "pass")
	ohc.sa.Param("pass", pass)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Delete(addr, key, token string) (*Response, []error) {

	ohc.sa.Delete(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "delete")
	ohc.sa.Param("key", key)
	ohc.sa.Param("token", token)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Expire(addr, key, exptime, token string) (*Response, []error) {

	ohc.sa.Put(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "expire")
	ohc.sa.Param("key", key)
	ohc.sa.Param("exptime", exptime)
	ohc.sa.Param("token", token)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Ping(addr, valuedata string) (*Response, []error) {

	ohc.sa.Get(common.JoinString(addr, "/data/"))
	ohc.sa.Param("cmd", "ping")
	if len(valuedata) > 0 {
		ohc.sa.Param("valuedata", valuedata)
	}

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

func (ohc *OwlcacheHttpclient) Group_Get(addr, key, valuedata string) (*Response, []error) {

	ohc.sa.Get(common.JoinString(addr, "/group_data/"))
	ohc.sa.Param("cmd", "get")
	ohc.sa.Param("key", key)
	if len(valuedata) > 0 {
		ohc.sa.Param("valuedata", "info")
	}

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}
