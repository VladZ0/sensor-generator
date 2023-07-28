package group

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
	basicPath           = "api/v1/group/:groupName"
	spiecesPath         = "/spieces"
	topNSpiecesPath     = spiecesPath + "/top/:N"
	temperatureAvgPath  = "/temperature/average"
	transparencyAvgPath = "/transparency/average"
)

type handler struct {
	sensorGroupService ISensorGroupService
	logger             *logging.Logger
}

func NewHandler(sensorGroupService ISensorGroupService, logger *logging.Logger) *handler {
	return &handler{
		sensorGroupService: sensorGroupService,
		logger:             logger,
	}
}

func (h *handler) Register(router *gin.Engine) {
	group := router.Group(basicPath)
	{
		group.GET(spiecesPath, h.GetSpiecesInGroup)
		group.GET(topNSpiecesPath, h.GetTopNSpiecesInGroup)
		group.GET(temperatureAvgPath, h.GetAvgTemperatureInGroup)
		group.GET(transparencyAvgPath, h.GetAvgTransparencyInGroup)
	}
}

// GetSpiecesInGroupHandler
// @Summary Spieces in group
// @Tags Groups
// @Param groupName path string true "Name of the group"
// @Param from query int false "from"
// @Param till query int false "till"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/group/{groupName}/spieces [get]
func (h *handler) GetSpiecesInGroup(c *gin.Context) {
	groupName := c.Param("groupName")

	filters := SensorGroupFilters{}

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

	spieces, err := h.sensorGroupService.GetSpiecesInGroup(context.Background(), groupName, filters)
	if err != nil {
		c.Error(err)
		return
	}

	spiecesJSON := make([]gin.H, 0)
	for spiece, count := range spieces {
		spieceJSON := gin.H{"name": spiece.Name, "count": count}
		spiecesJSON = append(spiecesJSON, spieceJSON)
	}

	c.JSON(http.StatusOK, gin.H{"spieces": spiecesJSON})
}

// GetTopNSpiecesInGroupHandler
// @Summary Top N spieces in group
// @Tags Groups
// @Param groupName path string true "Name of the group"
// @Param N path string true "Top N"
// @Param from query int false "from"
// @Param till query int false "till"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/group/{groupName}/spieces/top/{N} [get]
func (h *handler) GetTopNSpiecesInGroup(c *gin.Context) {
	var N int
	var err error

	groupName := c.Param("groupName")

	filters := SensorGroupFilters{}

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

	if N, err = strconv.Atoi(c.Param("N")); err != nil {
		c.Error(apperror.ErrorWithMessage(apperror.ErrBadRequest, "N is bad."))
	}

	if N < 0 {
		c.Error(apperror.ErrorWithMessage(apperror.ErrBadRequest, "N should be >= 0."))
	}

	filters.TopLimit = N

	spieces, err := h.sensorGroupService.GetSpiecesInGroup(context.Background(),
		groupName, filters)
	if err != nil {
		c.Error(err)
		return
	}

	spiecesJSON := make([]gin.H, 0)
	for spiece, count := range spieces {
		spieceJSON := gin.H{"name": spiece.Name, "count": count}
		spiecesJSON = append(spiecesJSON, spieceJSON)
	}

	c.JSON(http.StatusOK, gin.H{"spieces": spiecesJSON})
}

// GetAvgTransparencyInGroupHandler
// @Summary Average transparency in group
// @Tags Groups
// @Param groupName path string true "Name of the group"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/group/{groupName}/transparency/average [get]
func (h *handler) GetAvgTransparencyInGroup(c *gin.Context) {
	groupName := c.Param("groupName")

	avgTransparency, err := h.sensorGroupService.GetAvgTrasparencyInGroup(context.Background(), groupName, SensorGroupFilters{})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"average_transparency": avgTransparency})
}

// GetAvgTemperatureInGroupHandler
// @Summary Average temperature in group
// @Tags Groups
// @Param groupName path string true "Name of the group"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/group/{groupName}/temperature/average [get]
func (h *handler) GetAvgTemperatureInGroup(c *gin.Context) {
	groupName := c.Param("groupName")

	avgTemperature, err := h.sensorGroupService.GetAvgTemperatureInGroup(context.Background(), groupName, SensorGroupFilters{})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"average_temperature": avgTemperature})
}
