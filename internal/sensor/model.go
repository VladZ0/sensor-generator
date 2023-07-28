package sensor

import (
	"regexp"
	"sensors-generator/internal/apperror"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"strconv"
	"strings"
	"time"
)

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Codename struct {
	GroupName string `json:"group_name"`
	Index     int    `json:"index"`
}

type Sensor struct {
	ID             int             `json:"-"`
	CodeName       Codename        `json:"codename"`
	Coords         Coordinates     `json:"coordinates"`
	DataOutputRate time.Duration   `json:"data_output_rate"`
	Spieces        []spiece.Spiece `json:"spieces"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CreateSensorDTO struct {
	CodeName       Codename      `json:"codename"`
	Coords         Coordinates   `json:"coordinates"`
	DataOutputRate time.Duration `json:"data_output_rate"`
}

type SensorFilters struct {
	CodeName Codename
	FromDate time.Time
	TillDate time.Time
}

func NewCoordsFromString(x, y, z string) (Coordinates, error) {
	var X, Y, Z float64
	X, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return Coordinates{}, err
	}

	Y, err = strconv.ParseFloat(y, 64)
	if err != nil {
		return Coordinates{}, err
	}

	Z, err = strconv.ParseFloat(z, 64)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		X: X,
		Y: Y,
		Z: Z,
	}, nil
}

func NewCodenameFromString(codename string) (Codename, error) {
	splitted := strings.Split(codename, " ")

	pattern := `^[a-zA-Z]+\s\d+$`

	matched, err := regexp.MatchString(pattern, codename)
	if err != nil {
		logging.GetLogger().Errorf("Error while regex matching: %v", err)
		return Codename{}, apperror.ErrInternalSystem
	}

	if !matched {
		logging.GetLogger().Errorf("Regex not matched: %v", err)
		return Codename{}, apperror.ErrorWithMessage(apperror.ErrBadRequest, "Wrong codename.")
	}

	gName := splitted[0]
	index, err := strconv.Atoi(splitted[1])
	if err != nil {
		logging.GetLogger().Errorf("Cannot parse index, due to error: %v", err)
		return Codename{}, apperror.ErrInternalSystem
	}

	cName := Codename{
		GroupName: gName,
		Index:     index,
	}

	return cName, nil
}

func (cdn *Codename) IsEmpty() bool {
	if cdn.GroupName == "" && cdn.Index == 0 {
		return true
	}

	return false
}
