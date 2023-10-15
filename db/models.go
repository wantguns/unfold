package db

import (
	"time"
)

type Transactions struct {
	UUID           string `gorm:"primaryKey;"`
	Amount         float64
	CurrentBalance float64
	Timestamp      time.Time `gorm:"index:sortTimestamp,sort:desc"`
	Type           string
	Account        string
	Merchant       string
}
