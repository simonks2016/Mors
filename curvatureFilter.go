package main

import "sort"

type CurvatureFilter struct {
	windows    []float64
	windowSize int
}

func NewCurvatureFilter(windowSize int) *CurvatureFilter {

	return &CurvatureFilter{
		windows:    make([]float64, windowSize),
		windowSize: windowSize,
	}
}

// Append 添加新的分数
func (cf *CurvatureFilter) Append(score float64) {

	cf.windows = append(cf.windows, score)
	if len(cf.windows) > cf.windowSize {
		cf.windows = cf.windows[1:]
	}
}

// mean 均值
func mean(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}

// IsBestPoint 核心曲率检测函数
func (cf *CurvatureFilter) IsBestPoint(IsSeekTopPoint bool) bool {

	if len(cf.windows) < cf.windowSize {
		return false
	}

	side := func() string {

		if IsSeekTopPoint {
			return "up"
		}
		return "down"
	}()

	ws := cf.windows
	n := len(ws)

	// 1. 前5 后5 均值
	if n < 10 {
		return false
	}
	prev := ws[n-10 : n-5]
	last := ws[n-5:]

	mPrev := mean(prev)
	mLast := mean(last)

	// diff = mLast - mPrev
	diff := mLast - mPrev

	// 2. 形态过滤：不创新高 / 不创新低
	prevMax := prev[0]
	for _, v := range prev {
		if v > prevMax {
			prevMax = v
		}
	}
	lastMax := last[0]
	for _, v := range last {
		if v > lastMax {
			lastMax = v
		}
	}

	prevMin := prev[0]
	for _, v := range prev {
		if v < prevMin {
			prevMin = v
		}
	}
	lastMin := last[0]
	for _, v := range last {
		if v < lastMin {
			lastMin = v
		}
	}

	// 3. 严格定义信号：
	//    side == "up"   → 寻找顶部反转 → 曲率向下弯（diff < 0）且 lastMax < prevMax
	//    side == "down" → 寻找底部反转 → 曲率向上弯（diff > 0）且 lastMin > prevMin

	if side == "up" {
		if !(diff < 0 && lastMax < prevMax) {
			return false
		}
	} else if side == "down" {
		if !(diff > 0 && lastMin > prevMin) {
			return false
		}
	}

	// 4. 极端曲率过滤：diff 落在底部 10% 或顶部 10%
	curvs := make([]float64, cf.windowSize)
	for i := 10; i < cf.windowSize; i++ {
		p := ws[i-10 : i-5]
		l := ws[i-5 : i]
		curvs[i] = mean(l) - mean(p)
	}

	sort.Float64s(curvs)

	// 10% 分位数
	idx := int(float64(len(curvs)) * 0.1)
	if idx < 0 {
		idx = 0
	}
	low := curvs[idx]
	high := curvs[len(curvs)-idx-1]

	// up 要求 diff <= 低分位（极端向下弯）
	// down 要求 diff >= 高分位（极端向上弯）
	if side == "up" {
		return diff <= low
	} else {
		return diff >= high
	}
}
