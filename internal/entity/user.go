package entity

import "time"

type LimitType int

const DefaultUserTtl = time.Hour * 24

const (
	NoLimit LimitType = iota
	SpeedOnly
	TrafficOnly
	LimitSpeedAfterTrafficReached
)

type User struct {
	ID           uint32
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	TelegramId   int64     `json:"telegramId"`
	LimitSpeed   uint32    `json:"limitSpeed"`
	LimitTraffic uint32    `json:"limitTraffic"`
	LimitType    LimitType `json:"limitType"`
	Active       bool      `json:"active"`
	Ttl          time.Time `json:"ttl"`
	Updated      time.Time `json:"updated"`
}
