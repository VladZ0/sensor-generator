package generator

import (
	"sensors-generator/internal/group"
	"sensors-generator/internal/sensor"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
)

type Services struct {
	SensorService      sensor.ISensorService
	SensorGroupService group.ISensorGroupService
	SpieceService      spiece.ISpiecesService
	SensorDataService  sensordata.ISensorDataService
}
