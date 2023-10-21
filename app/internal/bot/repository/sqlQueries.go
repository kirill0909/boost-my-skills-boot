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
 FROM users.user u WHERE u.id = 1
  `

	queryUserActivation = `
  UPDATE users.user SET name = $1, tg_chat_id = $2, active = true WHERE tg_uuid = $3 AND active IS FALSE
  `
	querySetUpDirection = `
  UPDATE users.user SET direction_id = $1 WHERE tg_chat_id = $2
  `
	queryGetRandomQuestion = `
	SELECT
   ui.id,
   ui.question
   FROM users.info ui
   INNER JOIN users.user uu ON uu.direction_id = ui.direction_id
   INNER JOIN users.sub_directions usd ON usd.id = ui.sub_direction_id
   LEFT JOIN users.sub_sub_directions ussd ON ussd.sub_direction_id =  ui.sub_sub_direction_id
   WHERE
   uu.tg_chat_id = $1::BIGINT
   AND usd.id = $2::INTEGER
   AND ($3::INTEGER = 0 OR ui.sub_sub_direction_id = $3)
   ORDER BY RANDOM()
   LIMIT 1
	`

	queryGetSubdirectons = `
	SELECT 
	usd.sub_direction 
	FROM users.sub_directions usd
 INNER JOIN users.user uu ON uu.direction_id = usd.direction_id
 WHERE uu.tg_chat_id = $1
 ORDER BY usd.id
	`
	queryGetSubSubdirectons = `
	SELECT
  ussd.sub_sub_direction
  FROM users.sub_sub_directions ussd
  INNER JOIN users.sub_directions usd ON usd.id = ussd.direction_id
  INNER JOIN users.user uu ON uu.direction_id = ussd.direction_id
  WHERE
  sub_direction_id = $1
  AND uu.tg_chat_id = $2
  ORDER BY ussd.id
	`

	queryGetAnswer = `
	SELECT answer FROM users.info WHERE id = $1
	`

	querySaveQuestion = `
	INSERT INTO users.info (direction_id, sub_direction_id, sub_sub_direction_id, question, answer)
 VALUES (
 (SELECT direction_id FROM users.user WHERE tg_chat_id = $1),
 $2,
 (CASE WHEN $3 = 0 THEN NULL ELSE $3 END), 
 $4,
 '') RETURNING id
	`

	querySaveAnswer = `
	UPDATE users.info SET answer = $1 WHERE id = $2
	`
	queryGetDirectionsInfo = `
	SELECT
 ud.id AS direction_id,
 ud.direction AS direction_name
 FROM users.directions ud
 ORDER BY ud.id
	`
	queryGetSubdirectionsInfo = `
	SELECT
	usd.direction_id AS direction_id,
  usd.id AS sub_direction_id,
  usd.sub_direction AS sub_direction_name
  FROM users.sub_directions usd
  ORDER BY usd.id
	`

	queryGetSubSubdirectionsInfo = `
	 SELECT
	ussd.direction_id AS direction_id,
	ussd.sub_direction_id AS sub_direction_id,
  ussd.id AS sub_sub_direction_id,
  ussd.sub_sub_direction AS sub_sub_direction_name
  FROM users.sub_sub_directions ussd
  ORDER BY ussd.id
	`

	queryGetDirectionByChatID = `
	SELECT direction_id FROM users.user WHERE tg_chat_id = $1
	`

	queryPrintInfo = `
	SELECT
 ui.id AS id,
 ui.question AS question,
 ui.answer AS answer
 FROM users.info ui
 INNER JOIN users.user uu ON uu.tg_chat_id = $1
 WHERE ui.sub_direction_id = $2 AND ui.sub_sub_direction_id = $3
	`
)
