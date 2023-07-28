package sensor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sensors-generator/internal/sensor"
	"sensors-generator/pkg/logging"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Handler_MinTemperature(t *testing.T) {
	logging.Init("trace", true)
	mockSensorService := &MockSensorService{}
	handler := sensor.NewHandler(mockSensorService, logging.GetLogger())

	mockMinTemperature := float32(10.5)
	mockMinCoords := sensor.Coordinates{X: 1.0, Y: 2.0, Z: 3.0}
	mockMaxCoords := sensor.Coordinates{X: 10.0, Y: 20.0, Z: 30.0}

	mockSensorService.On("GetExtremumTemperatureForRegion", mock.Anything, mockMinCoords, mockMaxCoords, true).
		Return(mockMinTemperature, nil)

	req, err := http.NewRequest("GET", "/min_temperature?xMin=1&yMin=2&zMin=3&xMax=10&yMax=20&zMax=30", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.MinTemperature(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]float32
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	expectedResponse := map[string]float32{"min_temperature": mockMinTemperature}
	assert.Equal(t, expectedResponse, response)
}

func Test_Handler_MaxTemperature(t *testing.T) {
	mockService := &MockSensorService{}
	handler := sensor.NewHandler(mockService, logging.GetLogger())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/max_temperature?xMin=1&yMin=2&zMin=3&xMax=4&yMax=5&zMax=6", nil)

	expectedMinCoords := sensor.Coordinates{X: 1, Y: 2, Z: 3}
	expectedMaxCoords := sensor.Coordinates{X: 4, Y: 5, Z: 6}
	expectedTemperature := float32(25.5)
	mockService.On("GetExtremumTemperatureForRegion", mock.Anything, expectedMinCoords, expectedMaxCoords, false).
		Return(expectedTemperature, nil)

	handler.MaxTemperature(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]float32
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	assert.Equal(t, expectedTemperature, response["max_temperature"])
}
