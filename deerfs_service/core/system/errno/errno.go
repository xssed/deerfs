package errno

// 自定义的错误码
const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	NodeNameIsExisted  = 10001
	GetExistedNodeFail = 10002
	NodeIsNotExist     = 10003
	GetAllNodeFail     = 10004
	CountNodeFail      = 10005
	AddNodeFail        = 10006
	EditNodeFail       = 10007
	DeleteNodeFail     = 10008
	NodeNameRepeat     = 10009

	FileIsNotExist          = 20001
	CheckFileIsExistFail    = 20002
	AddFileFail             = 20003
	DeleteFileFail          = 20004
	EditFileFail            = 20005
	CountFileFail           = 20006
	GetFileListFail         = 20007
	GetFileFail             = 20008
	UploadFileFail          = 20009
	UploadFileLargeSizeFail = 20010
	UploadFileRepeatFail    = 20011
	UploadFileRepeatAllFail = 20012

	OwlcacheOffline         = 30001
	OwlcacheSetFileInfoFail = 30002
)

// 错误码对应的错误消息
var Msg = map[int]string{
	Success:       "Success",        //成功
	Error:         "Error",          //错误
	InvalidParams: "Invalid Params", //请求参数错误

	NodeNameIsExisted:  "The node name already exist",             //已存在该节点名称
	GetExistedNodeFail: "Failed to get the exist node",            //获取已存在节点失败
	NodeIsNotExist:     "The node does not exist",                 //该节点不存在
	GetAllNodeFail:     "Failed to get all nodes",                 //获取所有节点失败
	CountNodeFail:      "Count node failed",                       //统计节点失败
	AddNodeFail:        "Failed to add node",                      //新增节点失败
	EditNodeFail:       "Failed to modify node",                   //修改节点失败
	DeleteNodeFail:     "Failed to delete node",                   //删除节点失败
	NodeNameRepeat:     "The modify new node name already exists", //修改的新节点名称已存在

	FileIsNotExist:          "The file does not exist",                                                                                                                                                                  //该文件不存在
	CheckFileIsExistFail:    "Failed to check whether the file exists",                                                                                                                                                  //检查文件是否存在失败
	AddFileFail:             "Failed to add file",                                                                                                                                                                       //新增文件失败
	DeleteFileFail:          "Failed to delete file",                                                                                                                                                                    //删除文件失败
	EditFileFail:            "Failed to modify file",                                                                                                                                                                    //修改文件失败
	CountFileFail:           "Count file failed",                                                                                                                                                                        //统计文件失败
	GetFileListFail:         "Failed to get multiple files",                                                                                                                                                             //获取多个文件失败
	GetFileFail:             "Failed to get a single file",                                                                                                                                                              //获取单个文件失败
	UploadFileFail:          "Upload file failed",                                                                                                                                                                       //上传文件失败
	UploadFileLargeSizeFail: "The uploaded file size exceeds the set file size, please use chunks upload.If you use chunks upload, check that the file chunks size exceeds the maximum single file size on the server.", //上传文件大小超过设定文件大小,请使用分块上传
	UploadFileRepeatFail:    "This file already exists in the current node, no need to upload it again",                                                                                                                 //当前节点已经存在此文件，无需再次上传
	UploadFileRepeatAllFail: "This file already exists in node clusters. The current node config does not allow file redundancy, so there is no need to upload it again",                                                //该文件已存在于节点集群中。 当前节点配置不允许文件冗余，因此无需再次上传

	OwlcacheOffline:         "The status of owlcache is offline",   //owlcache为宕机状态
	OwlcacheSetFileInfoFail: "Failed to set file info in Owlcache", //向owlcache中设置文件信息失败
}
