DELIMITER //

DROP PROCEDURE IF EXISTS ListNodeModelsProviderInfos //

CREATE PROCEDURE ListNodeModelsProviderInfos(
    IN p_node_user_id BIGINT,
    IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255)
)
BEGIN
    -- 设置默认排序
    SET @order_clause = IF(p_order_by IS NOT NULL AND p_order_by <> '',
                          CONCAT(' ORDER BY ', p_order_by),
                          ' ORDER BY nm.id DESC');

    -- 构建主查询 SQL
    -- 1. 从 nodes 表查询用户的所有节点（通过 owner_id 关联）
    -- 2. 通过节点ID查询 node_models_info_maps 表，获取节点、模型、供应商的映射关系
    -- 3. 关联 models_info 表获取模型详细信息
    -- 4. 关联 models_provider 表获取供应商详细信息
    SET @query_sql = CONCAT(
        'SELECT ',
        '  nm.id AS map_id, ',
        '  nm.node_id, ',
        '  nm.node_name, ',
        '  nm.model_id, ',
        '  nm.model_provider_id, ',
        '  nm.created_at AS map_created_at, ',
        '  nm.updated_at AS map_updated_at, ',
        '  n.id AS node_pk_id, ',
        '  n.name AS node_full_name, ',
        '  n.owner_id, ',
        '  n.domain, ',
        '  n.company_id AS node_company_id, ',
        '  n.created_at AS node_created_at, ',
        '  n.lastupdate_at AS node_lastupdate_at, ',
        '  mi.id AS model_pk_id, ',
        '  mi.model_id AS model_identifier, ',
        '  mi.name AS model_name, ',
        '  mi.api_version, ',
        '  mi.deploy_name, ',
        '  mi.input_price AS model_input_price, ',
        '  mi.output_price AS model_output_price, ',
        '  mi.cache_price AS model_cache_price, ',
        '  mi.status AS model_status, ',
        '  mi.address AS model_address, ',
        '  mi.api_styles, ',
        '  mp.id AS provider_pk_id, ',
        '  mp.provider_id, ',
        '  mp.type AS provider_type, ',
        '  mp.name AS provider_name, ',
        '  mp.endpoint, ',
        '  mp.api_type, ',
        '  mp.model_name AS provider_model_name, ',
        '  mp.input_price AS provider_input_price, ',
        '  mp.output_price AS provider_output_price, ',
        '  mp.cache_price AS provider_cache_price, ',
        '  mp.api_keys ',
        'FROM node_models_info_maps nm ',
        'INNER JOIN nodes n ON nm.node_id = n.id ',
        'LEFT JOIN models_info mi ON nm.model_id = mi.id ',
        'LEFT JOIN models_provider mp ON nm.model_provider_id = mp.id ',
        'WHERE 1=1 ',
        IF(p_node_user_id IS NOT NULL AND p_node_user_id > 0,
           CONCAT('AND n.owner_id = ', p_node_user_id, ' '),
           ''),
        @order_clause,
        ' LIMIT ', p_limit,
        ' OFFSET ', p_offset
    );

    -- 执行主查询
    PREPARE stmt FROM @query_sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 查询总数
    SET @count_sql = CONCAT(
        'SELECT COUNT(*) AS total_count ',
        'FROM node_models_info_maps nm ',
        'INNER JOIN nodes n ON nm.node_id = n.id ',
        'WHERE 1=1 ',
        IF(p_node_user_id IS NOT NULL AND p_node_user_id > 0,
           CONCAT('AND n.owner_id = ', p_node_user_id),
           '')
    );

    PREPARE stmt FROM @count_sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END //

DELIMITER ;
