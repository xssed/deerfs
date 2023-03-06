package execute

import (
	//"fmt"
	// "log"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/extend/owlcache_httpclient"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

//绑定数据到OwlResponse
func bindDataToOwlResponse(cmd CommandType, status ResStatus, results, key string, data []byte, responsehost, keycreatetime string) *OwlResponse {

	var owlres *OwlResponse = &OwlResponse{}
	owlres.Cmd = cmd
	owlres.Status = status
	owlres.Results = results
	owlres.Key = key
	owlres.Data = data
	owlres.ResponseHost = responsehost

	//时间处理部分
	if len(keycreatetime) > 37 {
		//截取字符串的固定长度格式，并把前后两端的空格过滤
		keycreatetime = strings.TrimSpace(keycreatetime[0:37])
	}
	ct, terr := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", keycreatetime)
	if terr != nil {
		loger.Lg.Info("Error parsing time string to time:", zap.String("error", terr.Error())) //日志记录
	}
	owlres.KeyCreateTime = ct

	return owlres

}

func Get(key string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Get(config.OwlcacheAddr(), key) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Get error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Get error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Get ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	return bindDataToOwlResponse(GET, ResStatus(r_res.StatusCode), "", key, body, r_res.Header.Get("Responsehost"), r_res.Header.Get("Keycreatetime")), nil

}

func Set(key, valuedata, exptime string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Set(config.OwlcacheAddr(), key, valuedata, exptime, global.Token) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Set error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Set error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Set ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	r_owl := JsonStrToOwlResponse(string(body)) //数据绑定
	return &r_owl, nil

}

func Expire(key, exptime string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Expire(config.OwlcacheAddr(), key, exptime, global.Token) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Expire error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Expire error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Expire ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	r_owl := JsonStrToOwlResponse(string(body)) //数据绑定
	return &r_owl, nil

}

func Delete(key string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Delete(config.OwlcacheAddr(), key, global.Token) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Delete error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Delete error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Delete ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	r_owl := JsonStrToOwlResponse(string(body)) //数据绑定
	return &r_owl, nil

}

func Ping(valuedata string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Ping(config.OwlcacheAddr(), valuedata) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Ping error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Ping error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Ping ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	return bindDataToOwlResponse(PING, ResStatus(r_res.StatusCode), "", "", body, r_res.Header.Get("Responsehost"), r_res.Header.Get("Keycreatetime")), nil

}

func Exists(key string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Exists(config.OwlcacheAddr(), key) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Exists error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Exists error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Exists ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	r_owl := JsonStrToOwlResponse(string(body)) //数据绑定
	return &r_owl, nil

}

func Pass(pass string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Pass(config.OwlcacheAddr(), pass) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Pass error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Pass error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Pass ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	r_owl := JsonStrToOwlResponse(string(body)) //数据绑定
	return &r_owl, nil

}

func Group_Get(key, valuedata string) (*OwlResponse, error) {

	ohc := owlcache_httpclient.New()                 //初始化gorequest库
	ohc.Timeout(config.OwlcacheHttpRequestTimeout()) //设置超时

	r_res, r_err_slices := ohc.Group_Get(config.OwlcacheAddr(), key, valuedata) //请求
	if r_err_slices != nil {
		err_slices := common.ErrorSliceJoinToString(r_err_slices)
		loger.Lg.Info("owlclient method Group_Get error:", zap.String("error", err_slices)) //日志记录
		return nil, errors.New(common.JoinString("owlclient method Group_Get error:", err_slices))
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		return nil, errors.New(common.JoinString("Group_Get ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	return bindDataToOwlResponse(GET, ResStatus(r_res.StatusCode), "", key, body, r_res.Header.Get("Responsehost"), r_res.Header.Get("Keycreatetime")), nil

}
