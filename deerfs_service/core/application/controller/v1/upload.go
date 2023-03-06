package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/vm"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

func UploadFile(c *gin.Context) {

	// 获取上传文件
	file, post_err := c.FormFile(config.UploadFormField())
	if post_err != nil {
		loger.Lg.Error("upload file error:", zap.String("error", post_err.Error()))
		controller.Response(c, http.StatusInternalServerError, errno.UploadFileFail, post_err.Error())
		return
	}

	//校验判断上传的文件大小是否大于当前节点设定的单文件最大值
	if file.Size > config.FileStorageFileStorageMaxSize() {
		return_str := common.JoinString("The size of the current upload file:", file_manage.FormatFileSizeToString(file.Size), ".  The max size of a single file set by the current node:", file_manage.FormatFileSizeToString(config.FileStorageFileStorageMaxSize()), ".")
		loger.Lg.Error("upload file large size:", zap.String("error", errno.Msg[errno.UploadFileLargeSizeFail]))
		controller.Response(c, http.StatusInternalServerError, errno.UploadFileLargeSizeFail, return_str)
		return
	}

	//计算文件MD5值，然后开展工作
	file_md5, md5_err := file_manage.GetFileMD5(file)
	if md5_err != nil {
		loger.Lg.Error("get file md5 error:", zap.String("error", md5_err.Error()))
		controller.Response(c, http.StatusInternalServerError, errno.Error, md5_err.Error())
		return
	}
	//判断文件类型,返回文件类型字符串和MIME信息字符串
	file_type, file_mime := file_manage.CheckFileTypeByMultipart(file)

	//根据上传的文件MD5查询数据库中所有数据这个文件是否存在
	// 构造结构体
	tempVmFile1 := vm.File{}
	file_exist_all, file_exist_all_err := tempVmFile1.HasHashFlexibleByAll(file_md5, file_type, file.Size)
	if file_exist_all_err != nil {
		loger.Lg.Error("Failed to find file data in all nodes:", zap.String("error", file_exist_all_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_md5, errno.GetFileFail, file_exist_all_err.Error())
		return
	}
	//如果查找到集群中存在这个文件
	if file_exist_all {

		//判断当前节点是否允许冗余存储
		if config.SystemTagAllowDuplicates() == 0 {
			loger.Lg.Error("upload file repeat:", zap.String("error", errno.Msg[errno.UploadFileRepeatAllFail]))
			controller.FileResponse(c, http.StatusInternalServerError, file_md5, errno.UploadFileRepeatAllFail, nil)
			return
		}

		//再查询是否存在于当前节点
		// 构造结构体
		tempVmFile2 := vm.File{}
		//根据上传的文件MD5查询数据库中所有数据这个文件是否存在
		file_exist, file_exist_err := tempVmFile2.HasHashFlexible(global.Node_ID, file_md5, file_type, file.Size)
		if file_exist_err != nil {
			loger.Lg.Error("Failed to find file data in current node:", zap.String("error", file_exist_err.Error()))
			controller.FileResponse(c, http.StatusInternalServerError, file_md5, errno.GetFileFail, file_exist_err.Error())
			return
		}
		//如果查找到集群中存在这个文件
		if file_exist {
			loger.Lg.Error("upload file repeat:", zap.String("error", errno.Msg[errno.UploadFileRepeatFail]))
			controller.FileResponse(c, http.StatusInternalServerError, file_md5, errno.UploadFileRepeatFail, nil)
			return
		}
	}

	//创建一个存储文件夹
	saveFilePathString := file_manage.GetFilePathString() //获取存储路径
	saveFolder := path.Join(config.FileStorageDirectoryPath(), saveFilePathString)
	createFolder_err := file_manage.CreateFolder(saveFolder) //创建存储文件夹
	if createFolder_err != nil {
		fmt.Println("Failed to create storage directory:", createFolder_err.Error())
		loger.Lg.Error("Failed to create storage directory:", zap.String("error", createFolder_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_md5, errno.Error, createFolder_err.Error())
		return
	}

	//生成新的文件sign
	//sign格式为:文件MD5+加密(文件字节数+"-"+文件类型+"-"+四位随机数)
	str_filesize := strconv.FormatInt(file.Size, 10) //int64->string
	enCrypt_str := file_manage.EnCryptToString(common.JoinString(str_filesize, "-", file_type, "-", common.RandAllString(4)))
	file_sign := common.JoinString(file_md5, enCrypt_str) //计算sign
	//获取文件的后缀名
	extstring := path.Ext(file.Filename)
	//拼接新的名字
	newfilename := common.JoinString(file_sign, extstring)

	//保存上传文件
	saveFile := path.Join(saveFolder, "/", newfilename)
	saveErr := c.SaveUploadedFile(file, saveFile)
	if saveErr == nil {
		//保存文件成功
		// 构造结构体
		currentTime := time.Now()

		vmFileAdd := vm.File{
			NodeId:    global.Node_ID,
			FileHash:  file_md5,
			FileSign:  file_sign,
			FileName:  file.Filename,
			FileExt:   extstring,
			FileType:  file_type,
			FileMime:  file_mime,
			FileSize:  file.Size,
			FileAddr:  saveFile,
			CreatedBy: currentTime.Format("2006-01-02 15:04:05"), //form.CreatedBy
		}
		add_err := vmFileAdd.Add()
		if add_err != nil {
			loger.Lg.Error("add file data error:", zap.String("error", add_err.Error()))
			controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.AddFileFail, add_err.Error())
			return
		}
		//日志记录
		file_json_data, _ := json.Marshal(vmFileAdd)
		loger.Lg.Info("save file success:", zap.String("success", string(file_json_data)))

		//将数据占用加入当前已使用字节数
		global.Directory_Storage_Use_Size = global.Directory_Storage_Use_Size + file.Size

		//组合返回文件的访问路径
		return_path_str := path.Join(config.SystemTagUriAddress(), "/", file_sign)
		controller.FileResponse(c, http.StatusOK, file_sign, errno.Success, return_path_str)
		return
	} else {
		//保存文件失败
		loger.Lg.Error("save file error:", zap.String("error", saveErr.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, saveErr.Error())
		return
	}

}
