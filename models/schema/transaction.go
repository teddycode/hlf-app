package schema

type TransactionSwag struct {
	Timestamp int    `json:"timestamp" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
	Point     string `json:"point" binding:"required"`
}

type QueryTrancsNumSwag struct {
	Type string `json:"type" binding:"required"`
}