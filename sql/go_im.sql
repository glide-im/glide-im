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

 Date: 16/11/2021 16:01:51
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_chat
-- ----------------------------
DROP TABLE IF EXISTS `im_chat`;
CREATE TABLE `im_chat`  (
  `cid` bigint NOT NULL AUTO_INCREMENT,
  `chat_type` tinyint NULL DEFAULT NULL,
  `target_id` bigint NULL DEFAULT NULL,
  `current_mid` bigint NULL DEFAULT NULL,
  `new_message_at` datetime NULL DEFAULT NULL,
  `create_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`cid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_chat_message
-- ----------------------------
DROP TABLE IF EXISTS `im_chat_message`;
CREATE TABLE `im_chat_message`  (
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
-- Table structure for im_chat_message_id
-- ----------------------------
DROP TABLE IF EXISTS `im_chat_message_id`;
CREATE TABLE `im_chat_message_id`  (
  `cid` bigint NOT NULL AUTO_INCREMENT,
  `current_mid` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`cid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_contacts
-- ----------------------------
DROP TABLE IF EXISTS `im_contacts`;
CREATE TABLE `im_contacts`  (
  `fid` bigint NOT NULL AUTO_INCREMENT,
  `owner` bigint NULL DEFAULT NULL,
  `target_id` bigint NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `type` tinyint NULL DEFAULT NULL,
  `add_time` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`fid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group
-- ----------------------------
DROP TABLE IF EXISTS `im_group`;
CREATE TABLE `im_group`  (
  `gid` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `owner` bigint NULL DEFAULT NULL,
  `mute` tinyint(1) NULL DEFAULT NULL,
  `notice` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `chat_id` bigint NULL DEFAULT NULL,
  `create_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`gid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_member
-- ----------------------------
DROP TABLE IF EXISTS `im_group_member`;
CREATE TABLE `im_group_member`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `gid` bigint NULL DEFAULT NULL,
  `uid` bigint NULL DEFAULT NULL,
  `mute` bigint NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `flag` int NULL DEFAULT NULL,
  `join_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_member_msg_state
-- ----------------------------
DROP TABLE IF EXISTS `im_group_member_msg_state`;
CREATE TABLE `im_group_member_msg_state`  (
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
DROP TABLE IF EXISTS `im_group_message`;
CREATE TABLE `im_group_message`  (
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
DROP TABLE IF EXISTS `im_group_message_state`;
CREATE TABLE `im_group_message_state`  (
  `g_id` bigint NOT NULL AUTO_INCREMENT COMMENT '群id',
  `last_m_id` bigint NOT NULL COMMENT '最后一条群消息 id',
  `last_seq` bigint NOT NULL COMMENT '最后一条群消息 seq',
  `last_msg_at` bigint NOT NULL COMMENT '最后一条群消息 时间',
  PRIMARY KEY (`g_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_msg_seq
-- ----------------------------
DROP TABLE IF EXISTS `im_group_msg_seq`;
CREATE TABLE `im_group_msg_seq`  (
  `g_id` bigint NOT NULL AUTO_INCREMENT COMMENT '群id',
  `seq` bigint NOT NULL COMMENT '当前 seq\r\n',
  `step` bigint NOT NULL COMMENT 'seq 步长',
  PRIMARY KEY (`g_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_offline_message
-- ----------------------------
DROP TABLE IF EXISTS `im_offline_message`;
CREATE TABLE `im_offline_message`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '离线消息id',
  `m_id` bigint NOT NULL COMMENT '消息id',
  `uid` bigint NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_user
-- ----------------------------
DROP TABLE IF EXISTS `im_user`;
CREATE TABLE `im_user`  (
  `uid` bigint NOT NULL AUTO_INCREMENT,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_at` datetime NULL DEFAULT NULL,
  `update_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`uid`) USING BTREE,
  UNIQUE INDEX `account`(`account`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_user_chat
-- ----------------------------
DROP TABLE IF EXISTS `im_user_chat`;
CREATE TABLE `im_user_chat`  (
  `uc_id` bigint NOT NULL AUTO_INCREMENT,
  `ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cid` bigint NULL DEFAULT NULL,
  `owner` bigint NULL DEFAULT NULL,
  `target` bigint NULL DEFAULT NULL,
  `chat_type` tinyint NULL DEFAULT NULL,
  `unread` int NULL DEFAULT NULL,
  `new_message_at` datetime NULL DEFAULT NULL,
  `read_at` datetime NULL DEFAULT NULL,
  `create_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`uc_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
