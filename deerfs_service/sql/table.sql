-- 创建文件表
CREATE TABLE `deerfs_file` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `node_id` int(11) DEFAULT '0' COMMENT 'node id',
  `file_hash` char(32) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_sign` char(120) NOT NULL DEFAULT '' COMMENT '文件标识',
  `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '上传时的文件名',
  `file_ext` char(10) NOT NULL DEFAULT '' COMMENT '上传时的文件扩展名',
  `file_type` char(10) NOT NULL DEFAULT '' COMMENT '文件类型(后台识别)',
  `file_mime` char(57) NOT NULL DEFAULT '' COMMENT '文件MIME(后台识别)',
  `file_size` bigint DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(555) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  /* `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期', default '0000-00-00 00:00:00' */
  `disabled` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用0/禁用1)，禁止非法文件访问',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  `created_on` int NOT NULL  DEFAULT '0',
  `modified_on` int NOT NULL  DEFAULT '0',
  `deleted_on` int NOT NULL  DEFAULT '0',
  `created_by` varchar(255) NOT NULL  DEFAULT '',
  `modified_by` varchar(255) NOT NULL  DEFAULT '' , 
  PRIMARY KEY (`id`),
  KEY `index_node_id` (`node_id`),
  KEY `index_file_hash` (`file_hash`),
  UNIQUE  KEY `index_file_sign` (`file_sign`),
  KEY `index_file_ext` (`file_ext`),
  KEY `index_file_type` (`file_type`),
  KEY `index_file_mime` (`file_mime`),
  KEY `index_disabled` (`disabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建node表
CREATE TABLE `deerfs_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `node_name` varchar(255) NOT NULL DEFAULT '' COMMENT '节点名',
  `uri_address` varchar(255) NOT NULL DEFAULT '' COMMENT '资源定位地址',
  `use_cap` bigint DEFAULT '0' COMMENT '单节点已使用容量',
  `max_cap` bigint DEFAULT '0' COMMENT '单节点最大使用容量',
  `created_on` int NOT NULL  DEFAULT '0',
  `modified_on` int NOT NULL  DEFAULT '0',
  `deleted_on` int NOT NULL  DEFAULT '0',
  `created_by` varchar(255) NOT NULL  DEFAULT '',
  `modified_by` varchar(255) NOT NULL  DEFAULT '' , 
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_node_name` (`node_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;