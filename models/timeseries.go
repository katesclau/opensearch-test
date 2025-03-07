package models

import "time"

type TimeseriesString struct {
	TS  time.Time `json:"ts"`
	Val string    `json:"val"`
}
