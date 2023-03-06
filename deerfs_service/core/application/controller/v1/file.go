package v1

import (
	"net/http"
	//"time"
	//"fmt"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	//"github.com/unknwon/com"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/vm"

	//"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
)

//获取文件信息集合(当前Node下)
func GetFiles(c *gin.Context) {
	// 构造结构体
	vmFile := vm.File{PageNum: controller.PageNum(c), PageSize: controller.PageSize}

	files, err := vmFile.GetAll()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileListFail, nil)
		return
	}

	// 计数
	count, err := vmFile.Count()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.CountFileFail, nil)
		return
	}

	// 填充数据
	data := map[string]interface{}{"lists": files, "count": count}

	controller.Response(c, http.StatusOK, errno.Success, data)
}

//根据文件Sign值获取单个文件的信息。
func GetFile(c *gin.Context) {

	valid := validation.Validation{}
	//获取file_sign
	file_sign := c.Param("file_sign")
	valid.Required(file_sign, "file_sign").Message("file_sign is required.")
	// 表单验证错误
	if valid.HasErrors() {
		controller.LogErrors(valid.Errors)
		controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
		return
	}
	// 构造结构体
	vmFile := vm.File{FileSign: file_sign}

	exist, err := vmFile.HasSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileFail, nil)
		return
	}
	if !exist {
		controller.Response(c, http.StatusNotFound, errno.FileIsNotExist, nil)
		return
	}

	file, err := vmFile.GetToSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileFail, nil)
		return
	}
	controller.Response(c, http.StatusOK, errno.Success, file)

}

//添加一个文件。(当前Node下)
func AddFile(c *gin.Context) {}

//编辑一个文件。(当前Node下)
func EditFile(c *gin.Context) {}

//更改文件可访问状态，状态(可用0/禁用1)，禁止非法文件访问,通常用在色情暴力凶杀等非法文件一键封杀
func DisabledFile(c *gin.Context) {

	valid := validation.Validation{}
	//获取file_sign
	file_sign := c.Param("file_sign")
	valid.Required(file_sign, "file_sign").Message("file_sign is required.")
	// 表单验证错误
	if valid.HasErrors() {
		controller.LogErrors(valid.Errors)
		controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
		return
	}
	// 构造结构体
	vmFile := vm.File{FileSign: file_sign}

	exist, err := vmFile.HasSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileFail, nil)
		return
	}
	if !exist {
		controller.Response(c, http.StatusNotFound, errno.FileIsNotExist, nil)
		return
	}

	//更改文件可访问状态
	err = vmFile.DisabledFileBySign()
	if err != nil {
		//fmt.Println(err.Error())
		controller.Response(c, http.StatusInternalServerError, errno.EditFileFail, nil)
		return
	}

	controller.Response(c, http.StatusOK, errno.Success, nil)

}

//标记删除一个文件
func DeleteMarkFile(c *gin.Context) {

	valid := validation.Validation{}
	//获取file_sign
	file_sign := c.Param("file_sign")
	valid.Required(file_sign, "file_sign").Message("file_sign is required.")
	// 表单验证错误
	if valid.HasErrors() {
		controller.LogErrors(valid.Errors)
		controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
		return
	}
	// 构造结构体
	vmFile := vm.File{FileSign: file_sign}

	exist, err := vmFile.HasSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileFail, nil)
		return
	}
	if !exist {
		controller.Response(c, http.StatusNotFound, errno.FileIsNotExist, nil)
		return
	}

	//将文件状态标记为删除状态
	err = vmFile.DeleteForFileSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.DeleteFileFail, nil)
		return
	}
	controller.Response(c, http.StatusOK, errno.Success, nil)

}

//彻底删除一个文件
func DeleteFile(c *gin.Context) {

	valid := validation.Validation{}
	//获取file_sign
	file_sign := c.Param("file_sign")
	valid.Required(file_sign, "file_sign").Message("file_sign is required.")
	// 表单验证错误
	if valid.HasErrors() {
		controller.LogErrors(valid.Errors)
		controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
		return
	}
	// 构造结构体
	vmFile := vm.File{FileSign: file_sign}

	//判断要修改的文件是否存在
	exist, err := vmFile.HasDelFileSign()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetFileFail, nil)
		return
	}
	if !exist {
		controller.Response(c, http.StatusNotFound, errno.FileIsNotExist, nil)
		return
	}
	//将文件状态标记为删除状态
	err = vmFile.DeleteFileBySign_Unscoped()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.DeleteFileFail, nil)
		return
	}
	controller.Response(c, http.StatusOK, errno.Success, nil)

}
