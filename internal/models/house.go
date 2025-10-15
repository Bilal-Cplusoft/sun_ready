package models

import (
	"time"
)


type House struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Diameter   float64 `json:"diameter"`
	Probability float64 `json:"probability"`
	State      string  `json:"state"`
	TileID     *uint   `json:"tile_id,omitempty"`
}
