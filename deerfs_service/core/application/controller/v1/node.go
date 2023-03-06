package v1

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/vm"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
)

//获取节点信息集合
func GetNodes(c *gin.Context) {

	// 构造结构体
	vmNode := vm.Node{PageNum: controller.PageNum(c), PageSize: controller.PageSize} //ID, NodeName, UriAddress, UseCap, MaxCap, CreatedBy, ModifiedBy,

	nodes, err := vmNode.GetAll()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetAllNodeFail, nil)
		return
	}

	// 计数
	count, err := vmNode.Count()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.CountNodeFail, nil)
		return
	}

	// 填充数据
	data := map[string]interface{}{"lists": nodes, "count": count}

	controller.Response(c, http.StatusOK, errno.Success, data)
}

//获取单个节点的信息。param可以是ID,也可以是NodeName。
func GetNode(c *gin.Context) {
	// 获取param
	param := c.Param("param")
	//判断param是ID还是NodeName
	if common.IsNumber(param) {
		id := com.StrTo(param).MustInt()
		valid := validation.Validation{}
		valid.Min(id, 1, "id")

		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{ID: id}
		exist, err := vmNode.HasID()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}

		node, err := vmNode.GetToId()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		controller.Response(c, http.StatusOK, errno.Success, node)

	} else {
		// 构造结构体
		vmNode := vm.Node{NodeName: param}
		exist, err := vmNode.HasName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}

		node, err := vmNode.GetToName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		controller.Response(c, http.StatusOK, errno.Success, node)
	}

}

//添加一个节点信息
func AddNode(c *gin.Context) {

	var form AddNodeForm
	httpCode, errCode := controller.BindAndValid(c, &form)
	if errCode != errno.Success {
		controller.Response(c, httpCode, errCode, nil)
		return
	}
	currentTime := time.Now()

	vmNode := vm.Node{
		NodeName:   form.NodeName,
		UriAddress: form.UriAddress,
		UseCap:     com.StrTo(form.UseCap).MustInt64(),
		MaxCap:     com.StrTo(form.MaxCap).MustInt64(),
		CreatedBy:  currentTime.Format("2006-01-02 15:04:05"), //form.CreatedBy
	}
	exist, err := vmNode.HasName()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
		return
	}
	if exist {
		controller.Response(c, http.StatusInternalServerError, errno.NodeNameIsExisted, nil)
		return
	}
	err = vmNode.Add()
	if err != nil {
		controller.Response(c, http.StatusInternalServerError, errno.AddNodeFail, nil)
		return
	}

	controller.Response(c, http.StatusOK, errno.Success, nil)

}

//编辑一个节点信息。param可以是ID,也可以是NodeName。
func EditNode(c *gin.Context) {

	var form EditNodeForm
	httpCode, errCode := controller.BindAndValid(c, &form)
	if errCode != errno.Success {
		controller.Response(c, httpCode, errCode, nil)
		return
	}
	currentTime := time.Now()

	valid := validation.Validation{}

	// 获取param
	param := c.Param("param")
	//判断param是ID还是NodeName
	if common.IsNumber(param) {

		id := com.StrTo(param).MustInt()
		valid.Min(id, 1, "id")

		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			ID:         id,
			ModifiedBy: currentTime.Format("2006-01-02 15:04:05"), //form.ModifiedBy
			NodeName:   form.NodeName,
			UriAddress: form.UriAddress,
			UseCap:     com.StrTo(form.UseCap).MustInt64(),
			MaxCap:     com.StrTo(form.MaxCap).MustInt64(),
		}
		//根据ID判断节点内容是否存在
		old_exist, old_err := vmNode.HasID()
		if old_err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !old_exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}
		//判断新修改的节点名字是否存在
		exist, err := vmNode.HasNameFlexible(form.NodeName)
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if exist {
			controller.Response(c, http.StatusInternalServerError, errno.NodeNameRepeat, nil)
			return
		}
		//编辑节点
		err = vmNode.Edit() //根据ID来编辑
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.EditNodeFail, nil)
			return
		}

	} else {
		//识别为节点名后，通过节点名来更新数据
		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			ModifiedBy: currentTime.Format("2006-01-02 15:04:05"), //form.ModifiedBy
			NodeName:   form.NodeName,
			UriAddress: form.UriAddress,
			UseCap:     com.StrTo(form.UseCap).MustInt64(),
			MaxCap:     com.StrTo(form.MaxCap).MustInt64(),
		}
		//判断节点内容是否存在
		old_exist, old_err := vmNode.HasNameFlexible(param)
		if old_err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !old_exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}
		//判断要修改的节点名是否存在
		exist, err := vmNode.HasName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if exist {
			controller.Response(c, http.StatusInternalServerError, errno.NodeNameRepeat, nil)
			return
		}
		//编辑节点
		err = vmNode.EditForNodeName(param) //根据节点名字来编辑
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.EditNodeFail, nil)
			return
		}

	}

	controller.Response(c, http.StatusOK, errno.Success, nil)
}

