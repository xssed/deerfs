package owlcache_model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/application/model/owlcache_model/execute"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

//连接owlcache进行ping操作
func Ping() string {

	op_ping, op_ping_err := execute.Ping("")
	if op_ping_err != nil {
		fmt.Println("ping owlcache error! error:", op_ping_err.Error())
	}
	if op_ping != nil {
		if op_ping.Status == 200 {
			if string(op_ping.Data) == "PONG" {
				global.Owlcache_State = 1
				return "ok"
			}
		} /*else {
			fmt.Println("ping owlcache error! Status:", op_ping.Status)
		}*/
	}
	global.Owlcache_State = 0
	return ""

}

//连接owlcache获取一个Token
func GetToken() string {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return ""
	}

	op_pass, op_pass_err := execute.Pass(config.OwlcachePassword())
	if op_pass_err != nil {
		fmt.Println("owlcache get token error! error:", op_pass_err.Error())
	}
	if op_pass != nil {
		if op_pass.Status == 200 {
			global.Token = string(common.Base64_Encode(op_pass.Data))
			return global.Token
		} else {
			fmt.Println("owlcache get token error! Status:", op_pass.Status)
		}
	}
	return ""

}

//连接owlcache设置一个Key/Value
func SetKeyValue(key, valuedata, exptime string) string {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return ""
	}

	op_setkey, op_setkey_err := execute.Set(key, valuedata, exptime)
	if op_setkey_err != nil {
		fmt.Println("owlcache set Key/Value error! error:", op_setkey_err.Error())
	}
	//接下来判断请求后的响应值，
	if op_setkey != nil {

		if op_setkey.Status == 200 {
			//正常
			return "ok"
		} else if op_setkey.Status == 401 {
			//如果是没有权限的错误就重新获取一次token然后再set一次
			if GetToken() == "" {
				fmt.Println("Failed to get token. Please check <owlcache_password> in the config file.Or owlcache http service exception.")
			} else {
				op_setkey2, op_setkey_err2 := execute.Set(key, valuedata, exptime)
				if op_setkey_err2 != nil {
					fmt.Println("owlcache set Key/Value error! error:", op_setkey_err2.Error())
				}
				if op_setkey2 != nil {
					if op_setkey2.Status == 200 {
						//正常
						return "ok"
					} else {
						return ""
					}
				}
			}

		} else {
			fmt.Println("owlcache set Key/Value error! Status:", op_setkey.Status)
		}
	}
	return ""

}

//连接owlcache设置一个Key的Expire时间
func SetKeyExpire(key, exptime string) string {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return ""
	}

	op_expire, op_expire_err := execute.Expire(key, exptime)
	if op_expire_err != nil {
		fmt.Println("owlcache set key expire error! error:", op_expire_err.Error())
	}
	//接下来判断请求后的响应值，
	if op_expire != nil {

		if op_expire.Status == 200 {
			//正常
			return "ok"
		} else if op_expire.Status == 401 {
			//如果是没有权限的错误就重新获取一次token然后再set一次
			if GetToken() == "" {
				fmt.Println("Failed to get token. Please check <owlcache_password> in the config file.Or owlcache http service exception.")
			} else {
				op_expire2, op_expire_err2 := execute.Expire(key, exptime)
				if op_expire_err2 != nil {
					fmt.Println("owlcache set key expire error! error:", op_expire_err2.Error())
				}
				if op_expire2 != nil {
					if op_expire2.Status == 200 {
						//正常
						return "ok"
					} else {
						return ""
					}
				}
			}

		} else {
			fmt.Println("owlcache set key expire error! Status:", op_expire.Status)
		}
	}
	return ""

}

//连接owlcache删除一个Key
func DeleteKey(key string) string {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return ""
	}

	op_delete, op_delete_err := execute.Delete(key)
	if op_delete_err != nil {
		fmt.Println("owlcache delete Key error! error:", op_delete_err.Error())
	}
	//接下来判断请求后的响应值，
	if op_delete != nil {

		if op_delete.Status == 200 {
			//正常
			return "ok"
		} else if op_delete.Status == 401 {
			//如果是没有权限的错误就重新获取一次token然后再set一次
			if GetToken() == "" {
				fmt.Println("Failed to get token. Please check <owlcache_password> in the config file.Or owlcache http service exception.")
			} else {
				op_delete2, op_delete_err2 := execute.Delete(key)
				if op_delete_err2 != nil {
					fmt.Println("owlcache delete Key error! error:", op_delete_err2.Error())
				}
				if op_delete2 != nil {
					if op_delete2.Status == 200 {
						//正常
						return "ok"
					} else {
						return ""
					}
				}
			}

		} else {
			fmt.Println("owlcache delete Key error! Status:", op_delete.Status)
		}
	}
	return ""

}

