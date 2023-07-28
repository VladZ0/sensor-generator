package group

import (
	"context"
	"sensors-generator/internal/group"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GroupService_GetSpiecesInGroup(t *testing.T) {
	mockRepo := &MockGroupRepository{}

	logging.Init("trace", true)
	service := group.NewService(mockRepo, &MockCache{}, logging.GetLogger(), nil)

	groupName := "alpha"

	expectedSpieces := map[*spiece.Spiece]int{
		{ID: 1, Name: "Spiece1"}: 10,
		{ID: 2, Name: "Spiece2"}: 5,
	}

	mockRepo.On("FindSpiecesInGroup", mock.Anything, groupName, group.SensorGroupFilters{}).
		Return(expectedSpieces, nil)

	sp, err := service.GetSpiecesInGroup(context.Background(), groupName, group.SensorGroupFilters{})
	assert.NoError(t, err)
	assert.NotNil(t, sp)
	assert.Equal(t, len(expectedSpieces), len(sp))

	for expectedSpiece, expectedCount := range expectedSpieces {
		actualCount, found := sp[expectedSpiece]
		assert.True(t, found)
		assert.Equal(t, expectedCount, actualCount)
	}

	mockRepo.AssertExpectations(t)
}

func Test_GroupService_GetAvgTransparencyInGroup(t *testing.T) {
	mockCache := &MockCache{}
	mockRepo := &MockGroupRepository{}

	service := group.NewService(mockRepo, mockCache, logging.GetLogger(), nil)

	groupName := "alpha"
	expectedTransparency := uint8(80)

	mockCache.On("Get", mock.Anything, groupName+"AvgTransparency").
		Return("", nil)

	mockRepo.On("FindAvgTransparencyInGroup", mock.Anything, groupName, group.SensorGroupFilters{}).
		Return(expectedTransparency, nil)

	mockCache.On("Set", mock.Anything, groupName+"AvgTransparency", expectedTransparency, mock.Anything).
		Return(nil)

	transparency, err := service.GetAvgTrasparencyInGroup(context.Background(), groupName, group.SensorGroupFilters{})
	assert.NoError(t, err)
	assert.Equal(t, expectedTransparency, transparency)

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_GroupService_GetAvgTemperatureInGroup(t *testing.T) {
	mockCache := &MockCache{}
	mockRepo := &MockGroupRepository{}

	service := group.NewService(mockRepo, mockCache, logging.GetLogger(), nil)

	groupName := "alpha"
	expectedTemperature := float32(25.5)

	mockCache.On("Get", mock.Anything, groupName+"AvgTemperature").
		Return("", nil)

	mockRepo.On("FindAvgTemperatureInGroup", mock.Anything, groupName, group.SensorGroupFilters{}).
		Return(expectedTemperature, nil)

	mockCache.On("Set", mock.Anything, groupName+"AvgTemperature", expectedTemperature, mock.Anything).
		Return(nil)

	temperature, err := service.GetAvgTemperatureInGroup(context.Background(), groupName, group.SensorGroupFilters{})
	assert.NoError(t, err)
	assert.Equal(t, expectedTemperature, temperature)

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_GroupService_Create(t *testing.T) {
	mockRepo := &MockGroupRepository{}

	service := group.NewService(mockRepo, nil, logging.GetLogger(), nil)

	ctx := context.Background()
	group1 := group.CreateSensorGroupDTO{Name: "alpha"}
	group2 := group.CreateSensorGroupDTO{Name: "beta"}

	mockRepo.On("Create", ctx, group1).Return(nil)
	mockRepo.On("Create", ctx, group2).Return(nil)

	err := service.Create(ctx, group1, group2)
	assert.NoError(t, err)

	mockRepo.AssertNumberOfCalls(t, "Create", 2)
	mockRepo.AssertCalled(t, "Create", ctx, group1)
	mockRepo.AssertCalled(t, "Create", ctx, group2)
}
