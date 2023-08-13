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
	querySetUpFrontedDirection = `
  UPDATE users.user SET direction_id = 2 WHERE tg_chat_id = $1
  `

	queryGetRandomQuestion = `
	SELECT q.id, q.question, q.code
 FROM users.questions q
 INNER JOIN users.user u ON u.direction_id = q.direction_id
 WHERE u.tg_chat_id = $1
 ORDER BY RANDOM()
 LIMIT 1
	`
)
