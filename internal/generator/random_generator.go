package generator

import (
	"math"
	"math/rand"
)

type RandomGenerator struct {
}

func NewRandomGenerator() *RandomGenerator {
	return &RandomGenerator{}
}

func (rg *RandomGenerator) GenerateTemperatureBasedOnZ(z float64) float32 {
	depthMultiplier := 1.0 + float64(z)/100.0

	temperature := float32(rand.Intn(30) + 1)
	temperature *= rand.Float32()
	temperature = float32(math.Round(float64(temperature)*100) / 100)

	temperature *= float32(depthMultiplier)

	return temperature
}

func (rg *RandomGenerator) GenerateTemperature() float32 {
	temperature := float32(rand.Intn(30) + 1)

	return temperature
}

func (rg *RandomGenerator) GenerateTransparency() uint8 {
	baseTransparency := rand.Intn(101)
	offset := rand.Intn(11) - 5
	transparency := baseTransparency + offset

	if transparency < 0 {
		transparency = 0
	} else if transparency > 100 {
		transparency = 100
	}

	return uint8(transparency)
}
