DROP TABLE IF EXISTS friendships;
DROP TYPE IF EXISTS friend_status;
DROP FUNCTION IF EXISTS send_friend_request(INT, INT);
DROP FUNCTION IF EXISTS accept_friend_request(INT, INT);
DROP FUNCTION IF EXISTS cancel_friend_request(INT, INT);
