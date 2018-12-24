/*
Navicat MySQL Data Transfer

Source Server         : aq_server
Source Server Version : 50721
Source Host           : 120.77.146.125:3306
Source Database       : aq_shop

Target Server Type    : MYSQL
Target Server Version : 50721
File Encoding         : 65001

Date: 2018-12-24 10:41:35
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for aq_ads
-- ----------------------------
DROP TABLE IF EXISTS `aq_ads`;
CREATE TABLE `aq_ads` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `ads_pos` int(255) DEFAULT NULL COMMENT '对应广告位',
  `link` varchar(255) DEFAULT NULL COMMENT '链接',
  `pic` varchar(255) DEFAULT NULL COMMENT '图片',
  `order_id` int(11) DEFAULT '999',
  `item_id` varchar(255) DEFAULT NULL COMMENT '商品id',
  `is_del` tinyint(255) DEFAULT '0' COMMENT '是否废弃',
  `width` varchar(255) DEFAULT NULL,
  `post_id` int(11) DEFAULT NULL COMMENT '对应文章',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_ads_pos
-- ----------------------------
DROP TABLE IF EXISTS `aq_ads_pos`;
CREATE TABLE `aq_ads_pos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `is_del` tinyint(255) DEFAULT '0',
  `title_pic` varchar(255) DEFAULT NULL COMMENT '标题图片',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_album
-- ----------------------------
DROP TABLE IF EXISTS `aq_album`;
CREATE TABLE `aq_album` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT '',
  `default` tinyint(4) DEFAULT '0' COMMENT '是否是默认相册',
  `cover_pic` int(255) DEFAULT '0' COMMENT '封面图片id',
  `order_id` int(11) DEFAULT '999' COMMENT '排序id',
  `is_del` tinyint(4) DEFAULT '0' COMMENT '是否废弃',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_brand
-- ----------------------------
DROP TABLE IF EXISTS `aq_brand`;
CREATE TABLE `aq_brand` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pic` varchar(255) DEFAULT NULL COMMENT '图片',
  `name` varchar(255) DEFAULT NULL COMMENT '品牌名',
  `order_id` int(11) DEFAULT NULL,
  `is_del` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_config
-- ----------------------------
DROP TABLE IF EXISTS `aq_config`;
CREATE TABLE `aq_config` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '键名',
  `value` text COMMENT '值',
  `info` text COMMENT '描述',
  `type` varchar(255) DEFAULT '' COMMENT '类型',
  `tap` varchar(255) DEFAULT '' COMMENT '标签页',
  PRIMARY KEY (`id`,`name`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_database
-- ----------------------------
DROP TABLE IF EXISTS `aq_database`;
CREATE TABLE `aq_database` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `build_time` int(11) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `user_id` int(11) DEFAULT '1' COMMENT '用户id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_export
-- ----------------------------
DROP TABLE IF EXISTS `aq_export`;
CREATE TABLE `aq_export` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '名字',
  `value` text COMMENT '值',
  `model` varchar(255) DEFAULT '' COMMENT '模块名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_export_task
-- ----------------------------
DROP TABLE IF EXISTS `aq_export_task`;
CREATE TABLE `aq_export_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL COMMENT '用户id',
  `build_time` int(11) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `userid` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=202 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_item
-- ----------------------------
DROP TABLE IF EXISTS `aq_item`;
CREATE TABLE `aq_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pics` text COMMENT '介绍图片列表',
  `name` varchar(255) DEFAULT NULL,
  `price` float DEFAULT NULL COMMENT '价格',
  `store_num` int(11) DEFAULT NULL COMMENT '库存',
  `sell_num` int(11) DEFAULT '0' COMMENT '销量',
  `is_onsale` tinyint(4) DEFAULT '0' COMMENT '是否上架 1：上架 0 下架',
  `order_id` int(11) DEFAULT '999' COMMENT '排序',
  `item_type` int(11) DEFAULT NULL COMMENT '商品类别',
  `brand` int(255) DEFAULT NULL,
  `spec` text COMMENT '规格',
  `desc` text COMMENT '详情',
  `tag` text COMMENT '标签',
  `code` varchar(255) DEFAULT NULL COMMENT '商品编码',
  `icon` varchar(255) DEFAULT '' COMMENT '图标',
  `basenum` int(11) DEFAULT '3' COMMENT '基数(必须是这个数量倍数)',
  `group_price` text COMMENT '组价格',
  `is_sync_shipnum` tinyint(4) DEFAULT '1' COMMENT '导出erp时是否同时物流信息',
  `no_service` tinyint(4) DEFAULT '0' COMMENT '是否有特殊服务',
  `item_unit` varchar(255) DEFAULT NULL COMMENT '单位',
  `item_shelf_life` varchar(255) DEFAULT NULL COMMENT '保质期',
  `idnum_need` tinyint(4) DEFAULT '0' COMMENT '是否需要身份证',
  `supply_source` tinyint(4) DEFAULT '1' COMMENT '发货方式',
  `min_num` int(11) DEFAULT '1' COMMENT '最小数量',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=205 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_item_type
-- ----------------------------
DROP TABLE IF EXISTS `aq_item_type`;
CREATE TABLE `aq_item_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL COMMENT '分类名称',
  `code` varchar(255) NOT NULL COMMENT '代码',
  `info` text COMMENT '备注',
  `level` int(3) DEFAULT '1' COMMENT '第几层',
  `parent_id` int(11) DEFAULT '0' COMMENT '父节点id',
  `order_id` int(3) DEFAULT '100' COMMENT '排序id',
  `pic` varchar(255) DEFAULT '',
  `is_del` tinyint(4) DEFAULT '0',
  `intro_text` text CHARACTER SET utf8mb4 COMMENT '文本',
  PRIMARY KEY (`id`,`code`),
  UNIQUE KEY `id` (`id`),
  KEY `parent_id` (`parent_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_log
-- ----------------------------
DROP TABLE IF EXISTS `aq_log`;
CREATE TABLE `aq_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) DEFAULT NULL,
  `time` int(11) DEFAULT NULL,
  `info` text CHARACTER SET utf8mb4 COMMENT '详情',
  `controller` varchar(255) DEFAULT NULL,
  `method` varchar(255) DEFAULT NULL,
  `link` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=52280 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_logistics
-- ----------------------------
DROP TABLE IF EXISTS `aq_logistics`;
CREATE TABLE `aq_logistics` (
  `id` varchar(255) NOT NULL COMMENT '物流单号',
  `internal_ship_company_code` varchar(255) DEFAULT '' COMMENT '国内物流公司编码',
  `internal_ship_num` varchar(255) DEFAULT '' COMMENT '国内物流号',
  `build_time` int(11) DEFAULT '0' COMMENT '创建时间',
  `logistics_task` text,
  `logistics_task_starttime` int(11) DEFAULT '0',
  `state` tinyint(4) DEFAULT '0' COMMENT '物流状态',
  `idnum` varchar(255) DEFAULT '',
  `client_name` varchar(255) DEFAULT '' COMMENT '收件人姓名',
  `idnumpic1` varchar(255) DEFAULT '',
  `idnumpic2` varchar(255) DEFAULT '',
  `sync_erp_flag` tinyint(4) DEFAULT '0' COMMENT '同步到erp的状态',
  `client_phone` varchar(255) DEFAULT '' COMMENT '收件人电话',
  `client_address` text COMMENT '收件人地址',
  `is_del` tinyint(4) DEFAULT '0' COMMENT '是否废弃',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`) USING HASH,
  KEY `phone` (`client_phone`),
  KEY `name` (`client_name`),
  KEY `namephone` (`client_name`,`client_phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_logistics_task
-- ----------------------------
DROP TABLE IF EXISTS `aq_logistics_task`;
CREATE TABLE `aq_logistics_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tasklist` text CHARACTER SET utf8,
  `name` varchar(255) CHARACTER SET utf8 DEFAULT NULL COMMENT '物流名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for aq_module
-- ----------------------------
DROP TABLE IF EXISTS `aq_module`;
CREATE TABLE `aq_module` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT '',
  `controller` varchar(255) DEFAULT '',
  `method` varchar(255) DEFAULT '',
  `posid` int(11) DEFAULT '0',
  `need_auth` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `name` (`name`,`method`)
) ENGINE=InnoDB AUTO_INCREMENT=121 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_notice
-- ----------------------------
DROP TABLE IF EXISTS `aq_notice`;
CREATE TABLE `aq_notice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` text,
  `order_id` int(11) DEFAULT '999',
  `build_time` int(11) DEFAULT '0',
  `content` text CHARACTER SET utf8mb4,
  `is_del` int(255) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_order
-- ----------------------------
DROP TABLE IF EXISTS `aq_order`;
CREATE TABLE `aq_order` (
  `id` varchar(255) NOT NULL COMMENT '订单号',
  `order_time` int(11) DEFAULT NULL COMMENT '下单时间',
  `pay_time` int(11) DEFAULT '0' COMMENT '支付时间',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态',
  `client_address` text,
  `client_name` varchar(255) DEFAULT '' COMMENT '收货姓名',
  `user_id` int(11) DEFAULT NULL COMMENT '会员id',
  `idnum` varchar(255) DEFAULT '' COMMENT '收货人身份证号',
  `client_phone` varchar(255) DEFAULT '' COMMENT '收货人电话',
  `close_info` varchar(255) DEFAULT NULL COMMENT '卖家关闭订单注释',
  `close_type` tinyint(255) DEFAULT '0' COMMENT '订单关闭类型',
  `shipment_num` text COMMENT '物流单号(数组)',
  `client_info` varchar(255) DEFAULT NULL COMMENT '买家备注',
  `sell_info` varchar(255) DEFAULT NULL COMMENT '卖家备注',
  `close_time` int(11) DEFAULT '0' COMMENT '订单结束时间',
  `client_provice` varchar(255) DEFAULT NULL,
  `client_city` varchar(255) DEFAULT NULL,
  `client_area` varchar(255) DEFAULT NULL,
  `idnumpic1` varchar(255) DEFAULT '' COMMENT '身份证图片',
  `idnumpic2` varchar(255) DEFAULT '',
  `refund_info` varchar(255) DEFAULT NULL COMMENT '退款原因',
  `order_vip_type` tinyint(255) DEFAULT '0' COMMENT '特殊要求',
  `pay_id` varchar(255) NOT NULL COMMENT '支付id',
  `total_price` varchar(255) DEFAULT NULL COMMENT '订单总价',
  `send_user_name` varchar(255) DEFAULT NULL,
  `send_user_phone` varchar(255) DEFAULT NULL,
  `itemid` varchar(255) DEFAULT NULL COMMENT '商品id',
  `num` int(11) DEFAULT NULL COMMENT '商品数量',
  `specname` varchar(255) DEFAULT NULL COMMENT '商品规格',
  `itemcode` varchar(255) DEFAULT NULL COMMENT '商品编码',
  `itempic` varchar(255) DEFAULT NULL COMMENT '商品图片',
  `unitprice` varchar(255) DEFAULT NULL COMMENT '商品单价',
  `freight_price` int(11) DEFAULT '0' COMMENT '运费',
  `service_price` int(11) DEFAULT '0' COMMENT '服务费',
  `pay_type` tinyint(4) DEFAULT '0' COMMENT '支付方式',
  `pay_check_info` varchar(255) DEFAULT NULL COMMENT '支付审核描述',
  `supply_source` tinyint(4) DEFAULT '1' COMMENT '发货方式',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `payid` (`pay_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_paycode
-- ----------------------------
DROP TABLE IF EXISTS `aq_paycode`;
CREATE TABLE `aq_paycode` (
  `id` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '支付码',
  `order_list` text CHARACTER SET utf8 COMMENT '对应的订单列表',
  `money` float DEFAULT '0' COMMENT '对应金额',
  `user_id` int(11) DEFAULT NULL COMMENT '用户id',
  `build_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for aq_photo
-- ----------------------------
DROP TABLE IF EXISTS `aq_photo`;
CREATE TABLE `aq_photo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `path` varchar(255) DEFAULT '',
  `name` varchar(255) DEFAULT '',
  `upload_time` int(11) DEFAULT NULL,
  `upload_user` int(255) DEFAULT NULL,
  `width` int(11) DEFAULT '0',
  `height` int(11) DEFAULT '0',
  `album` int(11) DEFAULT '0' COMMENT '相册',
  `order_id` int(11) DEFAULT '999' COMMENT '排序id',
  `key` varchar(255) DEFAULT NULL COMMENT '七牛的key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=536 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_post
-- ----------------------------
DROP TABLE IF EXISTS `aq_post`;
CREATE TABLE `aq_post` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `build_user` int(255) DEFAULT NULL,
  `build_time` int(11) DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4,
  `type` int(255) DEFAULT NULL,
  `is_del` tinyint(255) DEFAULT '0' COMMENT '是否废弃',
  `summary` text COMMENT '摘要',
  `order_id` int(11) DEFAULT '0' COMMENT '排序',
  `pic` varchar(255) DEFAULT NULL COMMENT '对应图片',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_post_type
-- ----------------------------
DROP TABLE IF EXISTS `aq_post_type`;
CREATE TABLE `aq_post_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `is_del` tinyint(255) DEFAULT '0' COMMENT '是否废弃',
  `order_id` int(11) DEFAULT '999' COMMENT '排序id',
  `parent_id` int(11) DEFAULT '0' COMMENT '父类别',
  `level` int(11) DEFAULT '1' COMMENT '层级',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_tag
-- ----------------------------
DROP TABLE IF EXISTS `aq_tag`;
CREATE TABLE `aq_tag` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pic` varchar(255) DEFAULT '',
  `name` varchar(255) DEFAULT '',
  `is_del` tinyint(255) DEFAULT '0',
  `order_id` int(11) DEFAULT '999',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for aq_user
-- ----------------------------
DROP TABLE IF EXISTS `aq_user`;
CREATE TABLE `aq_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(190) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '账号名',
  `name` varchar(190) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '昵称',
  `mail` varchar(255) CHARACTER SET utf8 DEFAULT '' COMMENT '邮箱',
  `reg_time` int(11) DEFAULT '0' COMMENT '注册时间',
  `phone` varchar(255) CHARACTER SET utf8 DEFAULT '' COMMENT '手机号',
  `password` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '密码',
  `is_del` int(1) NOT NULL DEFAULT '0' COMMENT '是否在有效',
  `user_group` int(11) NOT NULL DEFAULT '0' COMMENT '用户所属用户组',
  `user_token` text CHARACTER SET utf8,
  `token_expire` int(11) DEFAULT '0' COMMENT 'token过期时间',
  `token_get_time` int(11) DEFAULT NULL COMMENT 'token获取时间',
  `shop_cart` text CHARACTER SET utf8 COMMENT '购物车',
  `wchat_openid` varchar(255) CHARACTER SET utf8 DEFAULT NULL COMMENT '微信号',
  `last_login_time` int(11) DEFAULT NULL COMMENT '上次登录时间',
  `head` varchar(255) CHARACTER SET utf8 DEFAULT NULL COMMENT '头像',
  `address` text CHARACTER SET utf8 COMMENT '收货地址',
  `country` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `province` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `city` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `wchat_unionid` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `sex` tinyint(255) DEFAULT NULL,
  `track_admin` int(255) DEFAULT '0' COMMENT '跟单员',
  PRIMARY KEY (`id`,`account`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=395 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for aq_usergroup
-- ----------------------------
DROP TABLE IF EXISTS `aq_usergroup`;
CREATE TABLE `aq_usergroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户组名称',
  `module_ids` text COMMENT '权限组',
  `group_type` int(1) NOT NULL DEFAULT '0' COMMENT '1000以上： 系统管理员    100以上：管理员  100以下会员',
  `expire_time` int(11) DEFAULT '86400' COMMENT '过期时间（0表示永不过期）',
  `limit_show_order` tinyint(4) DEFAULT '0' COMMENT '是否限制显示订单数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
