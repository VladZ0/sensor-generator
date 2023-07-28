package sensor

import (
	"context"
	"net/http"
	"sensors-generator/internal/apperror"
	"sensors-generator/pkg/logging"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	sensorPath         = "api/v1/sensor/:codeName"
	regionPath         = "api/v1/region"
	temperatureMinPath = "/temperature/min"
	temperatureMaxPath = "/temperature/max"
	temperatureAvgPath = "/temperature/average"
)

type handler struct {
	sensorService ISensorService
	logger        *logging.Logger
}

func NewHandler(sensorService ISensorService, logger *logging.Logger) *handler {
	return &handler{
		sensorService: sensorService,
		logger:        logger,
	}
}

func (h *handler) Register(router *gin.Engine) {
	region := router.Group(regionPath)
	{
		region.GET(temperatureMaxPath, h.MaxTemperature)
		region.GET(temperatureMinPath, h.MinTemperature)
	}

	sensor := router.Group(sensorPath)
	{
		sensor.GET(temperatureAvgPath, h.AvgTemperature)
	}
}

// MinTemperature
// @Summary Min Temperature
// @Tags Sensors
// @Param xMin query string true "Minimum value for x coordinate"
// @Param yMin query string true "Minimum value for y coordinate"
// @Param zMin query string true "Minimum value for z coordinate"
// @Param xMax query string true "Maximum value for x coordinate"
// @Param yMax query string true "Maximum value for y coordinate"
// @Param zMax query string true "Maximum value for z coordinate"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/region/temperature/min [get]
func (h *handler) MinTemperature(c *gin.Context) {
	h.logger.Info("MIN TEMPERATURE.")

	xMin := c.Query("xMin")
	yMin := c.Query("yMin")
	zMin := c.Query("zMin")

	xMax := c.Query("xMax")
	yMax := c.Query("yMax")
	zMax := c.Query("zMax")

	minCoords, err := NewCoordsFromString(xMin, yMin, zMin)
	if err != nil {
		h.logger.Errorf("Cannot convert coords, due to error: %v", err)
		c.Error(apperror.ErrInternalSystem)
		return
	}

	maxCoords, err := NewCoordsFromString(xMax, yMax, zMax)
	if err != nil {
		h.logger.Errorf("Cannot convert coords, due to error: %v", err)
		c.Error(apperror.ErrInternalSystem)
		return
	}

	minTemperature, err := h.sensorService.GetExtremumTemperatureForRegion(context.Background(),
		minCoords, maxCoords, true)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"min_temperature": minTemperature})
}

// MaxTemperature
// @Summary Max Temperature
// @Tags Sensors
// @Param xMin query string true "Minimum value for x coordinate"
// @Param yMin query string true "Minimum value for y coordinate"
// @Param zMin query string true "Minimum value for z coordinate"
// @Param xMax query string true "Maximum value for x coordinate"
// @Param yMax query string true "Maximum value for y coordinate"
// @Param zMax query string true "Maximum value for z coordinate"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/region/temperature/max [get]
func (h *handler) MaxTemperature(c *gin.Context) {
	h.logger.Info("MAX TEMPERATURE.")

	xMin := c.Query("xMin")
	yMin := c.Query("yMin")
	zMin := c.Query("zMin")

	xMax := c.Query("xMax")
	yMax := c.Query("yMax")
	zMax := c.Query("zMax")

	minCoords, err := NewCoordsFromString(xMin, yMin, zMin)
	if err != nil {
		h.logger.Errorf("Cannot convert coords, due to error: %v", err)
		c.Error(apperror.ErrBadRequest)
		return
	}

	maxCoords, err := NewCoordsFromString(xMax, yMax, zMax)
	if err != nil {
		h.logger.Errorf("Cannot convert coords, due to error: %v", err)
		c.Error(apperror.ErrBadRequest)
		return
	}

	maxTemperature, err := h.sensorService.GetExtremumTemperatureForRegion(context.Background(),
		minCoords, maxCoords, false)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_temperature": maxTemperature})
}

// AvgTemperature
// @Summary Average Temperature
// @Tags Sensors
// @Param codeName path string true "Name of the group"
// @Param from query int false "from"
// @Param till query int false "till"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/sensor/{codeName}/temperature/average [get]
func (h *handler) AvgTemperature(c *gin.Context) {
	h.logger.Info("AVERAGE TEMPERATURE.")

	codeNameQ := c.Param("codeName")
	codeName, err := NewCodenameFromString(codeNameQ)
	if err != nil {
		c.Error(err)
		return
	}

	filters := SensorFilters{
		CodeName: codeName,
	}

	if from, ok := c.GetQuery("from"); ok {
		fromTS, err := strconv.Atoi(from)
		if err != nil {
			h.logger.Errorf("Cannot parse string, due to error: %v", err)
			c.Error(apperror.ErrBadRequest)
			return
		}
		filters.FromDate = time.Unix(int64(fromTS), 0)
	}

	if till, ok := c.GetQuery("till"); ok {
		tillTS, err := strconv.Atoi(till)
		if err != nil {
			h.logger.Errorf("Cannot parse string, due to error: %v", err)
			c.Error(apperror.ErrBadRequest)
			return
		}
		filters.TillDate = time.Unix(int64(tillTS), 0)
	}

	avgTemperature, err := h.sensorService.GetAvgTemperatureForSensor(context.Background(), filters)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"avg_temperature": avgTemperature})
}
