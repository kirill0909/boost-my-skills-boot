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
	UPDATE bot.users SET active = TRUE WHERE tg_chat_id = $1
	`
)
