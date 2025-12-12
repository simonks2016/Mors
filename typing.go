package main

type PositionManagerGet func(string) (float64, error)
type HandlerByStrategy func(name string) error

type BasicStrategy interface {
	OnDataUpdate(parameters ...TickDataInterface) (*Signal, error)
}

type TickDataInterface interface {
	Name() string
	Value() float64
}

type TickData struct {
	name  string
	value float64
}

func (t *TickData) Name() string {
	return t.name
}
func (t *TickData) Value() float64 {
	return t.value
}
func NewTickData(name string, value float64) *TickData {
	return &TickData{
		name: name,
	}
}
