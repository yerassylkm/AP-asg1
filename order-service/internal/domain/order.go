package domain

import "time"

type Order struct {
	ID            string    `json:"id"             db:"id"`
	CustomerID    string    `json:"customer_id"    db:"customer_id"`
	CustomerEmail string    `json:"customer_email" db:"customer_email"`
	ItemName      string    `json:"item_name"      db:"item_name"`
	Amount        int64     `json:"amount"         db:"amount"`
	Status        string    `json:"status"         db:"status"`
	CreatedAt     time.Time `json:"created_at"     db:"created_at"`
}

const (
	StatusPending   = "Pending"
	StatusPaid      = "Paid"
	StatusFailed    = "Failed"
	StatusCancelled = "Cancelled"
)
