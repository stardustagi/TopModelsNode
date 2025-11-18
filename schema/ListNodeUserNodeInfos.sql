DELIMITER //

DROP PROCEDURE IF EXISTS ListNodeUserNodeInfos //

CREATE PROCEDURE ListNodeUserNodeInfos(
    IN p_node_user_id BIGINT,
    IN p_limit INT,
    IN p_offset INT,
    IN p_order_by VARCHAR(255)
)
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
        node_id VARCHAR(64),
        node_user_id
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
END //

DELIMITER ;
