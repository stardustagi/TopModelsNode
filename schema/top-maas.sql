/*
 Navicat Premium Data Transfer

 Source Server         : mysql-dev
 Source Server Type    : MySQL
 Source Server Version : 100612
 Source Host           : 127.0.0.1:3306
 Source Schema         : top-maas

 Target Server Type    : MySQL
 Target Server Version : 100612
 File Encoding         : 65001

 Date: 26/11/2025 14:14:06
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for company
-- ----------------------------
DROP TABLE IF EXISTS `company`;
CREATE TABLE `company`  (
  `id` bigint(12) NOT NULL,
  `company_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '公司名',
  `created_date` bigint(12) NULL DEFAULT NULL COMMENT '创建日期',
  `license_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '营业执照信息',
  `is_realname_authentication` tinyint(1) NULL DEFAULT NULL COMMENT '是否实名认证',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '公司地址',
  `phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '公司电话',
  `mail` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '公司邮件',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for invitation_code
-- ----------------------------
DROP TABLE IF EXISTS `invitation_code`;
CREATE TABLE `invitation_code`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `code` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `admin_id` bigint(12) NULL DEFAULT NULL,
  `created_at` bigint(12) NULL DEFAULT NULL,
  `is_use` tinyint(1) NULL DEFAULT NULL,
  `use_time` bigint(12) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for llm_usage_report
-- ----------------------------
DROP TABLE IF EXISTS `llm_usage_report`;
CREATE TABLE `llm_usage_report`  (
  `id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `node_id` bigint(12) NULL DEFAULT NULL,
  `model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `model_id` bigint(12) NULL DEFAULT NULL,
  `actual_model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `provider` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `actual_provider` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `caller` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `caller_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `client_version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `token_usage_input_tokens` int(11) NULL DEFAULT NULL,
  `token_usage_output_tokens` int(11) NULL DEFAULT NULL,
  `token_usage_cached_tokens` int(11) NULL DEFAULT NULL,
  `token_usage_reasoning_tokens` int(11) NULL DEFAULT NULL,
  `token_usage_tokens_per_sec` int(11) NULL DEFAULT NULL,
  `token_usage_latency` double NULL DEFAULT NULL,
  `agent_version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `stream` tinyint(1) NOT NULL DEFAULT 0,
  `updated_at` bigint(12) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for models_info
-- ----------------------------
DROP TABLE IF EXISTS `models_info`;
CREATE TABLE `models_info`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `model_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '模型ID',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '模型名',
  `api_version` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `deploy_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `input_price` int(10) NULL DEFAULT NULL,
  `output_price` int(10) NULL DEFAULT NULL,
  `cache_price` int(10) NULL DEFAULT NULL,
  `status` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '模型状态',
  `last_update` bigint(20) NULL DEFAULT NULL COMMENT '最后更新时间',
  `is_private` tinyint(1) NULL DEFAULT NULL COMMENT '是否私有化',
  `owner_id` bigint(12) NULL DEFAULT NULL COMMENT '为0则为平台模型',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '模型地址',
  `api_styles` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT 'Api风格',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 97 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for models_provider
-- ----------------------------
DROP TABLE IF EXISTS `models_provider`;
CREATE TABLE `models_provider`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `owner_id` bigint(12) NULL DEFAULT NULL COMMENT '模型供应商的nodeUserId',
  `provider_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `endpoint` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `api_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `model_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `input_price` int(10) NULL DEFAULT NULL,
  `output_price` int(10) NULL DEFAULT NULL,
  `cache_price` int(10) NULL DEFAULT NULL,
  `api_keys` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT 'apikeys列表',
  `deleted` bigint(12) NULL DEFAULT NULL,
  `last_update` bigint(12) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 76 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for node_models_info_maps
-- ----------------------------
DROP TABLE IF EXISTS `node_models_info_maps`;
CREATE TABLE `node_models_info_maps`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `node_id` bigint(12) NULL DEFAULT NULL,
  `node_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `model_id` bigint(12) NULL DEFAULT NULL,
  `model_provider_id` bigint(12) NULL DEFAULT NULL,
  `created_at` bigint(12) NULL DEFAULT NULL,
  `updated_at` bigint(12) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for node_users
-- ----------------------------
DROP TABLE IF EXISTS `node_users`;
CREATE TABLE `node_users`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '用户邮件',
  `salt` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '用户密码',
  `created_at` bigint(12) NULL DEFAULT NULL COMMENT '创建日期',
  `deleted` bigint(12) NULL DEFAULT NULL COMMENT '删除日期',
  `last_update` bigint(12) NULL DEFAULT NULL COMMENT '最后登录日期',
  `is_active` tinyint(1) NULL DEFAULT NULL COMMENT '是否激活',
  `company_id` bigint(12) NULL DEFAULT NULL COMMENT '公司ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for nodes
-- ----------------------------
DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `name` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `owner_id` bigint(12) NULL DEFAULT NULL COMMENT '节点所有者ID对应nodeUserID',
  `created_at` bigint(12) NULL DEFAULT NULL,
  `lastupdate_at` bigint(12) NULL DEFAULT NULL,
  `domain` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `access_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT 'ak',
  `security_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT 'sk',
  `company_id` bigint(12) NULL DEFAULT NULL COMMENT '企业ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for system_config
-- ----------------------------
DROP TABLE IF EXISTS `system_config`;
CREATE TABLE `system_config`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '配置key',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `created_at` bigint(12) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_agent_tokens
-- ----------------------------
DROP TABLE IF EXISTS `user_agent_tokens`;
CREATE TABLE `user_agent_tokens`  (
  `id` bigint(12) NOT NULL,
  `user_id` bigint(12) NULL DEFAULT NULL COMMENT '用户ID',
  `user_agent_tokens` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '用户端TokenList',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_api_keys
-- ----------------------------
DROP TABLE IF EXISTS `user_api_keys`;
CREATE TABLE `user_api_keys`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(12) NULL DEFAULT NULL COMMENT '用户ID',
  `last_update` bigint(12) NULL DEFAULT NULL COMMENT '最后更新时间',
  `api_keys` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT 'api表',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_consume_record
-- ----------------------------
DROP TABLE IF EXISTS `user_consume_record`;
CREATE TABLE `user_consume_record`  (
  `id` bigint(18) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(18) NULL DEFAULT NULL,
  `total_consumed` bigint(18) NULL DEFAULT NULL,
  `caller` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '调用方',
  `model` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '模型',
  `model_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '模型ID',
  `input_tokens` bigint(18) NULL DEFAULT NULL COMMENT '输入token数',
  `output_tokens` bigint(18) NULL DEFAULT NULL COMMENT '输出token数',
  `cache_tokens` bigint(18) NULL DEFAULT NULL COMMENT '缓存token数',
  `input_price` int(10) NULL DEFAULT NULL COMMENT '输入token价格',
  `output_price` int(10) NULL DEFAULT NULL COMMENT '输出token价格',
  `cache_price` int(10) NULL DEFAULT NULL COMMENT '缓存token价格\r\n',
  `created_at` bigint(18) NULL DEFAULT NULL COMMENT '创建时间',
  `update_at` bigint(18) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_models_infos
-- ----------------------------
DROP TABLE IF EXISTS `user_models_infos`;
CREATE TABLE `user_models_infos`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `node_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '节点ID',
  `model_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL COMMENT '模型IDS',
  `user_id` bigint(128) NULL DEFAULT NULL COMMENT '用户ID',
  `last_update` bigint(12) NULL DEFAULT NULL COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_pay_log
-- ----------------------------
DROP TABLE IF EXISTS `user_pay_log`;
CREATE TABLE `user_pay_log`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NULL DEFAULT NULL,
  `pay_amount` bigint(18) NULL DEFAULT NULL COMMENT '充值金额，单位分',
  `pay_time` bigint(18) NULL DEFAULT NULL COMMENT '充值时间',
  `pay_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '充值渠道',
  `pay_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '支付原因',
  `admin_user_id` bigint(20) NULL DEFAULT NULL COMMENT '充值用户ID 系统充值则ID为0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_wallet
-- ----------------------------
DROP TABLE IF EXISTS `user_wallet`;
CREATE TABLE `user_wallet`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(12) NULL DEFAULT NULL,
  `wallet_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `wallet_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `balance` bigint(12) NULL DEFAULT NULL,
  `created_at` bigint(12) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '用户名',
  `email` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '邮件地址',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '电话号码',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '密码',
  `real_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '真实性名',
  `id_number` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '身份证',
  `active` tinyint(4) NULL DEFAULT NULL COMMENT '是否激活',
  `created_at` bigint(12) NULL DEFAULT NULL COMMENT '创建日期',
  `last_update` bigint(12) NULL DEFAULT NULL COMMENT '最后更新日期',
  `company_id` bigint(12) NULL DEFAULT NULL COMMENT '公司Id',
  `salt` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '盐',
  `wallet_address_id` bigint(12) NULL DEFAULT NULL COMMENT '钱包地址',
  `spread_id` bigint(20) NULL DEFAULT NULL COMMENT '推广用户',
  `is_realname_authentication` tinyint(1) NULL DEFAULT NULL COMMENT '是否实名认证',
  `last_login_ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '最后登录IP',
  `is_ban` tinyint(1) NULL DEFAULT NULL COMMENT '是否封号',
  `deleted` tinyint(1) NULL DEFAULT NULL COMMENT '是否删除',
  `mail_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '邮件验证',
  `phone_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '电话验证',
  `is_admin` tinyint(1) NULL DEFAULT NULL COMMENT '管理员标志',
  `is_private` tinyint(1) NULL DEFAULT NULL COMMENT '私有化部署用户',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5097 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users_key
-- ----------------------------
DROP TABLE IF EXISTS `users_key`;
CREATE TABLE `users_key`  (
  `id` bigint(12) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(12) NULL DEFAULT NULL,
  `access_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `security_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `created_at` bigint(255) NULL DEFAULT NULL COMMENT '创建时间',
  `deleted` bigint(1) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Procedure structure for GetModelsWithProviders
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetModelsWithProviders`;
delimiter ;;
CREATE PROCEDURE `GetModelsWithProviders`(IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255))
BEGIN
    -- 如果 order by 参数为空，给一个默认排序
    IF p_order_by IS NULL OR p_order_by = '' THEN
        SET p_order_by = 'mi.id ASC';
    END IF;

    -- 拼接查询语句
    SET @sql = CONCAT(
        'SELECT
            mi.id            AS info_id,
            mi.model_id      AS info_model_id,
            mi.node_id       AS info_node_id,
            mi.name          AS info_name,
            mi.api_version   AS info_api_version,
            mi.deploy_name   AS info_deploy_name,
            mi.input_price   AS info_input_price,
            mi.output_price  AS info_output_price,
            mi.cache_price   AS info_cache_price,
            mi.status        AS info_status,
            mi.last_update   AS info_last_update,

            mp.id            AS provider_id,
            mp.node_id       AS provider_node_id,
            mp.provider_id   AS provider_provider_id,
            mp.type          AS provider_type,
            mp.name          AS provider_name,
            mp.endpoint      AS provider_endpoint,
            mp.api_type      AS provider_api_type,
            mp.model_name    AS provider_model_name,
            mp.input_price   AS provider_input_price,
            mp.output_price  AS provider_output_price,
            mp.cache_price   AS provider_cache_price,
            mp.deleted       AS provider_deleted,
            mp.last_update   AS provider_last_update
        FROM models_info mi
        LEFT JOIN models_provider mp
            ON mi.model_id = mp.model_id
           AND mi.node_id = mp.node_id
        ORDER BY ', p_order_by,
        ' LIMIT ', p_limit, ' OFFSET ', p_offset
    );

    -- 预处理并执行主查询
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 统计总数（只统计 models_info）
    SET @count_query = 'SELECT COUNT(*) AS total_count FROM models_info';
    PREPARE stmt FROM @count_query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetPrivateUserModels
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetPrivateUserModels`;
delimiter ;;
CREATE PROCEDURE `GetPrivateUserModels`(IN p_user_id BIGINT,
    IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255))
BEGIN
    -- 临时表存放最终结果
    CREATE TEMPORARY TABLE IF NOT EXISTS tmp_private_user_models (
        id BIGINT,
        model_id VARCHAR(128),
        node_id VARCHAR(64),
        name VARCHAR(128),
        api_version VARCHAR(24),
        deploy_name VARCHAR(128),
        input_price INT,
        output_price INT,
        cache_price INT,
        status VARCHAR(12),
        last_update BIGINT(12),
        is_private TINYINT(1),
        owner_id BIGINT(12),
        address VARCHAR(256)
    );
TRUNCATE TABLE tmp_private_user_models;

-- 只查私有模型
INSERT INTO tmp_private_user_models
SELECT *
FROM models_info
WHERE owner_id = p_user_id AND is_private = 1;

-- 拼接 ORDER BY / LIMIT / OFFSET
SET @final_sql = CONCAT(
        'SELECT * FROM tmp_private_user_models ',
        IF(p_order_by IS NOT NULL AND p_order_by <> '', CONCAT(' ORDER BY ', p_order_by), ''),
        ' LIMIT ', p_limit,
        ' OFFSET ', p_offset
    );

PREPARE stmt FROM @final_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 总数
SET @count_query = 'SELECT COUNT(*) AS total_count FROM tmp_private_user_models';
PREPARE stmt FROM @count_query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

DROP TEMPORARY TABLE IF EXISTS tmp_private_user_models;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetUserInfos
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetUserInfos`;
delimiter ;;
CREATE PROCEDURE `GetUserInfos`(IN p_user_id BIGINT)
BEGIN

    SELECT
        u.user_name,
        u.email,
        u.phone,
        u.real_name,
        u.id_number,
        u.created_at,
        u.company_id,
        u.last_login_ip,
        u.spread_id,
        u.is_private,
        w.wallet_type,
        w.wallet_address,
        w.balance
    FROM
        users u
            LEFT JOIN user_wallet w ON u.id = w.user_id
    WHERE
        u.id = p_user_id;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetUserModels
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetUserModels`;
delimiter ;;
CREATE PROCEDURE `GetUserModels`(IN p_user_id BIGINT,
    IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255))
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE v_model_ids TEXT;

    -- 游标遍历 user_models_infos，只取 model_ids
    DECLARE cur CURSOR FOR
SELECT model_ids
FROM user_models_infos
WHERE user_id = p_user_id;

DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    -- 临时表存放最终结果
    CREATE TEMPORARY TABLE IF NOT EXISTS tmp_user_models (
        id BIGINT,
        model_id VARCHAR(128),
        node_id VARCHAR(64),
        name VARCHAR(128),
        api_version VARCHAR(24),
        deploy_name VARCHAR(128),
        input_price INT,
        output_price INT,
        cache_price INT,
        status VARCHAR(12),
        last_update BIGINT(12),
        is_private TINYINT(1),
        owner_id BIGINT(12),
        address VARCHAR(256)
    );
    TRUNCATE TABLE tmp_user_models;

    OPEN cur;
    read_loop: LOOP
            FETCH cur INTO v_model_ids;
            IF done THEN
                LEAVE read_loop;
    END IF;

            -- 动态拼接 SQL，只根据 model_ids 查 models_info
            SET @sql = CONCAT(
                'INSERT INTO tmp_user_models ',
                'SELECT * FROM models_info ',
                'WHERE FIND_IN_SET(model_id, ''', v_model_ids, ''') > 0'
            );

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    END LOOP;

    CLOSE cur;

    -- 拼接 ORDER BY / LIMIT / OFFSET
    SET @final_sql = CONCAT(
            'SELECT * FROM tmp_user_models ',
            IF(p_order_by IS NOT NULL AND p_order_by <> '', CONCAT(' ORDER BY ', p_order_by), ''),
            ' LIMIT ', p_limit,
            ' OFFSET ', p_offset
        );

    PREPARE stmt FROM @final_sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 总数
    SET @count_query = 'SELECT COUNT(*) AS total_count FROM tmp_user_models';
    PREPARE stmt FROM @count_query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    DROP TEMPORARY TABLE IF EXISTS tmp_user_models;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for ListNodeUserNodeInfos
-- ----------------------------
DROP PROCEDURE IF EXISTS `ListNodeUserNodeInfos`;
delimiter ;;
CREATE PROCEDURE `ListNodeUserNodeInfos`(IN p_node_user_id BIGINT,   -- 目标用户ID（node_users.id）
    IN p_page         INT,      -- 页码，从1开始
    IN p_page_size    INT,      -- 每页条数
    IN p_order_by     VARCHAR(64))
BEGIN
    DECLARE v_page INT DEFAULT 1;
    DECLARE v_page_size INT DEFAULT 20;
    DECLARE v_offset INT DEFAULT 0;
    DECLARE v_order_by_sql VARCHAR(64);

    -- 参数处理
    IF p_page IS NOT NULL AND p_page > 0 THEN
        SET v_page = p_page;
    END IF;

    IF p_page_size IS NOT NULL AND p_page_size > 0 THEN
        SET v_page_size = p_page_size;
    END IF;

    IF v_page_size > 1000 THEN
        SET v_page_size = 1000;
    END IF;

    SET v_offset = (v_page - 1) * v_page_size;

    -- 设置排序参数（防SQL注入，仅允许合法字段）
    IF p_order_by IS NULL OR p_order_by = '' THEN
        SET v_order_by_sql = 'n.created_at DESC';
    ELSE
        -- 简单校验：只允许字母数字下划线和空格
        IF p_order_by REGEXP '^[a-zA-Z0-9_\\.\\s]+(ASC|DESC)?$' THEN
            SET v_order_by_sql = p_order_by;
        ELSE
            SET v_order_by_sql = 'n.created_at DESC';
        END IF;
    END IF;

    -- 构建动态 SQL
    SET @sql = CONCAT(
        'SELECT ',
        'n.id AS node_id, ',
        'n.ids AS node_code, ',
        'n.node_user_id AS node_user_id, ',
        'n.created_at AS node_created_at, ',
        'n.lastupdate_at AS node_lastupdate_at, ',
        'n.domain AS node_domain, ',
        'u.id AS user_id, ',
        'u.email AS user_email, ',
        'u.created_at AS user_created_at, ',
        'u.last_update AS user_last_login, ',
        'u.is_active AS user_is_active, ',
        'u.company_id AS user_company_id ',
        'FROM nodes AS n ',
        'INNER JOIN node_users AS u ON u.id = n.node_user_id ',
        'WHERE u.id = ? ',
        'ORDER BY ', v_order_by_sql, ' ',
        'LIMIT ? OFFSET ?'
    );

    PREPARE stmt FROM @sql;
    EXECUTE stmt USING p_node_user_id, v_page_size, v_offset;
    DEALLOCATE PREPARE stmt;

    -- 第二个结果集：总数
    SELECT COUNT(*) AS total
    FROM nodes AS n
    INNER JOIN node_users AS u
        ON u.id = n.node_user_id
    WHERE u.id = p_node_user_id;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
