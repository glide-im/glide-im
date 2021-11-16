/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80022
 Source Host           : localhost:3306
 Source Schema         : go_im

 Target Server Type    : MySQL
 Target Server Version : 80022
 File Encoding         : 65001

 Date: 16/11/2021 14:46:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_chat_message
-- ----------------------------
drop table IF EXISTS `im_chat_message`;
create TABLE `im_chat_message` (
  `m_id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID, 全局唯一自增',
  `receive_seq` bigint NOT NULL COMMENT '接收者全局 seq',
  `cli_seq` bigint NOT NULL COMMENT '发送者消息 seq',
  `from` bigint NOT NULL COMMENT '发送者 id',
  `to` bigint NOT NULL COMMENT '接收者 id',
  `type` bigint NOT NULL COMMENT '消息类型',
  `send_at` bigint NOT NULL COMMENT '发时间戳',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发送内容',
  PRIMARY KEY (`m_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_member_msg_state
-- ----------------------------
drop table IF EXISTS `im_group_member_msg_state`;
create TABLE `im_group_member_msg_state`  (
  `mb_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '群成员 ID, gid+uid 拼接',
  `g_id` bigint NOT NULL COMMENT '群 id',
  `uid` bigint NOT NULL COMMENT '成员id',
  `last_ack_m_id` bigint NOT NULL COMMENT '最后一次确认收到群消息id',
  `last_ack_seq` bigint NOT NULL COMMENT '最后一次确认收到消息的seq',
  PRIMARY KEY (`mb_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_message
-- ----------------------------
drop table IF EXISTS `im_group_message`;
create TABLE `im_group_message`  (
  `m_id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息id',
  `seq` bigint NOT NULL COMMENT '群消息seq序列号',
  `to` bigint NOT NULL COMMENT '群id',
  `from` bigint NOT NULL COMMENT '发送者id',
  `type` bigint NOT NULL COMMENT '消息类型',
  `send_at` bigint NOT NULL COMMENT '发送时间',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '消息内容',
  PRIMARY KEY (`m_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_message_state
-- ----------------------------
drop table IF EXISTS `im_group_message_state`;
create TABLE `im_group_message_state`  (
  `g_id` bigint NOT NULL AUTO_INCREMENT COMMENT '群id',
  `last_m_id` bigint NOT NULL COMMENT '最后一条群消息 id',
  `last_seq` bigint NOT NULL COMMENT '最后一条群消息 seq',
  `last_msg_at` bigint NOT NULL COMMENT '最后一条群消息 时间',
  PRIMARY KEY (`g_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_msg_seq
-- ----------------------------
drop table IF EXISTS `im_group_msg_seq`;
create TABLE `im_group_msg_seq`  (
  `g_id` bigint NOT NULL AUTO_INCREMENT COMMENT '群id',
  `seq` bigint NOT NULL COMMENT '当前 seq',
  `step` bigint NOT NULL COMMENT 'seq 步长',
  PRIMARY KEY (`g_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_offline_message
-- ----------------------------
drop table IF EXISTS `im_offline_message`;
create TABLE `im_offline_message`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '离线消息id',
  `m_id` bigint NOT NULL COMMENT '消息id',
  `uid` bigint NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
