package main

type TrendState string

const (
	StateStrongUp   TrendState = "strong_up"
	StateUp         TrendState = "up"
	StateNeutral    TrendState = "neutral"
	StateDown       TrendState = "down"
	StateStrongDown TrendState = "strong_down"
)
