package domain

type Payment struct {
	ID            string `db:"id"`
	OrderID       string `db:"order_id"`
	TransactionID string `db:"transaction_id"`
	Amount        int64  `db:"amount"`
	Status        string `db:"status"`
	CustomerEmail string `db:"customer_email"`
}

const (
	StatusAuthorized = "Authorized"
	StatusDeclined   = "Declined"
)
