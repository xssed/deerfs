package application

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
)

//配置初始化时,配置节点ID,并向数据库中更新当前的数据信息
func Init_Deerfs_Info() {

	//如果没有创建文件存储目录则创建它
	err := file_manage.CreateFolder(config.FileStorageDirectoryPath())
	if err != nil {
		//创建文件夹失败
		fmt.Println(err.Error())
		loger.Lg.Fatal(err.Error()) // 记录日志
	}
	//启动服务时获取当前文件存储目录的字节数使用情况
	fmt.Println("Checking the used size of the directory......")
	file_manage.ListDirSize(config.FileStorageDirectoryPath())

	//系统最大使用容量
	maxcap := common.JoinString("The max storage cap set in the deerfs system is ", file_manage.FormatFileSizeToString(config.FileStorageDirectoryStorageMaxSize()), ".")
	fmt.Println(maxcap)
	loger.Lg.Info(maxcap) // 记录日志

	//系统当前已使用容量
	usecap := common.JoinString(file_manage.FormatFileSizeToString(global.Directory_Storage_Use_Size), " files have been stored in deerfs system.")
	fmt.Println(usecap)
	loger.Lg.Info(usecap) // 记录日志

	//对已经存储文件的数量的统计
	stored_file_number := common.JoinString("The number of files stored in deerfs system is ", common.Int64ToString(global.Directory_Storage_File_Number), ".")
	fmt.Println(stored_file_number)
	loger.Lg.Info(stored_file_number) // 记录日志

	//如果当前使用容量大于设定的最大存储容量deerfs系统则异常退出
	if global.Directory_Storage_Use_Size > config.FileStorageDirectoryStorageMaxSize() {
		fmt.Println("The current used cap is greater than the max storage cap set by the deerfs system.")
		loger.Lg.Fatal("The current used cap is greater than the max storage cap set by the deerfs system.") // 记录日志
	}

	//获取节点名字
	nodename := config.SystemTagName()
	if len(nodename) <= 1 {
		fmt.Println("<name> in the [system_tag] group in the config file has no config parameters.")
		loger.Lg.Fatal("<name> in the [system_tag] group in the config file has no config parameters.") // 记录日志
	}
	//通过查询节点名获取节点数据来判断节点是否存在
	node, err := mysql_model.GetNodeToName(nodename)
	if err != nil {
		fmt.Println("Error executing Init_Global_Deerfs_Info() function.")
		loger.Lg.Fatal("Error executing Init_Global_Deerfs_Info() function.") // 记录日志
	}
	if node.ID < 1 {
		//不存在则新增
		//构建新增节点数据
		add_data := map[string]interface{}{}
		add_data["node_name"] = config.SystemTagName()
		add_data["uri_address"] = config.SystemTagUriAddress()
		add_data["use_cap"] = global.Directory_Storage_Use_Size
		add_data["max_cap"] = config.FileStorageDirectoryStorageMaxSize()

		currentTime := time.Now()
		add_data["created_by"] = currentTime.Format("2006-01-02 15:04:05")
		add_err := mysql_model.AddNodeToMap(add_data)
		if add_err == nil {
			//通过查询节点名获取节点数据来判断节点是否存在
			new_node, _ := mysql_model.GetNodeToName(config.SystemTagName())
			//赋值全局
			global.Node_ID = new_node.ID
			fmt.Println("Write node data to database success.")
			loger.Lg.Info("Write node data to database success.") // 记录日志
		} else {
			fmt.Println("Failed to write node data to database.")
			loger.Lg.Fatal("Failed to write node data to database.")
		}

	} else {
		//存在则更新
		//构建更新节点数据
		up_data := map[string]interface{}{}
		up_data["node_name"] = config.SystemTagName()
		up_data["uri_address"] = config.SystemTagUriAddress()
		up_data["use_cap"] = global.Directory_Storage_Use_Size
		up_data["max_cap"] = config.FileStorageDirectoryStorageMaxSize()

		currentTime := time.Now()
		up_data["modified_by"] = currentTime.Format("2006-01-02 15:04:05")

		err = mysql_model.EditNode(node.ID, up_data)
		if err == nil {
			//赋值全局
			global.Node_ID = node.ID
			fmt.Println("Update node data to database success.")
			loger.Lg.Info("Update node data to database success.") // 记录日志
		} else {
			fmt.Println("Failed to update node data to database.")
			loger.Lg.Fatal("Failed to update node data to database.")
		}
	}
	//输出全局节点ID
	nodeid_str := common.JoinString("Global Node_ID :", strconv.Itoa(global.Node_ID))
	fmt.Println(nodeid_str)
	loger.Lg.Info(nodeid_str) // 记录日志
	//获取节点ID,进行最后的校验
	global.Node, _ = mysql_model.GetNodeToId(global.Node_ID)
	if global.Node.ID < 1 {
		fmt.Println("The deerfs system failed to initialize and set node_id, and exited.")
		loger.Lg.Fatal("The deerfs system failed to initialize and set node_id, and exited.") // 记录日志
	}
}
