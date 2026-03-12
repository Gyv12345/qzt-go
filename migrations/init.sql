-- QZT-GO 数据库初始化脚本
-- 创建时间: 2026-03-11
-- 说明: QZT CRM 系统基础数据表结构

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(30) NOT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(加密)',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `department_id` varchar(30) DEFAULT NULL COMMENT '部门ID',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` varchar(30) NOT NULL COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ----------------------------
-- Table structure for departments
-- ----------------------------
DROP TABLE IF EXISTS `departments`;
CREATE TABLE `departments` (
  `id` varchar(30) NOT NULL COMMENT '部门ID',
  `name` varchar(50) NOT NULL COMMENT '部门名称',
  `code` varchar(50) DEFAULT NULL COMMENT '部门编码',
  `parent_id` varchar(30) DEFAULT NULL COMMENT '父部门ID',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ----------------------------
-- Table structure for menus
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
  `id` varchar(30) NOT NULL COMMENT '菜单ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `code` varchar(50) DEFAULT NULL COMMENT '菜单编码',
  `parent_id` varchar(30) DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(255) DEFAULT NULL COMMENT '路径',
  `icon` varchar(50) DEFAULT NULL COMMENT '图标',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `type` varchar(20) DEFAULT NULL COMMENT '类型 menu/button',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `user_id` varchar(30) NOT NULL COMMENT '用户ID',
  `role_id` varchar(30) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ----------------------------
-- Table structure for role_menus
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus` (
  `role_id` varchar(30) NOT NULL COMMENT '角色ID',
  `menu_id` varchar(30) NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ----------------------------
-- Table structure for customers
-- ----------------------------
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` varchar(30) NOT NULL COMMENT '客户ID',
  `name` varchar(100) NOT NULL COMMENT '公司名称',
  `short_name` varchar(50) DEFAULT NULL COMMENT '简称',
  `code` varchar(50) DEFAULT NULL COMMENT '客户编码',
  `industry` varchar(50) DEFAULT NULL COMMENT '行业',
  `scale` varchar(50) DEFAULT NULL COMMENT '规模',
  `address` varchar(255) DEFAULT NULL COMMENT '地址',
  `website` varchar(255) DEFAULT NULL COMMENT '网站',
  `customer_level` varchar(20) DEFAULT 'LEAD' COMMENT '客户等级 LEAD/PROSPECT/CUSTOMER/VIP',
  `source_channel` varchar(50) DEFAULT NULL COMMENT '来源渠道',
  `follow_user_id` varchar(30) DEFAULT NULL COMMENT '跟进人ID',
  `tags` text COMMENT '标签(JSON)',
  `remark` text COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `first_contact_date` datetime DEFAULT NULL COMMENT '首次联系日期',
  `contract_date` datetime DEFAULT NULL COMMENT '签约日期',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_follow_user` (`follow_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户表';

-- ----------------------------
-- Table structure for contacts
-- ----------------------------
DROP TABLE IF EXISTS `contacts`;
CREATE TABLE `contacts` (
  `id` varchar(30) NOT NULL COMMENT '联系人ID',
  `customer_id` varchar(30) NOT NULL COMMENT '客户ID',
  `name` varchar(50) NOT NULL COMMENT '姓名',
  `position` varchar(50) DEFAULT NULL COMMENT '职位',
  `department` varchar(50) DEFAULT NULL COMMENT '部门',
  `phone` varchar(20) DEFAULT NULL COMMENT '固定电话',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `wechat` varchar(50) DEFAULT NULL COMMENT '微信号',
  `qq` varchar(20) DEFAULT NULL COMMENT 'QQ号',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要联系人',
  `is_decision` tinyint(1) DEFAULT '0' COMMENT '是否决策人',
  `gender` varchar(10) DEFAULT NULL COMMENT '性别',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `address` varchar(255) DEFAULT NULL COMMENT '地址',
  `remark` text COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='联系人表';

-- ----------------------------
-- Table structure for contracts
-- ----------------------------
DROP TABLE IF EXISTS `contracts`;
CREATE TABLE `contracts` (
  `id` varchar(30) NOT NULL COMMENT '合同ID',
  `customer_id` varchar(30) NOT NULL COMMENT '客户ID',
  `contract_no` varchar(50) DEFAULT NULL COMMENT '合同编号',
  `title` varchar(200) NOT NULL COMMENT '合同标题',
  `type` varchar(50) DEFAULT NULL COMMENT '合同类型',
  `amount` decimal(12,2) DEFAULT NULL COMMENT '合同金额',
  `paid_amount` decimal(12,2) DEFAULT '0.00' COMMENT '已付金额',
  `status` varchar(20) DEFAULT 'DRAFT' COMMENT '状态 DRAFT/PENDING/ACTIVE/COMPLETED/CANCELLED',
  `start_date` date DEFAULT NULL COMMENT '开始日期',
  `end_date` date DEFAULT NULL COMMENT '结束日期',
  `sign_date` date DEFAULT NULL COMMENT '签约日期',
  `owner_user_id` varchar(30) DEFAULT NULL COMMENT '负责人ID',
  `tags` text COMMENT '标签(JSON)',
  `remark` text COMMENT '备注',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合同表';

-- ----------------------------
-- Table structure for api_keys
-- ----------------------------
DROP TABLE IF EXISTS `api_keys`;
CREATE TABLE `api_keys` (
  `id` varchar(30) NOT NULL COMMENT 'API Key ID',
  `name` varchar(100) NOT NULL COMMENT '名称',
  `key_hash` varchar(64) NOT NULL COMMENT 'Key Hash',
  `prefix` varchar(10) NOT NULL COMMENT 'Key 前缀',
  `user_id` varchar(30) NOT NULL COMMENT '所属用户ID',
  `permissions` text COMMENT '权限(JSON)',
  `ip_whitelist` text COMMENT 'IP白名单(JSON)',
  `rate_limit` int(11) DEFAULT '100' COMMENT '速率限制',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  `last_used_at` datetime DEFAULT NULL COMMENT '最后使用时间',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key_hash` (`key_hash`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API Key表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 插入管理员用户 (密码: admin123)
-- 密码经过 bcrypt 加密: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
INSERT INTO `users` (`id`, `username`, `password`, `email`, `real_name`, `status`, `created_at`, `updated_at`)
VALUES ('1', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@qzt.com', '系统管理员', 1, NOW(), NOW());

-- 插入默认角色
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`)
VALUES 
('1', '超级管理员', 'ADMIN', '拥有所有权限', 1, NOW(), NOW()),
('2', '销售经理', 'SALES_MANAGER', '销售管理权限', 1, NOW(), NOW()),
('3', '销售人员', 'SALES', '销售人员权限', 1, NOW(), NOW());

-- 插入管理员角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES ('1', '1');

SET FOREIGN_KEY_CHECKS = 1;
