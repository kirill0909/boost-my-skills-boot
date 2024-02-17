package repository

const (
	querySetStatusActive = `
	UPDATE users SET active = TRUE, tg_chat_id = $1, tg_name = $3 
	WHERE tg_uuid = $2 AND active IS FALSE
	`

	queryGetMainButtons = `
	SELECT
	name AS name,
	only_for_admin AS only_for_admin
	FROM main_buttons
	`

	queryGetActiveUsers = `
	SELECT 
	tg_chat_id AS tg_chat_id, 
	is_admin AS is_admin
	FROM users WHERE active
	`

	queryGetUpdatedButtons = `
	SELECT
    name AS name,
    only_for_admin AS only_for_admin
FROM main_buttons
            WHERE EXTRACT(EPOCH FROM created_at) > $1 OR
                    EXTRACT(EPOCH FROM updated_at) > $1;
	`

	querySetUserActive = `
	UPDATE users SET 
	active = TRUE
	, tg_chat_id = $1
	, tg_name = $3
	, updated_at = NOW()
	WHERE tg_uuid = $2 AND active IS FALSE
	`
)
