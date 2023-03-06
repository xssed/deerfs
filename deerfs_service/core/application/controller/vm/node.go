package vm

import (
	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
)

//Node
type Node struct {
	ID         int
	NodeName   string
	UriAddress string
	UseCap     int64
	MaxCap     int64
	CreatedBy  string
	ModifiedBy string
	PageNum    int
	PageSize   int
}

//通过节点名判断节点是否存在
func (n *Node) HasName() (bool, error) {
	return mysql_model.HasNodeByNodeName(n.NodeName)
}

//通过节点名判断节点是否存在,动态传值
func (n *Node) HasNameFlexible(node_name string) (bool, error) {
	return mysql_model.HasNodeByNodeName(node_name)
}

//通过ID判断节点是否存在
func (n *Node) HasID() (bool, error) {
	return mysql_model.HasNodeByID(n.ID)
}

//添加节点数据
func (n *Node) Add() error {
	return mysql_model.AddNode(n.NodeName, n.UriAddress, n.CreatedBy, n.UseCap, n.MaxCap)
}

//修改节点数据，根据ID
func (n *Node) Edit() error {
	data := map[string]interface{}{
		"modified_by": n.ModifiedBy,
		"node_name":   n.NodeName,
		"uri_address": n.UriAddress,
		"use_cap":     n.UseCap,
		"max_cap":     n.MaxCap,
	}
	return mysql_model.EditNode(n.ID, data)
}

//修改节点数据，根据节点名字
func (n *Node) EditForNodeName(node_name string) error {
	data := map[string]interface{}{
		"modified_by": n.ModifiedBy,
		"node_name":   n.NodeName,
		"uri_address": n.UriAddress,
		"use_cap":     n.UseCap,
		"max_cap":     n.MaxCap,
	}
	return mysql_model.EditNodeForNodeName(node_name, data)
}

//删除节点数据
func (n *Node) Delete() error {
	return mysql_model.DeleteNode(n.ID)
}

//删除节点数据,条件是根据节点名
func (n *Node) DeleteForNodeName() error {
	return mysql_model.DeleteNodeForNodeName(n.NodeName)
}

//将节点状态标记为删除状态
func (n *Node) DeleteNodeChangeMark() error {
	return mysql_model.DeleteNodeChangeMark(n.ID)
}

//将节点状态标记为删除状态,条件是根据节点名
func (n *Node) DeleteNodeChangeMarkForNodeName() error {
	return mysql_model.DeleteNodeChangeMarkForNodeName(n.NodeName)
}

//统计数据
func (n *Node) Count() (int, error) {
	return mysql_model.GetNodeCount(n.toMap())
}

//获取所有数据
func (n *Node) GetAll() ([]*mysql_model.Node, error) {
	var nodes []*mysql_model.Node
	nodes, err := mysql_model.GetNodes(n.PageNum, n.PageSize, n.toMap())
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

//获取单条数据
func (n *Node) GetToId() (*mysql_model.Node, error) {
	var node *mysql_model.Node
	node, err := mysql_model.GetNodeToId(n.ID)
	if err != nil {
		return nil, err
	}
	return node, nil
}

//获取单条数据
func (n *Node) GetToName() (*mysql_model.Node, error) {
	var node *mysql_model.Node
	node, err := mysql_model.GetNodeToName(n.NodeName)
	if err != nil {
		return nil, err
	}
	return node, nil
}

//查询条件
func (n *Node) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["deleted_on"] = 0
	if n.NodeName != "" {
		m["node_name"] = n.NodeName
	}
	if n.UriAddress != "" {
		m["uri_address"] = n.UriAddress
	}
	return m
}
