package repository

const (
	queryGetUUID = `
		INSERT INTO users.user DEFAULT VALUES RETURNING tg_uuid
  `

	queryIsAdmin = `
  SELECT
 CASE
     WHEN u.name = $1 AND tg_chat_id = $2 THEN TRUE
     ELSE FALSE
 END AS result
 FROM users.user u
  `

	queryUserActivation = `
  UPDATE users.user SET name = $1, tg_chat_id = $2, active = true WHERE tg_uuid = $3
  `
	querySetUpBackendDirection = `
  UPDATE users.user SET direction_id = 1 WHERE tg_chat_id = $1
  `
)
