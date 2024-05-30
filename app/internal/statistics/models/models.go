package models

type GetStatisticsRequest struct {
	DateFrom int64
	DateTo   int64
}

type GetStatisticsResult struct {
	InfosAdded int64
}
