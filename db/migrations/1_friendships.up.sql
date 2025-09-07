CREATE TYPE friend_status AS ENUM ('pending', 'accepted');

CREATE TABLE friendships (
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    status friend_status NOT NULL DEFAULT 'pending',
    sent_by_user_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT chk_canonical_id_order CHECK (user_id_1 < user_id_2),

    PRIMARY KEY (user_id_1, user_id_2),

    FOREIGN KEY (user_id_1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id_2) REFERENCES users(id) ON DELETE CASCADE,
    
    FOREIGN KEY (sent_by_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION send_friend_request(p_sender_id INT, p_receiver_id INT)
RETURNS TEXT AS $$
DECLARE
    v_user_1 INT;
    v_user_2 INT;
BEGIN
    IF p_sender_id = p_receiver_id THEN
        RETURN 'Error: Cannot send a friend request to yourself.';
    END IF;

    IF p_sender_id < p_receiver_id THEN
        v_user_1 := p_sender_id;
        v_user_2 := p_receiver_id;
    ELSE
        v_user_1 := p_receiver_id;
        v_user_2 := p_sender_id;
    END IF;

    IF EXISTS (SELECT 1 FROM friendships WHERE user_id_1 = v_user_1 AND user_id_2 = v_user_2) THEN
        RETURN 'Error: Friendship already exists or a request is pending.';
    END IF;

    INSERT INTO friendships (user_id_1, user_id_2, sent_by_user_id, status)
    VALUES (v_user_1, v_user_2, p_sender_id, 'pending');

    RETURN 'Friend request sent successfully.';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION accept_friend_request(p_responder_id INT, p_sender_id INT)
RETURNS TEXT AS $$
DECLARE
    v_user_1 INT;
    v_user_2 INT;
    v_friendship friendships%ROWTYPE;
BEGIN
    IF p_responder_id = p_sender_id THEN
        RETURN 'Error: Cannot accept a friend request from yourself.';
    END IF;

    IF p_responder_id < p_sender_id THEN
        v_user_1 := p_responder_id;
        v_user_2 := p_sender_id;
    ELSE
        v_user_1 := p_sender_id;
        v_user_2 := p_responder_id;
    END IF;

    SELECT * INTO v_friendship FROM friendships WHERE user_id_1 = v_user_1 AND user_id_2 = v_user_2;

    IF v_friendship IS NULL THEN
        RETURN 'Error: No pending friend request found.';
    END IF;

    IF v_friendship.sent_by_user_id = p_responder_id THEN
        RETURN 'Error: Cannot accept a request you sent.';
    END IF;

    IF v_friendship.status = 'accepted' THEN
        RETURN 'Error: Friendship already exists.';
    END IF;

    UPDATE friendships
    SET status = 'accepted', updated_at = NOW()
    WHERE user_id_1 = v_user_1 AND user_id_2 = v_user_2;
    RETURN 'Friend request accepted.';
END;
$$ LANGUAGE plpgsql;

-- helper function to cancel/delete friendship by two user ids (order-insensitive)
CREATE OR REPLACE FUNCTION cancel_friend_request(p_user_a INT, p_user_b INT)
RETURNS TEXT AS $$
DECLARE
    v_user_1 INT;
    v_user_2 INT;
    v_rows INT;
BEGIN
    IF p_user_a = p_user_b THEN
        RETURN 'Error: Cannot cancel friendship with yourself.';
    END IF;

    IF p_user_a < p_user_b THEN
        v_user_1 := p_user_a;
        v_user_2 := p_user_b;
    ELSE
        v_user_1 := p_user_b;
        v_user_2 := p_user_a;
    END IF;

    DELETE FROM friendships WHERE user_id_1 = v_user_1 AND user_id_2 = v_user_2;
    GET DIAGNOSTICS v_rows = ROW_COUNT;

    IF v_rows = 0 THEN
        RETURN 'No friendship or pending request found.';
    END IF;

    RETURN 'Friendship/request cancelled.';
END;
$$ LANGUAGE plpgsql;
