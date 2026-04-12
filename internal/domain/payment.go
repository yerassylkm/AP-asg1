package domain

type Payment struct {
	ID            string `json:"id"             db:"id"`
	OrderID       string `json:"order_id"       db:"order_id"`
	TransactionID string `json:"transaction_id" db:"transaction_id"`
	Amount        int64  `json:"amount"         db:"amount"`
	Status        string `json:"status"         db:"status"`
}

const (
	StatusAuthorized = "Authorized"
	StatusDeclined   = "Declined"
)
