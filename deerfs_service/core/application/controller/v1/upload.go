package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/filter"
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

//上传块文件
func UploadChunksFile(c *gin.Context) {

	// 获取上传文件
	file, post_err := c.FormFile(config.UploadFormChunksField())
	if post_err != nil {
		loger.Lg.Error("upload chunks file error:", zap.String("error", post_err.Error()))
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

	//接收文件Key值
	file_sign, fs_exist := c.GetPostForm("file_sign")
	if !fs_exist {
		loger.Lg.Error("get file sign error")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "get file sign error")
		return
	}

	// //接收文件块索引值
	// file_index, fi_exist := c.GetQuery("file_index")
	// if !fi_exist {
	// 	loger.Lg.Error("get file index error")
	// 	controller.Response(c, http.StatusInternalServerError, errno.Error, "get file index error")
	// 	return
	// }

	//判断上传的文件片段后缀是否合规，必须是“.deerfs”
	//获取文件的后缀名
	extstring := path.Ext(file.Filename)
	if extstring != ".deerfs" {
		loger.Lg.Error("upload chunks filename extension is not '. deerfs' ")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "upload chunks filename extension is not '. deerfs' ")
		return
	}

	//创建一个存储文件夹
	saveFolder := path.Join(config.FileStorageDirectoryPath(), "temp", file_sign)
	createFolder_err := file_manage.CreateFolder(saveFolder) //创建存储文件夹
	if createFolder_err != nil {
		fmt.Println("Failed to create storage directory:", createFolder_err.Error())
		loger.Lg.Error("Failed to create storage directory:", zap.String("error", createFolder_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, createFolder_err.Error())
		return
	}

	//保存上传文件
	saveFile := path.Join(saveFolder, "/", file.Filename)
	saveErr := c.SaveUploadedFile(file, saveFile)
	if saveErr != nil {
		//保存失败
		loger.Lg.Info("save chunks file error:", zap.String("error", saveErr.Error())) //日志记录
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, saveErr.Error())
		return
	} else {
		//保存文件成功
		controller.FileResponse(c, http.StatusOK, file_sign, errno.Success, "ok")
		return
	}

}

