package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
)

//获取节点系统信息
func GetSystemInfo(c *gin.Context) {

	// 填充数据
	data := map[string]interface{}{}

	//current version
	data["current_version"] = global.VERSION

	//current version date
	data["current_version_date"] = global.VERSION_DATE

	//owlcache状态值，1在线,0未在线。
	data["owlcache_state"] = global.Owlcache_State

	//Node_ID
	data["node_id"] = global.Node_ID

	//Node_Name
	data["node_name"] = config.SystemTagName()

	//Uri_Address
	data["uri_address"] = config.SystemTagUriAddress()

	//系统最大使用容量
	maxcap := file_manage.FormatFileSizeToString(config.FileStorageDirectoryStorageMaxSize())
	data["max_cap"] = maxcap
	data["max_cap_basic"] = config.FileStorageDirectoryStorageMaxSize()

	//当前节点单个文件最大的存储字节数
	file_max_size := file_manage.FormatFileSizeToString(config.FileStorageFileStorageMaxSize())
	data["file_storage_max_size"] = file_max_size
	data["file_storage_max_size_basic"] = config.FileStorageFileStorageMaxSize()

	//系统当前已使用容量
	usecap := file_manage.FormatFileSizeToString(global.Directory_Storage_Use_Size)
	data["use_cap"] = usecap
	data["use_cap_basic"] = global.Directory_Storage_Use_Size

	//对已经存储文件的数量的统计
	stored_file_number := common.Int64ToString(global.Directory_Storage_File_Number)
	data["stored_file_number"] = stored_file_number

	//当前deerfs主机的缓存目录使用量,单位字节
	cachefile_storage_use_size := file_manage.FormatFileSizeToString(global.CacheFile_Storage_Use_Size)
	data["cachefile_storage_use_size_basic"] = global.CacheFile_Storage_Use_Size
	data["cachefile_storage_use_size"] = cachefile_storage_use_size

	controller.Response(c, http.StatusOK, errno.Success, data)
}

//系统辅助小工具
func Tools(c *gin.Context) {

	//加密解密字符串
	encrypt, encrypt_exist := c.GetQuery("encrypt") //加密
	if encrypt_exist {
		c.Data(http.StatusOK, "text/html;charset=utf-8", []byte(file_manage.EnCryptToString(encrypt)))
		return
	}
	decrypt, decrypt_exist := c.GetQuery("decrypt") //解密
	if decrypt_exist {
		de, de_err := file_manage.DeCryptString(decrypt)
		if de_err != nil {
			c.Data(http.StatusServiceUnavailable, "", []byte(""))
			return
		}
		c.Data(http.StatusOK, "text/html;charset=utf-8", de)
		return
	}

}

//服务欢迎页
func SayHello(c *gin.Context) {
	c.Data(http.StatusOK, "text/html;charset=utf-8", []byte(common.JoinString("<style type='text/css'>*{ padding: 0; margin: 0; } div{ padding: 4px 48px;} a{color:#2E5CD5;cursor: pointer;text-decoration: none} a:hover{text-decoration:underline; } body{ background: #fff; font-family: 'Century Gothic','Microsoft yahei'; color: #333;font-size:18px;} h1{ font-size: 100px; font-weight: normal; margin-bottom: 12px; } p{ line-height: 1.6em; font-size: 42px }</style><div style='padding: 24px 48px;'><h1>:)</h1><p>Welcome to use deerfs. Version:", global.VERSION, "<br/><span style='font-size:25px'>If you have any questions,Please contact us: <a href=\"mailto:xsser@xsser.cc\">xsser@xsser.cc</a><br>Project Home : <a href=\"https://github.com/xssed/deerfs\" target=\"_blank\">https://github.com/xssed/deerfs</a></span></p><div>")))
}
