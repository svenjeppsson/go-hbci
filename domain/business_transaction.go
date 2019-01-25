package domain

import "github.com/shopspring/decimal"

// Amount represents a value associated with a currency
type Amount struct {
	Amount   decimal.Decimal
	Currency string
}

// BusinessTransaction provides information about a transaction and whether
// there is a signature needed or not
type BusinessTransaction struct {
	ID               string
	NeededSignatures int
	Limit            *AccountLimit
}
