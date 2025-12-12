package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MorStrategy struct {
	cfg       *MORSConfig
	lastState TrendState
	windows   []float64
	emaScores []float64

	getPosition PositionManagerGet

	curvatureFilter *CurvatureFilter
}

func NewMorStrategy(cfg *MORSConfig) *MorStrategy {
	return &MorStrategy{
		cfg:             cfg,
		emaScores:       []float64{},
		windows:         make([]float64, cfg.windowsSize),
		curvatureFilter: NewCurvatureFilter(300),
	}
}

func (s *MorStrategy) updateWindow(v float64) {
	s.windows = append(s.windows, v)
	if len(s.windows) > s.cfg.windowsSize {
		s.windows = s.windows[1:]
	}
}

func (s *MorStrategy) updateEMA(rawScore float64) float64 {
	alpha := s.cfg.emaAlpha

	if len(s.emaScores) == 0 {
		// 第一个点，直接用原始 score 初始化
		s.emaScores = append(s.emaScores, rawScore)
		return rawScore
	}
	prev := s.emaScores[len(s.emaScores)-1]
	now := alpha*rawScore + (1-alpha)*prev
	s.emaScores = append(s.emaScores, now)
	return now
}

func (s *MorStrategy) classifyWindow() TrendState {
	n := len(s.windows)
	if n == 0 || n < s.cfg.MinWindowSamples {
		return StateNeutral
	}

	var sum float64
	var pos, neg int
	for _, v := range s.windows {
		sum += v
		if v > 0 {
			pos++
		} else if v < 0 {
			neg++
		}
	}

	avg := sum / float64(n)
	posRatio := float64(pos) / float64(n)
	negRatio := float64(neg) / float64(n)

	// 阈值逻辑：你可以根据自己 Python 的实验微调
	if avg > s.cfg.StrongUpAvg && posRatio >= s.cfg.PosRatioStrong {
		return StateStrongUp
	} else if avg > s.cfg.UpAvg && posRatio >= s.cfg.PosRatioUp {
		return StateUp
	} else if avg < s.cfg.StrongDownAvg && negRatio >= s.cfg.NegRatioStrong {
		return StateStrongDown
	} else if avg < s.cfg.DownAvg && negRatio >= s.cfg.NegRatioDown {
		return StateDown
	}
	return StateNeutral
}

func (s *MorStrategy) OnDataUpdate(p ...TickDataInterface) (*Signal, error) {
	// 1. 提取 score
	score := 0.0
	hasScore := false
	for _, v := range p {
		if strings.ToLower(v.Name()) == "score" {
			score = v.Value()
			hasScore = true
			break
		}
	}
	if !hasScore {
		// 没有 score，直接无信号
		return nil, nil
	}

	// 2. 更新 EMA & 窗口
	emaScore := s.updateEMA(score)
	s.updateWindow(emaScore)

	// 3. 所有点都喂给曲率过滤器，让它自己决定是否要用
	s.curvatureFilter.Append(emaScore)

	// 4. 计算当前趋势状态
	state := s.classifyWindow()

	// 统一的 meta
	meta := map[string]float64{
		"score":     score,
		"ema_score": emaScore,
	}

	// 5. 根据状态 + 曲率过滤器生成信号
	switch state {
	case StateNeutral, StateUp, StateDown:
		// 趋势不强 → 一律无信号
		return nil, nil

	case StateStrongUp:
		// strong_up 区间：只在“最佳切入点”才做空
		if s.curvatureFilter.IsBestPoint(true) {
			return &Signal{
				SignalType: ShortEntry,
				Ts:         time.Now().UnixMilli(),
				Strength:   5,
				Reason:     "strong_up + curvature best point",
				Meta:       meta,
				Id:         uuid.New().String(),
			}, nil
		}
		return nil, nil

	case StateStrongDown:
		// strong_down 区间：只在“最佳切入点”才做多
		if s.curvatureFilter.IsBestPoint(false) {
			return &Signal{
				SignalType: LongEntry,
				Ts:         time.Now().UnixMilli(),
				Strength:   5,
				Reason:     "strong_down + curvature best point",
				Meta:       meta,
				Id:         uuid.New().String(),
			}, nil
		}
		return nil, nil

	default:
		// 兜底：未知状态直接无信号
		return nil, nil
	}
}
