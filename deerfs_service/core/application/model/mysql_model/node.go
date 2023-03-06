package mysql_model

import (
	"github.com/jinzhu/gorm"
)

//Node节点模型
type Node struct {
	ID         int    `gorm:"primary_key" json:"id"`                                                     // 主键，根据约定不需要
	NodeName   string `gorm:"type:varchar(255);not null;default:'';COMMENT:'节点名';" json:"node_name"`     //节点名
	UriAddress string `gorm:"type:varchar(255);not null;default:'';COMMENT:'资源定位地址'" json:"uri_address"` //资源定位地址
	UseCap     int64  `gorm:"type:bigint;not null;default:'0';COMMENT:'单节点已使用容量'" json:"use_cap"`        //单节点已使用容量
	MaxCap     int64  `gorm:"type:bigint;not null;default:'0';COMMENT:'单节点最大使用容量'" json:"max_cap"`       //单节点最大使用容量
	Model             // 自定义的 Model
}

//获取Node的数据集合
func GetNodes(offset, limit int, cond map[string]interface{}) ([]*Node, error) {
	var nodes []*Node
	// 根据 cond 多表查询
	err := db.Where(cond).Offset(offset).Order("id desc").Limit(limit).Find(&nodes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return nodes, nil
}

//根据ID获取一个Node信息
func GetNodeToId(id int) (*Node, error) {
	var node Node
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&node).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &node, nil
}

//根据NodeName获取一个Node信息
func GetNodeToName(node_name string) (*Node, error) {
	var node Node
	err := db.Where("node_name = ? AND deleted_on = ? ", node_name, 0).First(&node).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &node, nil
}

//通过ID判断节点是否存在
func HasNodeByID(id int) (bool, error) {
	var node Node
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&node).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if node.ID > 0 {
		return true, nil
	}
	return false, nil
}

//通过节点名判断节点是否存在
func HasNodeByNodeName(node_name string) (bool, error) {
	var node Node
	err := db.Select("id").Where("node_name = ? AND deleted_on = ?", node_name, 0).First(&node).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if node.ID > 0 {
		return true, nil
	}
	return false, nil
}

//统计节点数据条数
func GetNodeCount(cond map[string]interface{}) (int, error) {
	var count int
	if err := db.Model(&Node{}).Where(cond).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//添加节点数据
func AddNodeToMap(data map[string]interface{}) error {
	// 根据 data 参数构造 node 结构体
	node := Node{
		NodeName:   data["node_name"].(string),
		UriAddress: data["uri_address"].(string),
		UseCap:     data["use_cap"].(int64),
		MaxCap:     data["max_cap"].(int64),
		Model: Model{
			CreatedBy: data["created_by"].(string),
		},
	}
	// 插入记录
	if err := db.Create(&node).Error; err != nil {
		return err
	}
	return nil
}

//添加节点数据
func AddNode(node_name, uri_address, createdBy string, use_cap, max_cap int64) error {
	// 根据参数构造 node 结构体
	node := Node{
		NodeName:   node_name,
		UriAddress: uri_address,
		UseCap:     use_cap,
		MaxCap:     max_cap,
		Model: Model{
			CreatedBy: createdBy,
		},
	}

	// 插入记录
	if err := db.Create(&node).Error; err != nil {
		return err
	}
	return nil
}

//修改节点数据,根据ID来修改
func EditNode(id int, data map[string]interface{}) error {
	if err := db.Model(&Node{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//修改节点数据,根据节点名修改
func EditNodeForNodeName(node_name string, data map[string]interface{}) error {
	if err := db.Model(&Node{}).Where("node_name = ? AND deleted_on = ?", node_name, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//将节点状态标记为删除状态
func DeleteNodeChangeMark(id int) error {
	if err := db.Model(&Node{}).Where("id = ? AND deleted_on = ?", id, 0).Update("deleted_on", 1).Error; err != nil {
		return err
	}
	return nil
}

//将节点状态标记为删除状态,条件是根据节点名
func DeleteNodeChangeMarkForNodeName(node_name string) error {
	if err := db.Model(&Node{}).Where("node_name  = ? AND deleted_on = ?", node_name, 0).Update("deleted_on", 1).Error; err != nil {
		return err
	}
	return nil
}

//永久删除节点数据
func DeleteNode(id int) error {
	if err := db.Where("id = ?", id).Delete(&Node{}).Error; err != nil {
		return err
	}
	return nil
}

//永久删除节点数据,条件是根据节点名
func DeleteNodeForNodeName(node_name string) error {
	if err := db.Where("node_name = ?", node_name).Delete(&Node{}).Error; err != nil {
		return err
	}
	return nil
}

//永久删除过期数据
func DeleteNodes() error {
	// Unscoped 返回所有记录，包含软删除的记录
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Node{}).Error; err != nil {
		return err
	}
	return nil
}
