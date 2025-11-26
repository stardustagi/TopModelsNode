DELIMITER //

DROP PROCEDURE IF EXISTS ListNodeUserNodeInfos //

CREATE PROCEDURE ListNodeUserNodeInfos(
    IN p_node_user_id BIGINT,
    IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255)
)
BEGIN
    -- 查询节点和节点用户的关联信息
    -- 如果 p_order_by 为空，默认按 nodes.id 降序排序
    SET @order_clause = IF(p_order_by IS NOT NULL AND p_order_by <> '',
                          CONCAT(' ORDER BY ', p_order_by),
                          ' ORDER BY n.id DESC');

    -- 构建主查询 SQL，关联 nodes 和 node_users 表
    SET @query_sql = CONCAT(
        'SELECT ',
        '  n.id AS node_id, ',
        '  n.name AS node_name, ',
        '  n.owner_id, ',
        '  n.created_at AS node_created_at, ',
        '  n.lastupdate_at AS node_lastupdate_at, ',
        '  n.domain, ',
        '  n.access_key, ',
        '  n.security_key, ',
        '  n.company_id AS node_company_id, ',
        '  nu.id AS user_id, ',
        '  nu.email, ',
        '  nu.created_at AS user_created_at, ',
        '  nu.deleted, ',
        '  nu.last_update AS user_last_update, ',
        '  nu.is_active, ',
        '  nu.company_id AS user_company_id ',
        'FROM nodes n ',
        'LEFT JOIN node_users nu ON n.owner_id = nu.id ',
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
        'FROM nodes n ',
        'LEFT JOIN node_users nu ON n.owner_id = nu.id ',
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
