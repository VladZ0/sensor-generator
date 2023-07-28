package sensor

import (
	"context"
	"sensors-generator/internal/sensor"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SensorService_GetAll(t *testing.T) {
	logging.Init("trace", true)
	repo := &MockSensorRepository{}

	service := sensor.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	filters := sensor.SensorFilters{
		FromDate: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		TillDate: time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
	}

	expectedSensors := []sensor.Sensor{
		{
			ID: 1,
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 10.0,
				Y: 20.0,
				Z: 5.0,
			},
			DataOutputRate: time.Second * 10,
			Spieces: []spiece.Spiece{
				{ID: 1, Name: "Species1"},
				{ID: 2, Name: "Species2"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID: 2,
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 15.0,
				Y: 25.0,
				Z: 8.0,
			},
			DataOutputRate: time.Second * 5,
			Spieces: []spiece.Spiece{
				{ID: 3, Name: "Species3"},
				{ID: 4, Name: "Species4"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	repo.On("FindAll", ctx, filters).Return(expectedSensors, nil)

	sensors, err := service.GetAll(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, expectedSensors, sensors)

	repo.AssertExpectations(t)
}

func Test_SensorService_Create(t *testing.T) {
	repo := &MockSensorRepository{}

	service := sensor.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	sensors := []sensor.CreateSensorDTO{
		{
			CodeName: sensor.Codename{
				GroupName: "alpha",
				Index:     1,
			},
			Coords: sensor.Coordinates{
				X: 10.0,
				Y: 20.0,
				Z: 5.0,
			},
			DataOutputRate: time.Second * 10,
		},
		{
			CodeName: sensor.Codename{
				GroupName: "beta",
				Index:     2,
			},
			Coords: sensor.Coordinates{
				X: 15.0,
				Y: 25.0,
				Z: 8.0,
			},
			DataOutputRate: time.Second * 5,
		},
	}

	repo.On("Create", ctx, sensors[0]).Return(nil)
	repo.On("Create", ctx, sensors[1]).Return(nil)

	err := service.Create(ctx, sensors...)

	assert.NoError(t, err)

	repo.AssertCalled(t, "Create", ctx, sensors[0])
	repo.AssertCalled(t, "Create", ctx, sensors[1])
}

func Test_SensorService_AddSensorToGroup(t *testing.T) {
	repo := &MockSensorRepository{}

	service := sensor.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	sensorID := 123
	groupID := 456

	repo.On("AddSensorToGroup", ctx, sensorID, groupID).Return(nil)

	err := service.AddSensorToGroup(ctx, sensorID, groupID)

	assert.NoError(t, err)

	repo.AssertCalled(t, "AddSensorToGroup", ctx, sensorID, groupID)
}

func Test_SensorService_GetExtremumTemperatureForRegion(t *testing.T) {
	repo := &MockSensorRepository{}

	service := sensor.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	minCoords := sensor.Coordinates{X: 0, Y: 0, Z: 0}
	maxCoords := sensor.Coordinates{X: 10, Y: 10, Z: 10}
	minTemperature := float32(25.5)
	maxTemperature := float32(30.0)

	repo.On("FindMinTemperatureForRegion", ctx, minCoords, maxCoords).Return(minTemperature, nil)
	temp, err := service.GetExtremumTemperatureForRegion(ctx, minCoords, maxCoords, true)
	assert.NoError(t, err)
	assert.Equal(t, minTemperature, temp)

	repo.On("FindMaxTemperatureForRegion", ctx, minCoords, maxCoords).Return(maxTemperature, nil)
	temp, err = service.GetExtremumTemperatureForRegion(ctx, minCoords, maxCoords, false)
	assert.NoError(t, err)
	assert.Equal(t, maxTemperature, temp)

	repo.AssertExpectations(t)
}

func Test_SensorService_GetAvgTemperatureForSensor(t *testing.T) {
	repo := &MockSensorRepository{}

	service := sensor.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	filters := sensor.SensorFilters{
		CodeName: sensor.Codename{GroupName: "alpha", Index: 1},
		FromDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		TillDate: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	avgTemperature := float32(27.8)

	repo.On("FindAvgTemperatureForSensor", ctx, filters).Return(avgTemperature, nil)
	temp, err := service.GetAvgTemperatureForSensor(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, avgTemperature, temp)

	repo.AssertExpectations(t)
}
