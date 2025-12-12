package main

type SignalType string

const (
	ShortEntry SignalType = "short_entry"
	ShortExit  SignalType = "short_exit"
	LongEntry  SignalType = "long_entry"
	LongExit   SignalType = "long_exit"
	FLAT       SignalType = "flat"
	None       SignalType = "none"
)

type Signal struct {
	SignalType SignalType         `json:"signal_type"`
	Ts         int64              `json:"ts"` // 毫秒级时间戳
	Strength   float64            `json:"strength"`
	Reason     string             `json:"reason"`
	Meta       map[string]float64 `json:"meta"`
	Id         string             `json:"id"`
}