//合并块文件
func MergeChunksFile(c *gin.Context) {

	//文件Key值
	file_sign, fs_exist := c.GetPostForm("file_sign")
	if !fs_exist {
		loger.Lg.Error("Failed to obtain the 'file_sign' value for the form.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "Failed to obtain the 'file_sign' value for the form.")
		return
	}
	//文件名
	file_name, fn_exist := c.GetPostForm("file_name")
	if !fn_exist {
		loger.Lg.Error("Failed to obtain the 'file_name' value for the form.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "Failed to obtain the 'file_name' value for the form.")
		return
	}
	//文件的块数
	file_chunk_number_str, fcn_exist := c.GetPostForm("file_chunk_number")
	if !fcn_exist {
		loger.Lg.Error("Failed to obtain the 'file_chunk_number' value for the form.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "Failed to obtain the 'file_chunk_number' value for the form.")
		return
	}
	//转换文件块数的值类型
	file_chunk_number, c_fcn_error := filter.Float(file_chunk_number_str)
	if c_fcn_error != nil {
		loger.Lg.Error("The 'file_chunk_number' value of the form is not a number.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "The 'file_chunk_number' value of the form is not a number.")
		return
	}
	//临时存储文件夹目录
	tempFolder := path.Join(config.FileStorageDirectoryPath(), "temp", file_sign)
	//判断文件夹下的文件数量是否与提交的值相等
	if file_chunk_number != float64(file_manage.ListDirFileNumber(tempFolder)) {
		loger.Lg.Error("The number of files in the temporary folder is not equal to the value of 'file_chunk_number'. Please maintain the integrity of the file data.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "The number of files in the temporary folder is not equal to the value of 'file_chunk_number'. Please maintain the integrity of the file data.")
		return
	}

	//截取file_sign值的前32位获取文件MD5值
	file_md5 := file_sign[0:32]
	//截取file_sign值32位md5后面的加密sign值
	sign_code := file_sign[32:]
	//解密sign_code
	deCrypt_sign_code, de_err := file_manage.DeCryptString(sign_code)
	if de_err != nil {
		//解密失败
		loger.Lg.Error("DeCrypt sign_code error")
		controller.Response(c, http.StatusInternalServerError, errno.Error, de_err.Error())
		return
	}
	//解密出sign_code值的信息
	file_de_info := strings.Split(string(deCrypt_sign_code), "-")
	//获取sign_code中文件的size
	file_de_info_size := file_de_info[0]
	filter_de_filesize, filter_file_size_err := filter.Int64(file_de_info_size)
	if filter_file_size_err != nil {
		fmt.Println("Failed to convert decrypted file bytes to int64")
		loger.Lg.Error("Failed to convert decrypted file bytes to int64")
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, "Failed to convert decrypted file bytes to int64")
		return
	}
	//获取sign_code中文件的type
	file_de_info_type := file_de_info[1]

	//根据上传的文件MD5查询数据库中所有数据这个文件是否存在
	// 构造结构体
	tempVmFile1 := vm.File{}
	file_exist_all, file_exist_all_err := tempVmFile1.HasHashFlexibleByAll(file_md5, file_de_info_type, filter_de_filesize)
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
		file_exist, file_exist_err := tempVmFile2.HasHashFlexible(global.Node_ID, file_md5, file_de_info_type, filter_de_filesize)
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

	//======合并文件======
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

	//获取文件的后缀名
	extstring := path.Ext(file_name)
	//拼接新的名字
	newfilename := common.JoinString(file_sign, extstring)

	//上传文件的保存路径
	saveFile := path.Join(saveFolder, "/", newfilename)

	//合并前先检查文件存在吗？存在就删除
	sfex, _ := file_manage.PathExists(saveFile)
	if sfex == true {
		os.RemoveAll(saveFile)
	}

	//merge_file
	deerfs_file, createfile_err := os.OpenFile(saveFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if createfile_err != nil {
		fmt.Println("Failed to create file", createfile_err.Error())
		loger.Lg.Error("Failed to create file:", zap.String("error", createfile_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, createfile_err.Error())
		return
	}

	//读取临时文件夹的数据，合并文件
	for i := 1; i <= int(file_chunk_number); i++ {
		f, ofs_err := os.OpenFile(path.Join(tempFolder, strconv.Itoa(int(i))+".deerfs"), os.O_RDONLY, os.ModePerm)
		if ofs_err != nil {
			fmt.Println("Failed to open chunk file", ofs_err.Error())
			loger.Lg.Error("Failed to open chunk file:", zap.String("error", ofs_err.Error()))
			controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, ofs_err.Error())
			return
		}
		b, rcf_err := ioutil.ReadAll(f)
		if rcf_err != nil {
			fmt.Println("Failed to read chunk file", rcf_err.Error())
			loger.Lg.Error("Failed to read chunk file:", zap.String("error", rcf_err.Error()))
			controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, rcf_err.Error())
			return
		}
		deerfs_file.Write(b)
		f.Close()
	}
	deerfs_file.Close()

	//文件合并成功后校验MD5文件
	check_filemd5, getmd5_err := file_manage.GetFileMD5ByPath(saveFile)
	if getmd5_err != nil {
		file_manage.DeleteFile(saveFile) //获取失败后删除文件
		fmt.Println("Calculate the MD5 value of the file", getmd5_err.Error())
		loger.Lg.Error("Calculate the MD5 value of the file:", zap.String("error", getmd5_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, getmd5_err.Error())
		return
	}
	if check_filemd5 != file_md5 {
		file_manage.DeleteFile(saveFile) //校验失败后删除文件
		fmt.Println("The merged file is inconsistent with the submitted file MD5")
		loger.Lg.Error("The merged file is inconsistent with the submitted file MD5")
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, "The merged file is inconsistent with the submitted file MD5")
		return
	}
	//文件合并成功后校验文件字节数
	check_filesize := file_manage.FileSize(saveFile)
	if filter_de_filesize != check_filesize {
		file_manage.DeleteFile(saveFile) //校验失败后删除文件
		fmt.Println("The bytes of the merged file and the submitted file do not match")
		loger.Lg.Error("The bytes of the merged file and the submitted file do not match")
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, "The bytes of the merged file and the submitted file do not match")
		return
	}
	//文件合并成功后校验文件类型
	newfile_type, file_mime := file_manage.CheckFileType(saveFile)
	if newfile_type != file_de_info_type {
		file_manage.DeleteFile(saveFile) //校验失败后删除文件
		fmt.Println("The type of the merged file does not match the submitted file")
		loger.Lg.Error("The type of the merged file does not match the submitted file")
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.Error, "The type of the merged file does not match the submitted file")
		return
	}
	//通过校验后上传数据库
	// 构造结构体
	currentTime := time.Now()

	vmFileAdd := vm.File{
		NodeId:    global.Node_ID,
		FileHash:  file_md5,
		FileSign:  file_sign,
		FileName:  file_name,
		FileExt:   extstring,
		FileType:  newfile_type,
		FileMime:  file_mime,
		FileSize:  check_filesize,
		FileAddr:  saveFile,
		CreatedBy: currentTime.Format("2006-01-02 15:04:05"), //form.CreatedBy
	}
	add_err := vmFileAdd.Add()
	if add_err != nil {
		file_manage.DeleteFile(saveFile) //添加数据库失败后删除文件
		loger.Lg.Error("[database]add file data error:", zap.String("error", add_err.Error()))
		controller.FileResponse(c, http.StatusInternalServerError, file_sign, errno.AddFileFail, add_err.Error())
		return
	}

	//合并文件成功后删除临时文件夹，节省资源
	os.RemoveAll(tempFolder)

	//日志记录
	file_json_data, _ := json.Marshal(vmFileAdd)
	loger.Lg.Info("save file success:", zap.String("success", string(file_json_data)))

	//将数据占用加入当前已使用字节数
	global.Directory_Storage_Use_Size = global.Directory_Storage_Use_Size + check_filesize

	//组合返回文件的访问路径
	return_path_str := path.Join(config.SystemTagUriAddress(), "/", file_sign)
	controller.FileResponse(c, http.StatusOK, file_sign, errno.Success, return_path_str)
	return

}

//清理临时块文件夹
func ClearChunksFile(c *gin.Context) {
	//文件Key值
	file_sign, fs_exist := c.GetPostForm("file_sign")
	if !fs_exist {
		loger.Lg.Error("Failed to obtain the 'file_sign' value for the form.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "Failed to obtain the 'file_sign' value for the form.")
		return
	}
	if len(file_sign) < 50 || len(file_sign) > 80 {
		loger.Lg.Error("The 'file_sign' value you submitted is not within a legal range.")
		controller.Response(c, http.StatusInternalServerError, errno.Error, "The 'file_sign' value you submitted is not within a legal range.")
		return
	}
	//临时存储文件夹目录
	tempFolder := path.Join(config.FileStorageDirectoryPath(), "temp", file_sign)
	//合并文件成功后删除临时文件夹，节省资源
	os.RemoveAll(tempFolder)
	controller.FileResponse(c, http.StatusOK, file_sign, errno.Success, tempFolder)
	return
}
