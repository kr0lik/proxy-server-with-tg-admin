package entity

import "time"

type UserStat struct {
	UserID          uint32
	TrafficInDay    uint64    `json:"trafficInDay"`
	TrafficOutDay   uint64    `json:"trafficOutDay"`
	TrafficInTotal  uint64    `json:"trafficInTotal"`
	TrafficOutTotal uint64    `json:"trafficOutTotal"`
	DaysActive      uint      `json:"daysActive"`
	Updated         time.Time `json:"updated"`
}
