package repository

const (
	querySetStatusActive = `
	UPDATE users SET active = TRUE, tg_chat_id = $1, tg_name = $3, updated_at = NOW()
	WHERE tg_uuid = $2 AND active IS FALSE
	`

	queryGetMainButtons = `
	SELECT
	name AS name,
	only_for_admin AS only_for_admin
	FROM main_buttons ORDER BY name;
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
	queryGetUserDirection = `
	SELECT
    d.id AS id,
    d.direction AS direction,
    COALESCE(d.parent_direction_id, 0) AS parent_directon_id,
    CAST(EXTRACT(EPOCH FROM d.created_at) AS BIGINT) AS created_at,
    CAST(COALESCE(EXTRACT (EPOCH FROM d.updated_at), 0) AS BIGINT) AS updated_at
FROM directions d
INNER JOIN users u ON u.id = d.user_id
WHERE
	u.tg_chat_id = $1 AND
	CASE
		WHEN $2::INTEGER = 0 THEN (d.parent_direction_id IS NULL)
		ELSE $2::INTEGER = d.parent_direction_id
	END ORDER BY direction;
	`

	queryCreateDirection = `
	INSERT INTO directions (direction, user_id, parent_direction_id)
	VALUES ($1, (SELECT id FROM users WHERE tg_chat_id = $2), $3)
	RETURNING direction;
	`

	querySaveQuestion = `
	INSERT INTO infos(question, direction_id) VALUES ($1, $2) RETURNING id;
	`

	querySaveAnswer = `
	UPDATE infos SET answer = $1, updated_at = NOW() WHERE id = $2;
	`
	queryGetQuestionsByDirectionID = `
	SELECT
		id AS id,
		question AS text
		FROM infos WHERE direction_id = $1 ORDER BY created_at;
	`
	queryGetAnswerByInfoID = `
	SELECT answer FROM infos WHERE id = $1;
	`
)
