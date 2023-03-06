package v1

import (
	"fmt"
	"net/http"
	"path"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/vm"

	"github.com/xssed/deerfs/deerfs_service/core/common"
	//"github.com/xssed/deerfs/deerfs_service/core/system/config"
	//"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

func DownloadFile(c *gin.Context) {

	//内部处理方法，校验输入数据和查询文件信息
	file := base_input(c)

	//判断文件处理流程
	action, res_action := c.GetQuery("action")
	//如果没有执行流程则直接输出文件
	if !res_action {
		//输出文件
		c.File(path.Join(file.FileAddr))
		return
	}

	//流程执行
	if action == "imageView" {

		//过滤无意义的图像缓存生成，节约资源
		if c.Request.RequestURI != common.JoinString("/", file.FileSign, "?action=imageView") {
			ImageMainHandler(c, file)
			return
		}

	}

	//无有效指令,还是直接输出文件
	c.File(path.Join(file.FileAddr))
	return

}

//内部处理方法，校验输入数据和查询文件信息
func base_input(c *gin.Context) *mysql_model.File {

	valid := validation.Validation{}
	//获取file_sign
	file_sign := c.Param("file_sign")
	valid.MinSize(file_sign, 40, "file_sign").Message("file_sign value length error.")

	// 表单验证错误
	if valid.HasErrors() {
		controller.LogErrors(valid.Errors)
		//c.Data(http.StatusBadRequest, "text/html;charset=utf-8", []byte(controller.ErrorsSliceJoinToString(valid.Errors)))
		c.Data(http.StatusNotFound, "", []byte(""))
		return nil
	}
	// 构造结构体
	vmFile := vm.File{FileSign: file_sign}
	//查询文件数据
	file, err := vmFile.GetToSign()
	if err != nil {
		//这里的返回错误已经优化到文件不存在
		fmt.Println("get download file info error:", zap.String("error", err.Error()))
		loger.Lg.Error("get download file info error:", zap.String("error", err.Error()))
		c.Data(http.StatusNotFound, "", []byte(""))
		//c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(err.Error()))
		return nil
	}
	//判断文件是否存在
	if file == nil || file.ID < 1 {
		c.Data(http.StatusNotFound, "", []byte(""))
		return nil
	}

	return file

}
