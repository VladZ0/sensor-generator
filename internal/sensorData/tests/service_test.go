package sensordata

import (
	"context"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SensorDataService_Create(t *testing.T) {
	sensorDataList := []sensordata.CreateSensorDataDTO{
		{SensorID: 1, Temperature: 25.5, Transparency: 90},
		{SensorID: 2, Temperature: 28.0, Transparency: 80},
	}
	logging.Init("trace", true)

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockSensorDataRepository{}

		mockRepo.On("Create", mock.Anything, mock.Anything).
			Return(42, nil).
			Times(len(sensorDataList))

		service := sensordata.NewService(mockRepo, logging.GetLogger(), nil)
		ids, err := service.Create(context.Background(), sensorDataList...)

		assert.NoError(t, err)
		assert.Len(t, ids, len(sensorDataList))
		for _, id := range ids {
			assert.NotEqual(t, 0, id)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestAddDetectedSpieces(t *testing.T) {
	repo := &MockSensorDataRepository{}

	service := sensordata.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	sensorDataID := 123
	species1 := spiece.Spiece{Name: "Species1"}
	species2 := spiece.Spiece{Name: "Species2"}

	repo.On("AddDetectedSpiece", ctx, sensorDataID, species1).Return(nil)
	repo.On("AddDetectedSpiece", ctx, sensorDataID, species2).Return(nil)

	err := service.AddDetectedSpieces(ctx, sensorDataID, species1, species2)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestGetOneByID(t *testing.T) {
	repo := &MockSensorDataRepository{}

	service := sensordata.NewService(repo, logging.GetLogger(), nil)

	ctx := context.Background()
	id := 123
	filters := sensordata.SensorDataFilters{}
	expectedSensorData := &sensordata.SensorData{
		ID:           id,
		SensorID:     456,
		Temperature:  25.5,
		Transparency: 80,
		DetectedSpieces: []spiece.Spiece{
			{ID: 1, Name: "Species1"},
			{ID: 2, Name: "Species2"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("FindOneByID", ctx, id, filters).Return(expectedSensorData, nil)

	sensorData, err := service.GetOneByID(ctx, id, filters)

	assert.NoError(t, err)
	assert.Equal(t, expectedSensorData, sensorData)

	repo.AssertExpectations(t)
}
