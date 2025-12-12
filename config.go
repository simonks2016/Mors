package main

type MORSConfig struct {
	maxLoss     float64
	windowsSize int
	emaAlpha    float64

	StrongUpAvg   float64
	UpAvg         float64
	StrongDownAvg float64
	DownAvg       float64

	PosRatioStrong   float64
	PosRatioUp       float64
	NegRatioStrong   float64
	NegRatioDown     float64
	MinWindowSamples int
}

func NewConfig(windowSize int, emaAlpha float64) *MORSConfig {

	return &MORSConfig{
		windowsSize:      windowSize,
		emaAlpha:         emaAlpha,
		StrongUpAvg:      0.35,
		UpAvg:            0.15,
		StrongDownAvg:    -0.35,
		DownAvg:          -0.15,
		PosRatioStrong:   0.7,
		PosRatioUp:       0.6,
		NegRatioStrong:   0.7,
		NegRatioDown:     0.6,
		MinWindowSamples: 100,
	}
}

func (c *MORSConfig) SetScale(strongUp, up, StrongDown, down float64) *MORSConfig {

	c.StrongUpAvg = strongUp
	c.UpAvg = up
	c.StrongDownAvg = StrongDown
	c.DownAvg = down
	return c
}
func (c *MORSConfig) SetPosThreshold(strong, up float64) *MORSConfig {
	c.PosRatioStrong = strong
	c.PosRatioUp = up
	return c
}
func (c *MORSConfig) SetNegThreshold(strong, down float64) *MORSConfig {
	c.NegRatioStrong = strong
	c.NegRatioDown = down
	return c
}
func (c *MORSConfig) SetMinWindowSamples(minWindowSamples int) *MORSConfig {
	c.MinWindowSamples = minWindowSamples
	return c
}
