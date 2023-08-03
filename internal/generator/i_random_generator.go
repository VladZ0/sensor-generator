package generator

type IRandomGenerator interface {
	GenerateTemperatureBasedOnZ(z float64) float32
	GenerateTemperature() float32
	GenerateTransparency() uint8
}
