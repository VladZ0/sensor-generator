package utils

import (
	"context"
	"math/rand"
	"sensors-generator/internal/apperror"
	"sensors-generator/internal/group"
	"sensors-generator/internal/sensor"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"time"
)

type Services struct {
	SensorService      sensor.ISensorService
	SensorGroupService group.ISensorGroupService
	SpieceService      spiece.ISpiecesService
	SensorDataService  sensordata.ISensorDataService
}

func GenerateData(services Services) error {
	sensors, err := services.SensorService.GetAll(context.Background(), sensor.SensorFilters{})
	if err != nil {
		return apperror.ErrInternalSystem
	}

	for _, sensor := range sensors {
		go generateData(sensor, services)
	}

	return nil
}

func generateData(sensor sensor.Sensor, services Services) {
	for {
		sdata := sensordata.CreateSensorDataDTO{
			SensorID:     sensor.ID,
			Temperature:  GenerateTemperature(sensor.Coords.Z),
			Transparency: GenerateTransparency(),
		}

		sensorDataIDS, err := services.SensorDataService.Create(context.Background(), sdata)
		if err != nil {
			logging.GetLogger().Errorf("Sensor data generetor error: %v", err)
		}
		if err := generateDetectedSpieces(sensorDataIDS[0], services); err != nil {
			logging.GetLogger().Errorf("Sensor data spieces generetor error: %v", err)
		}
		time.Sleep(sensor.DataOutputRate * time.Second)
	}
}

func generateDetectedSpieces(sensorDataID int, services Services) error {
	spieces, err := services.SpieceService.GetAll(context.Background(), spiece.SpieceFilters{})
	if err != nil {
		return err
	}

	detectedSpieces := make([]spiece.Spiece, 0)

	count := rand.Intn(50)

	for i := 0; i < count; i++ {
		index := rand.Intn(len(spieces))
		detectedSpieces = append(detectedSpieces, spieces[index])
	}

	services.SensorDataService.AddDetectedSpieces(context.Background(), sensorDataID, detectedSpieces...)
	return nil
}

func GenerateGroupsSensorsSpieces(groups []group.CreateSensorGroupDTO, sensors []sensor.CreateSensorDTO,
	spieces []spiece.CreateSpieceDTO, services Services) error {

	if err := services.SensorGroupService.Create(context.Background(), groups...); err != nil {
		return err
	}

	if err := services.SensorService.Create(context.Background(), sensors...); err != nil {
		return err
	}

	if err := services.SpieceService.Create(context.Background(), spieces...); err != nil {
		return err
	}

	return nil
}
