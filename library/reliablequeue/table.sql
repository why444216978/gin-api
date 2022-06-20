CREATE TABLE `reliable_mq_message` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `create_user` varchar(64) NOT NULL DEFAULT '' COMMENT '创建方标识',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_user` varchar(64) NOT NULL DEFAULT '' COMMENT '更新方标识',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0-未删除，1-已删除',
  `scene` varchar(255) NOT NULL COMMENT '唯一消息scene',
  `scene_desc` varchar(255) NOT NULL COMMENT '描述信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_scene` (`scene`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='消息实体表';

CREATE TABLE `reliable_mq_message_distribute` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `create_user` varchar(64) NOT NULL DEFAULT '' COMMENT '创建方标识',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_user` varchar(64) NOT NULL DEFAULT '' COMMENT '更新方标识',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0-未删除，1-已删除',

  `message_id` bigint(20) unsigned  NOT NULL COMMENT '关联message表主键id',
  `scene` varchar(255) NOT NULL COMMENT '唯一消息scene',
  `bns` varchar(64) NOT NULL COMMENT 'bns',
  `uri` varchar(255) NOT NULL COMMENT 'uri',
  `method` varchar(32) NOT NULL COMMENT 'http method',
  PRIMARY KEY (`id`),
  KEY `idx_message_id` (`message_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='消息分发关联表';

CREATE TABLE `reliable_mq_message_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `create_user` varchar(64) NOT NULL DEFAULT '' COMMENT '创建方标识',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_user` varchar(64) NOT NULL DEFAULT '' COMMENT '更新方标识',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '版本号',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0-未删除，1-已删除',

  `message_id` bigint(20) unsigned NOT NULL COMMENT '关联message表主键id',
  `message_distribute_id` bigint(20) unsigned NOT NULL COMMENT '关联reliable_mq_message_distribute表主键id',

  `log_id` varchar(64) NOT NULL COMMENT '消息产生的log_id',
  `uuid` varchar(64) NOT NULL COMMENT '消息唯一id',
  `bns` varchar(64) NOT NULL COMMENT 'bns',
  `uri` varchar(255) NOT NULL COMMENT 'uri',
  `method` varchar(32) NOT NULL COMMENT 'http method',
  `body` text NOT NULL COMMENT '请求body',

  `delay` bigint(20) unsigned NOT NULL DEFAULT '1' COMMENT '重试间隔',
  `retry_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次重试时间',
  `next_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下次重试时间',

  `is_success` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否消费成功',
  PRIMARY KEY (`id`),
  KEY `idx_is_success` (`is_success`) USING BTREE,
  UNIQUE KEY `idx_uuid` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='业务侧可靠消息表';