//标记删除一个节点信息。param可以是ID,也可以是NodeName。
func DeleteMarkNode(c *gin.Context) {

	valid := validation.Validation{}

	// 获取param
	param := c.Param("param")
	//判断param是ID还是NodeName
	if common.IsNumber(param) {

		id := com.StrTo(param).MustInt()
		valid.Min(id, 1, "id").Message("must be greater than 0.")

		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			ID: id,
		}
		//根据ID判断节点内容是否存在
		exist, err := vmNode.HasID()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}
		//将节点状态标记为删除状态
		err = vmNode.DeleteNodeChangeMark()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.EditNodeFail, nil)
			return
		}

	} else {
		//识别为节点名后，通过节点名来更新数据
		valid.MinSize(param, 2, "node_name").Message("length must be greater than 2.")
		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			NodeName: param, //form.NodeName
		}
		//判断要修改的节点名是否存在
		exist, err := vmNode.HasName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusInternalServerError, errno.NodeIsNotExist, nil)
			return
		}
		//将节点状态标记为删除状态
		err = vmNode.DeleteNodeChangeMarkForNodeName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.EditNodeFail, nil)
			return
		}

	}

	controller.Response(c, http.StatusOK, errno.Success, nil)

}

func DeleteNode(c *gin.Context) {

	valid := validation.Validation{}

	// 获取param
	param := c.Param("param")
	//判断param是ID还是NodeName
	if common.IsNumber(param) {

		id := com.StrTo(param).MustInt()
		valid.Min(id, 1, "id").Message("must be greater than 0.")

		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			ID: id,
		}
		//根据ID判断节点内容是否存在
		exist, err := vmNode.HasID()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusNotFound, errno.NodeIsNotExist, nil)
			return
		}
		//将节点删除
		err = vmNode.Delete()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.DeleteNodeFail, nil)
			return
		}

	} else {
		//识别为节点名后，通过节点名来删除数据
		valid.MinSize(param, 2, "node_name").Message("length must be greater than 2.")
		// 表单验证错误
		if valid.HasErrors() {
			controller.LogErrors(valid.Errors)
			controller.Response(c, http.StatusBadRequest, errno.InvalidParams, nil)
			return
		}
		// 构造结构体
		vmNode := vm.Node{
			NodeName: param, //form.NodeName
		}
		//判断要修改的节点名是否存在
		exist, err := vmNode.HasName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.GetExistedNodeFail, nil)
			return
		}
		if !exist {
			controller.Response(c, http.StatusInternalServerError, errno.NodeIsNotExist, nil)
			return
		}
		//将节点删除
		err = vmNode.DeleteForNodeName()
		if err != nil {
			controller.Response(c, http.StatusInternalServerError, errno.DeleteNodeFail, nil)
			return
		}

	}

	controller.Response(c, http.StatusOK, errno.Success, nil)

}
