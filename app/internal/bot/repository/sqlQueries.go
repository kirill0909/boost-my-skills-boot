package repository

const (
	queryCompareUUID = `
	SELECT 
	  CASE
	    WHEN $2 = (SELECT tg_uuid FROM bot.users WHERE tg_chat_id = $1) THEN TRUE
	    ELSE FALSE
	  END
	`

	querySetStatusActive = `
	UPDATE bot.users SET active = TRUE, tg_chat_id = $1, tg_name = $3 
	WHERE tg_uuid = $2 AND active IS FALSE
	`
)
