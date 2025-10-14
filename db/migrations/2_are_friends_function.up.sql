CREATE OR REPLACE FUNCTION are_friends(p_user_id_1 INT, p_user_id_2 INT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM friendships
        WHERE
            user_id_1 = LEAST(p_user_id_1, p_user_id_2) AND
            user_id_2 = GREATEST(p_user_id_1, p_user_id_2) AND
            status = 'accepted'
    );
END;
$$ LANGUAGE plpgsql;
