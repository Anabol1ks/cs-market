package inventory

import "time"

type Skin struct {
	MarketHashName string    `json:"market_hash_name" gorm:"primaryKey"`
	MinPrice       *float64  `json:"min_price"`
	AvgPrice       *float64  `json:"mean_price"`
	MaxPrice       *float64  `json:"max_price"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
