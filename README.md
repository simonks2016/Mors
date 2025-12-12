# MORS: A Microstructure-Based Reversal Signal Engine

MORS is a lightweight trading signal engine built on market microstructure and transient impact decay theory.
It identifies short-cycle reversal points by combining a momentum score with a local curvature filter, detecting the exact moment when order-flow momentum weakens.

## ⚠️ MORS Does Not Manage Positions or Orders

MORS focuses only on generating high-quality reversal signals.
Position sizing, risk management, and execution should be handled by your own trading engine.

## 🔍 How It Works
### 1. Momentum Score

A real-time score derived from order flow, price movement, or other microstructure-driven features.

### 2. Trend Classification

A sliding window classifies the local trend state into five categories:

``` go 
    StateStrongUp   TrendState = "strong_up"
	StateUp         TrendState = "up"
	StateNeutral    TrendState = "neutral"
	StateDown       TrendState = "down"
	StateStrongDown TrendState = "strong_down""
```
This determines whether the system should look for potential short or long reversal points.

### 3. Curvature Filter

1. The core innovation of MORS.
2. It analyzes local curvature by comparing:
3. Backward 5-point average
4. Forward 5-point average

This detects momentum decay, a critical sign of microstructure overreaction exhaustion.
When curvature bends in the opposite direction, a high-confidence reversal opportunity emerges.

### 4. Reversal Signals

MORS generates signals only at statistically meaningful exhaustion points:

ShortEntry → when strong_up + curvature bending downward

LongEntry → when strong_down + curvature bending upward

These represent direction changes, not trade orders.

## 📡 Signal Output Example
``` go 
type Signal struct {
    SignalType SignalType
    Ts         int64
    Strength   float64
    Reason     string
    Meta       map[string]float64
}
```

### Example Meaning:

* Short Entry: trend is strong_up and curvature indicates momentum exhaustion

* Long Entry: trend is strong_down and curvature shows reversal forming


## 🧠 Design Philosophy

* Strategy-layer only → MORS does not manage positions

* Composable → can integrate with any trading engine or execution layer

* Microstructure-driven → signals reflect real-time energy decay

* Robust → avoids noise, avoids trend traps, filters false moves


## 🚀 Why It Works

MORS captures microstructure overreaction, a known phenomenon where short-term order-flow imbalance pushes price too far before snapping back.
The curvature filter isolates the impact decay moment, producing exceptionally clean, high-quality reversal entries.

In testing, this drastically reduces noise and improves profitability over naive reversal logic.

## 📦 Coming Soon

Backtesting utilities

Go code examples

Visualization tools

Real-time streaming integration

Plug-and-play signal dashboards

## ✍️ Author
* Authored by Simon Liang, Dec 12, 2025. 
