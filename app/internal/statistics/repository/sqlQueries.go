package repository

const (
	queryGetStatistics = `
	SELECT COUNT(1) FROM infos WHERE EXTRACT(EPOCH FROM created_at) between $1 AND $2;
	`
)
