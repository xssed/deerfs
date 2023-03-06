package application

import (
	"github.com/gin-gonic/gin"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/v1"
	"github.com/xssed/deerfs/deerfs_service/core/application/middleware"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
)

func StartRoutes() *gin.Engine {

	gin.SetMode(config.HttpMode())

	r := gin.Default()
	// 注册zap相关中间件
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	//过滤无效请求
	r.GET("/favicon.ico", func(c *gin.Context) { return })

	//数据库数据处理
	apiv1 := r.Group("/deerfs_db/v1")
	{
		//System 系统基本信息
		apiv1.GET("/system", v1.GetSystemInfo)
		apiv1.GET("/tools", v1.Tools)

		//File 文件路由
		apiv1.GET("/files", v1.GetFiles)          //获取文件信息集合
		apiv1.GET("/file/:file_sign", v1.GetFile) //获取单个文件的信息
		// apiv1.POST("/file", v1.AddFile)                    //添加一个文件信息
		// apiv1.PUT("/file/:file_sign", v1.EditFile)               //编辑一个文件信息
		apiv1.PUT("/file/disabled/:file_sign", v1.DisabledFile)  //禁用一个文件，通常是违禁文件
		apiv1.DELETE("/file/mark/:file_sign", v1.DeleteMarkFile) //删除一个文件信息,隐式删除
		apiv1.DELETE("/file/:file_sign", v1.DeleteFile)          //删除一个文件信息,彻底删除

		//Node 节点路由
		apiv1.GET("/nodes", v1.GetNodes)                     //获取节点信息集合
		apiv1.GET("/node/:param", v1.GetNode)                //获取单个节点的信息。param可以是ID,也可以是NodeName。
		apiv1.POST("/node", v1.AddNode)                      //添加一个节点信息
		apiv1.PUT("/node/:param", v1.EditNode)               //编辑一个节点信息。param可以是ID,也可以是NodeName。
		apiv1.DELETE("/node/mark/:param", v1.DeleteMarkNode) //删除一个节点信息,隐式删除。param可以是ID,也可以是NodeName。
		apiv1.DELETE("/node/:param", v1.DeleteNode)          //删除一个节点信息,彻底删除。param可以是ID,也可以是NodeName。
	}

	//上传文件部分
	upload_v1 := r.Group("/deerfs_upload")
	{
		upload_v1.POST("/upload", v1.UploadFile)
	}

	//文件访问
	r.GET("/:file_sign", v1.DownloadFile)
	//欢迎页
	r.GET("/", v1.SayHello)

	return r

}
