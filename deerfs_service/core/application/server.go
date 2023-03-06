package application

import (
	"fmt"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/application/model/owlcache_model"
	"github.com/xssed/deerfs/deerfs_service/core/system/cachefile_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
)

//开启deerfs服务
func StartServer() {

	//执行步骤信息
	fmt.Println("deerfs system http server initialization...")

	//初始化连接owlcache的http服务检测其状态并获取可用token
	owlcache_model.Init()
	//初始化Mysql数据库连接数据
	mysql_model.Conn()
	//配置初始化时,配置节点ID,并向数据库中更新当前的数据信息
	Init_Deerfs_Info()
	//缓存文件管理启动
	go cachefile_manage.Cachefile_manage_start()
	//开启THHP服务(使用的GIN框架)
	StartRoutes().Run(config.HttpAddr() + ":" + config.HttpPort())

}
