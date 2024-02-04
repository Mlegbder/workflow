/*
 Navicat Premium Data Transfer

 Source Server         : lianfei-dev
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : 118.190.154.32:3306
 Source Schema         : womata_sys

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 04/02/2024 17:33:35
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for proc_def
-- ----------------------------
DROP TABLE IF EXISTS `proc_def`;
CREATE TABLE `proc_def` (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                            `proc_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '流程名称',
                            `version` int NOT NULL DEFAULT '1' COMMENT '版本 默认1',
                            `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
                            `resource` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '流程JSON',
                            `created_at` datetime NOT NULL COMMENT '创建时间',
                            `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
                            `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程定义表';

SET FOREIGN_KEY_CHECKS = 1;



/*
 Navicat Premium Data Transfer

 Source Server         : lianfei-dev
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : 118.190.154.32:3306
 Source Schema         : womata_sys

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 04/02/2024 17:33:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for proc_inst
-- ----------------------------
DROP TABLE IF EXISTS `proc_inst`;
CREATE TABLE `proc_inst` (
                             `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `proc_def_id` bigint NOT NULL COMMENT '流程定义ID',
                             `proc_def_version` int NOT NULL COMMENT '流程版本号',
                             `node_info` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '流程节点',
                             `node_id` bigint NOT NULL COMMENT '当前节点',
                             `is_complete` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否完成 1:进行中 2已完成',
                             `assignee` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '发起人',
                             `created_at` datetime NOT NULL COMMENT '创建时间',
                             `updated_at` datetime DEFAULT NULL COMMENT '更新时间/最后审批时间',
                             PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2234 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程实例表';

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : lianfei-dev
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : 118.190.154.32:3306
 Source Schema         : womata_sys

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 04/02/2024 17:33:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for proc_task
-- ----------------------------
DROP TABLE IF EXISTS `proc_task`;
CREATE TABLE `proc_task` (
                             `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `proc_inst_id` bigint NOT NULL COMMENT '流程实例ID',
                             `node_id` bigint NOT NULL COMMENT '节点ID',
                             `assignee` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务处理人',
                             `member_count` int NOT NULL DEFAULT '1' COMMENT '表示当前任务需要多少人审批之后才能结束，默认是 1',
                             `un_complete_num` int NOT NULL DEFAULT '1' COMMENT '表示还有多少人没有审批，默认是1',
                             `agree_num` int NOT NULL DEFAULT '0' COMMENT '表示通过的人数',
                             `act_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '表示任务类型 "or"表示或签，即一个人通过或者驳回就结束，"and"表示会签，要所有人通过就流\n转到下一步，如果有一个人驳回那么就跳转到上一步',
                             `is_complete` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否完成 1:进行中 2已完成',
                             `created_at` datetime NOT NULL COMMENT '创建时间',
                             `updated_at` datetime DEFAULT NULL COMMENT '更新时间/最后审批时间',
                             PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4179 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程任务表';

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : lianfei-dev
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : 118.190.154.32:3306
 Source Schema         : womata_sys

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 04/02/2024 17:33:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for proc_his
-- ----------------------------
DROP TABLE IF EXISTS `proc_his`;
CREATE TABLE `proc_his` (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                            `assignee` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务处理人',
                            `proc_inst_id` bigint NOT NULL COMMENT '流程实例ID',
                            `node_id` bigint NOT NULL COMMENT '节点id',
                            `approval_status` tinyint(1) NOT NULL COMMENT '审批状态:  1.审批通过 2.审批驳回  3.发起流程 4.中止流程 5.结束流程',
                            `created_at` datetime NOT NULL COMMENT '创建时间',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=247473 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程历史记录表';

SET FOREIGN_KEY_CHECKS = 1;
