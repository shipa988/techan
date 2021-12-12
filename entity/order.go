package entity

import (
	"github.com/sdcoffey/big"
	"time"
)

// OrderSide is a simple enumeration representing the side of an Order (buy or sell)
type OrderSide int

// BUY and SELL enumerations
const (
	BUY OrderSide = iota
	SELL
)

// Order represents a trade execution (buy or sell) with associated metadata.
type Order struct {
	OrderID  uint64      `json:"order_id"`
	ClientId uint64      `json:"client_id"`
	Created  time.Time   `json:"created"`
	Type     OrderSide   `json:"type"`
	Pair     string      `json:"pair"`
	Price    big.Decimal `json:"price"`
	Quantity big.Decimal `json:"quantity"`
	Amount   big.Decimal `json:"amount"`
}

type StopOrder struct {
	ParentOrderID string      `json:"parent_order_id"`
	OrderID       uint64      `json:"order_id"`
	ClientId      uint64      `json:"client_id"`
	Created       time.Time   `json:"created"`
	Type          OrderSide   `json:"type"`
	Pair          string      `json:"pair"`
	TriggerPrice  big.Decimal `json:"trigger_price"`
	Quantity      big.Decimal `json:"quantity"`
	Amount        big.Decimal `json:"amount"`
}
