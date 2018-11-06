CREATE TABLE `visual_dispatch` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '调度任务id',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `description` varchar(512) NOT NULL DEFAULT '' COMMENT '描述',
  `crontab` varchar(255) NOT NULL DEFAULT '' COMMENT '定时',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1-未执行，2-执行中，3-执行结束，4-执行失败，5-待执行',
  `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
  `node` text NOT NULL COMMENT '节点',
  `relation` text NOT NULL COMMENT '关系',
  `create_uid` int(11) NOT NULL COMMENT '添加者',
  `update_uid` int(11) NOT NULL COMMENT '修改者',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0-未删除  1-已删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_unique_title` (`title`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='调度服务-调度任务信息主表';

CREATE TABLE `visual_dispatch_callback` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `dispatch_id` int(11) NOT NULL DEFAULT '0' COMMENT '调度任务id 关联 visual_dipatch.id',
  `node_id` int(11) NOT NULL DEFAULT '0' COMMENT '调度任务中每个单元执行块的排序id',
  `logic_id` int(11) NOT NULL DEFAULT '0' COMMENT '调度任务每个单元块的具体逻辑id关联data_sync_logic_configuration.id',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务状态：1 执行中，2 执行结束，3 执行失败',
  `progress` tinyint(4) NOT NULL DEFAULT '0' COMMENT '执行进度 整数',
  `remark` text NOT NULL COMMENT '备注信息',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `index_dispatch_date` (`dispatch_id`,`create_time`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='调度服务-调度任务体执行状况回调表';

CREATE TABLE `visual_dispatch_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `title` varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
  `origin` text COMMENT '源',
  `target` text COMMENT '目标',
  `params` text COMMENT '参数',
  `where_condition` text COMMENT 'where条件',
  `remark` text COMMENT '备注',
  `create_uid` int(10) NOT NULL DEFAULT '0' COMMENT '创建用户ID',
  `update_uid` int(10) NOT NULL DEFAULT '0' COMMENT '修改用户ID',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0：正常 1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `inx_title` (`title`),
  KEY `inx_create_uid` (`create_uid`),
  KEY `inx_update_uid` (`update_uid`),
  KEY `inx_is_deleted` (`is_deleted`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='数据同步业务配置';

CREATE TABLE `visual_dispatch_data` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id无意义',
  `dispatch_id` int(11) NOT NULL DEFAULT '0' COMMENT '调度任务id 关联 visual_dipatch.id',
  `node_id` tinyint(4) NOT NULL DEFAULT '1' COMMENT '调度任务中每个单元执行块的排序id',
  `logic_id` int(11) NOT NULL DEFAULT '0' COMMENT '调度任务每个单元块的具体逻辑id关联data_sync_logic_configuration.id',
  `pre_node` varchar(255) NOT NULL DEFAULT '0' COMMENT '调度叶子节点的上一个节点排序id 多个英文逗号分割',
  `clock` tinyint(4) NOT NULL DEFAULT '0' COMMENT '打卡次数',
  `is_lock` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否枷锁   0  未枷锁  1 枷锁',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 未执行，2 执行中，3 执行结束，4 执行失败，5 待执行',
  PRIMARY KEY (`id`),
  KEY `dispatch_index_parent_sort_id` (`pre_node`),
  KEY `dispatch_index_logic_id` (`logic_id`),
  KEY `diapatch_index_dispatch_node_id` (`dispatch_id`,`node_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='调度服务-调度任务具体节点信息表';


CREATE TABLE `visual_dispatch_logic` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `title` varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
  `content` text COMMENT 'shell语句',
  `remark` text COMMENT '备注',
  `business_id` text COMMENT '业务配置ID：可多个 逗号分隔',
  `server_id` int(10) NOT NULL DEFAULT '0' COMMENT '服务器管理ID',
  `create_uid` int(10) NOT NULL DEFAULT '0' COMMENT '创建用户ID',
  `update_uid` int(10) NOT NULL DEFAULT '0' COMMENT '修改用户ID',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0：正常 1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `inx_title` (`title`),
  KEY `inx_create_uid` (`create_uid`),
  KEY `inx_update_uid` (`update_uid`),
  KEY `inx_is_deleted` (`is_deleted`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='调度任务-逻辑配置表';

