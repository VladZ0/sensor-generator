package utils

import (
	"math"
	"math/rand"
	"sensors-generator/internal/sensor"
)

func PointsDistance(point1 sensor.Coordinates, point2 sensor.Coordinates) float64 {
	distance := math.Sqrt(math.Pow(point2.X-point1.X, 2) +
		math.Pow(point2.Y-point1.Y, 2) + math.Pow(point2.Z-point1.Z, 2))

	return distance
}

func GenerateTemperature(z float64) float32 {
	depthMultiplier := 1.0 + float64(z)/100.0

	temperature := float32(rand.Intn(30) + 1)
	temperature *= rand.Float32()
	temperature = float32(math.Round(float64(temperature)*100) / 100)

	temperature *= float32(depthMultiplier)

	return temperature
}

func GenerateTransparency() uint8 {
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

func minIndex(arr []float64) int {
	minIndex := 0

	if len(arr) <= 0 {
		return -1
	}

	for i := 0; i < len(arr); i++ {
		if arr[minIndex] < arr[i] {
			minIndex = i
		}
	}

	return minIndex
}