//连接owlcache进行get操作
func GetKey(key string) string {

	op_get, op_get_err := execute.Get(key)
	if op_get_err != nil {
		fmt.Println("owlcache get key error! error:", op_get_err.Error())
	}
	if op_get != nil {
		if op_get.Status == 200 {
			return string(op_get.Data)
		}
	}
	return ""

}

//根据文件Sign获取一个File信息（连接owlcache进行get操作）
func GetFileToSign(file_sign string) (*mysql_model.File, error) {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return nil, errors.New(errno.Msg[errno.OwlcacheOffline])
	}

	var file mysql_model.File
	var query_key string = common.JoinString(global.Owlcache_deerfs_key_storage_prefix, file_sign)

	op_get, op_get_err := execute.Get(query_key)
	if op_get_err != nil {
		fmt.Println("owlcache get key error! error:", op_get_err.Error())
		return nil, op_get_err
	}
	if op_get != nil {
		if op_get.Status == 200 {
			unmarshal_err := json.Unmarshal(op_get.Data, &file)
			if unmarshal_err != nil {
				fmt.Println("Unmarshal Error:", unmarshal_err)
				return nil, errors.New(common.JoinString("Unmarshal Error:", unmarshal_err.Error()))
			}
			return &file, nil
		}
	}
	return nil, errors.New(errno.Msg[errno.FileIsNotExist])
}

//向owlcache中存一个File信息(连接owlcache设置一个Key/Value)
func SetFileToSign(file_sign string, valuedata *mysql_model.File) error {

	//主机状态先行判定减少请求错误率
	if global.Owlcache_State == 0 {
		return errors.New(errno.Msg[errno.OwlcacheOffline])
	}

	var set_key string = common.JoinString(global.Owlcache_deerfs_key_storage_prefix, file_sign) //定义要存储的key名称
	valuedata.Node = *global.Node                                                                //定义要存储的key值
	jsonBytes, json_err := json.Marshal(valuedata)
	if json_err != nil {
		fmt.Println("Json Marshal Error: ", json_err.Error())
		return json_err
	}

	op_setkey, op_setkey_err := execute.Set(set_key, string(jsonBytes), config.OwlcacheKeyStorageExpire())
	if op_setkey_err != nil {
		fmt.Println("owlcache set Key/Value error! error:", op_setkey_err.Error())
		loger.Lg.Error("SetFileToSign error! error:", zap.String("error", op_setkey_err.Error()))
	}
	//接下来判断请求后的响应值，
	if op_setkey != nil {

		if op_setkey.Status == 200 {
			//正常
			return nil
		} else if op_setkey.Status == 401 {
			//如果是没有权限的错误就重新获取一次token然后再set一次
			if GetToken() == "" {
				fmt.Println("Failed to get token. Please check <owlcache_password> in the config file.Or owlcache http service exception.")
			} else {
				op_setkey2, op_setkey_err2 := execute.Set(set_key, string(jsonBytes), config.OwlcacheKeyStorageExpire())
				if op_setkey_err2 != nil {
					fmt.Println("owlcache set Key/Value error! error:", op_setkey_err2.Error())
				}
				if op_setkey2 != nil {
					if op_setkey2.Status == 200 {
						//正常
						return nil
					}
				}
			}

		} else {
			fmt.Println("owlcache set Key/Value error! Status:", op_setkey.Status)
		}
	}
	return errors.New(errno.Msg[errno.OwlcacheSetFileInfoFail])

}

//连接owlcache进行group_get操作
func GetGroupKey(key, valuedata string) string {

	op_group_get, op_group_get_err := execute.Group_Get(key, valuedata)
	if op_group_get_err != nil {
		fmt.Println("owlcache get group key error! error:", op_group_get_err.Error())
	}
	if op_group_get != nil {
		if op_group_get.Status == 200 {
			return string(op_group_get.Data)
		}
	}
	return ""

}
