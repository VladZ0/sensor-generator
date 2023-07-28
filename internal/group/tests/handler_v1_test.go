package group

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sensors-generator/internal/group"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func sortBySpieceName(species []map[string]interface{}) {
	sort.Slice(species, func(i, j int) bool {
		name1 := species[i]["name"].(string)
		name2 := species[j]["name"].(string)
		return name1 < name2
	})
}

func Test_Handler_GetSpiecesInGroup(t *testing.T) {
	logging.Init("trace", true)
	mockService := &MockSensorGroupService{}
	handler := group.NewHandler(mockService, logging.GetLogger())

	expectedSpieces := map[spiece.Spiece]int{
		{ID: 1, Name: "Spiece1"}: 5,
		{ID: 2, Name: "Spiece2"}: 10,
	}
	mockService.On("GetSpiecesInGroup", mock.Anything, "alpha", group.SensorGroupFilters{}).Return(expectedSpieces, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "groupName", Value: "alpha"})

	handler.GetSpiecesInGroup(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	for _, item := range response["spieces"] {
		item["count"] = int(item["count"].(float64))
	}

	sortBySpieceName(response["spieces"])

	expectedResponse := map[string][]map[string]interface{}{
		"spieces": {
			{"name": "Spiece1", "count": 5},
			{"name": "Spiece2", "count": 10},
		},
	}
	assert.Equal(t, expectedResponse, response)
}

func Test_Handler_GetTopNSpiecesInGroup(t *testing.T) {
	logging.Init("trace", true)
	mockService := &MockSensorGroupService{}
	handler := group.NewHandler(mockService, logging.GetLogger())

	expectedSpieces := map[spiece.Spiece]int{
		{ID: 1, Name: "Spiece1"}: 10,
		{ID: 2, Name: "Spiece2"}: 5,
	}
	mockService.On("GetSpiecesInGroup", mock.AnythingOfType("*context.emptyCtx"), "alpha", mock.Anything).Return(expectedSpieces, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "groupName", Value: "alpha"})
	c.Params = append(c.Params, gin.Param{Key: "N", Value: "2"})

	handler.GetTopNSpiecesInGroup(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	for _, spiece := range response["spieces"] {
		spiece["count"] = int(spiece["count"].(float64))
	}

	sortBySpieceName(response["spieces"])

	expectedResponse := map[string][]map[string]interface{}{
		"spieces": {
			{"name": "Spiece1", "count": 10},
			{"name": "Spiece2", "count": 5},
		},
	}
	assert.Equal(t, expectedResponse, response)
}

func Test_Handler_GetAvgTransparencyInGroup(t *testing.T) {
	mockService := &MockSensorGroupService{}
	handler := group.NewHandler(mockService, logging.GetLogger())

	expectedAvgTransparency := uint8(80)
	mockService.On("GetAvgTrasparencyInGroup", mock.Anything, "alpha", group.SensorGroupFilters{}).Return(expectedAvgTransparency, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "groupName", Value: "alpha"})

	handler.GetAvgTransparencyInGroup(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]uint8
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	expectedResponse := map[string]uint8{
		"average_transparency": expectedAvgTransparency,
	}
	assert.Equal(t, expectedResponse, response)
}

func Test_Handler_GetAvgTemperatureInGroup(t *testing.T) {
	mockService := &MockSensorGroupService{}
	handler := group.NewHandler(mockService, logging.GetLogger())

	expectedAvgTemperature := float32(25.5)
	mockService.On("GetAvgTemperatureInGroup", mock.Anything, "alpha", group.SensorGroupFilters{}).Return(expectedAvgTemperature, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "groupName", Value: "alpha"})

	handler.GetAvgTemperatureInGroup(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]float32
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	expectedResponse := map[string]float32{
		"average_temperature": expectedAvgTemperature,
	}
	assert.Equal(t, expectedResponse, response)
}
