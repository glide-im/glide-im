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

 Date: 12/01/2022 16:22:21
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_chat_message
-- ----------------------------
DROP TABLE IF EXISTS `im_chat_message`;
CREATE TABLE `im_chat_message`  (
  `m_id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cli_seq` bigint NULL DEFAULT NULL,
  `from` bigint NULL DEFAULT NULL,
  `to` bigint NULL DEFAULT NULL,
  `type` int NULL DEFAULT NULL,
  `send_at` bigint NULL DEFAULT NULL,
  `create_at` bigint NULL DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`m_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 123432 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_contacts
-- ----------------------------
DROP TABLE IF EXISTS `im_contacts`;
CREATE TABLE `im_contacts`  (
  `fid` char(254) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `uid` bigint NOT NULL,
  `id` bigint NOT NULL,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` int NOT NULL,
  PRIMARY KEY (`fid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_member_model
-- ----------------------------
DROP TABLE IF EXISTS `im_group_member_model`;
CREATE TABLE `im_group_member_model`  (
  `mb_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `gid` bigint NULL DEFAULT NULL,
  `uid` bigint NULL DEFAULT NULL,
  `flag` bigint NULL DEFAULT NULL,
  `type` bigint NULL DEFAULT NULL,
  `remark` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`mb_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_member_msg_state
-- ----------------------------
DROP TABLE IF EXISTS `im_group_member_msg_state`;
CREATE TABLE `im_group_member_msg_state`  (
  `mb_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `g_id` bigint NULL DEFAULT NULL,
  `uid` bigint NULL DEFAULT NULL,
  `last_ack_m_id` bigint NULL DEFAULT NULL,
  `last_ack_seq` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`mb_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_message
-- ----------------------------
DROP TABLE IF EXISTS `im_group_message`;
CREATE TABLE `im_group_message`  (
  `m_id` bigint NOT NULL AUTO_INCREMENT,
  `seq` bigint NULL DEFAULT NULL,
  `to` bigint NULL DEFAULT NULL,
  `from` bigint NULL DEFAULT NULL,
  `type` bigint NULL DEFAULT NULL,
  `send_at` bigint NULL DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`m_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1231241239 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_message_state
-- ----------------------------
DROP TABLE IF EXISTS `im_group_message_state`;
CREATE TABLE `im_group_message_state`  (
  `gid` bigint NOT NULL AUTO_INCREMENT,
  `last_m_id` bigint NULL DEFAULT NULL,
  `last_seq` bigint NULL DEFAULT NULL,
  `last_msg_at` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`gid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_model
-- ----------------------------
DROP TABLE IF EXISTS `im_group_model`;
CREATE TABLE `im_group_model`  (
  `gid` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `mute` tinyint(1) NULL DEFAULT NULL,
  `flag` int NULL DEFAULT NULL,
  `create_at` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`gid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_msg_seq
-- ----------------------------
DROP TABLE IF EXISTS `im_group_msg_seq`;
CREATE TABLE `im_group_msg_seq`  (
  `gid` bigint NOT NULL AUTO_INCREMENT,
  `seq` bigint NULL DEFAULT NULL,
  `step` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`gid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_offline_message
-- ----------------------------
DROP TABLE IF EXISTS `im_offline_message`;
CREATE TABLE `im_offline_message`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `m_id` bigint NULL DEFAULT NULL,
  `uid` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 82 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_session
-- ----------------------------
DROP TABLE IF EXISTS `im_session`;
CREATE TABLE `im_session`  (
  `session_id` char(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `uid` bigint NOT NULL,
  `update_at` bigint NOT NULL,
  `uid2` bigint NOT NULL,
  `create_at` bigint NOT NULL,
  `last_m_id` bigint NOT NULL,
  PRIMARY KEY (`session_id`) USING BTREE
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
  `create_at` bigint NULL DEFAULT NULL,
  `update_at` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`uid`) USING BTREE,
  UNIQUE INDEX `account`(`account`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 543623 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
