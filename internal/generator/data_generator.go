package generator

import (
	"context"
	"math/rand"
	"sensors-generator/internal/apperror"
	"sensors-generator/internal/sensor"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"sync"
	"time"
)

type DataGenerator struct {
	services  Services
	randomGen IRandomGenerator
	sync.Mutex
}

func NewDataGenerator(services Services, randomGen IRandomGenerator) *DataGenerator {
	return &DataGenerator{
		services:  services,
		randomGen: randomGen,
	}
}

func (dg *DataGenerator) Generate() error {
	sensors, err := dg.services.SensorService.GetAll(context.Background(), sensor.SensorFilters{})
	if err != nil {
		return apperror.ErrInternalSystem
	}

	for _, sensor := range sensors {
		go dg.generateData(sensor)
	}

	return nil
}

func (dg *DataGenerator) generateData(sensor sensor.Sensor) {
	for {
		dg.Lock()
		sdata := sensordata.CreateSensorDataDTO{
			SensorID:     sensor.ID,
			Temperature:  dg.randomGen.GenerateTemperatureBasedOnZ(sensor.Coords.Z),
			Transparency: dg.randomGen.GenerateTransparency(),
		}
		dg.Unlock()

		sensorDataIDS, err := dg.services.SensorDataService.Create(context.Background(), sdata)
		if err != nil {
			logging.GetLogger().Errorf("Sensor data generetor error: %v", err)
		}
		if err := dg.generateDetectedSpieces(sensorDataIDS[0], dg.services); err != nil {
			logging.GetLogger().Errorf("Sensor data spieces generetor error: %v", err)
		}
		time.Sleep(sensor.DataOutputRate * time.Second)
	}
}

func (dg *DataGenerator) generateDetectedSpieces(sensorDataID int, services Services) error {
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
