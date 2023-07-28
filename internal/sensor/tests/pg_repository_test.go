package sensor

import (
	"context"
	"sensors-generator/internal/sensor"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_SensorRepository_Create(t *testing.T) {
	mockDTO := sensor.CreateSensorDTO{
		CodeName: sensor.Codename{
			GroupName: "alpha",
			Index:     1,
		},
		Coords: sensor.Coordinates{
			X: 34.33,
			Y: 27.24,
			Z: 4,
		},
		DataOutputRate: 15,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logging.Init("trace", true)

	repo := sensor.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectQuery("SELECT id FROM sensor_groups WHERE name = \\$1").
		WithArgs(mockDTO.CodeName.GroupName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO sensors`).WithArgs(
		1, mockDTO.CodeName.Index, mockDTO.Coords.X, mockDTO.Coords.Y,
		mockDTO.Coords.Z, mockDTO.DataOutputRate, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Create(context.Background(), mockDTO); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorRepository_AddSensorToGroup(t *testing.T) {
	sensorID, groupID := 23, 444

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensor.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectExec("UPDATE sensors SET group_id=\\$1 WHERE id=\\$2").
		WithArgs(groupID, sensorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.AddSensorToGroup(context.Background(), sensorID, groupID); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorRepository_FindMaxTemperatureForRegion(t *testing.T) {
	minCoords := sensor.Coordinates{X: 10.0, Y: 20.0, Z: 5.0}
	maxCoords := sensor.Coordinates{X: 20.0, Y: 30.0, Z: 15.0}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensor.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	expectedTemperature := float32(25.0)
	mockRows := sqlmock.NewRows([]string{"max_temperature"}).AddRow(expectedTemperature)
	mock.ExpectQuery(`SELECT MAX\(sd\.temperature\) FROM sensors as sens(.+)`).WithArgs(
		maxCoords.X, minCoords.X, maxCoords.Y, minCoords.Y, maxCoords.Z, minCoords.Z,
	).WillReturnRows(mockRows)

	temperature, err := repo.FindMaxTemperatureForRegion(context.Background(), minCoords, maxCoords)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if temperature != expectedTemperature {
		t.Errorf("expected temperature %f, got %f", expectedTemperature, temperature)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorRepository_FindMinTemperatureForRegion(t *testing.T) {
	minCoords := sensor.Coordinates{
		X: 10.0,
		Y: 20.0,
		Z: 5.0,
	}

	maxCoords := sensor.Coordinates{
		X: 20.0,
		Y: 30.0,
		Z: 10.0,
	}

	expectedTemperature := float32(15.5)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensor.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectQuery(`SELECT MIN\(sd\.temperature\) FROM sensors as sens`).
		WithArgs(maxCoords.X, minCoords.X, maxCoords.Y, minCoords.Y, maxCoords.Z, minCoords.Z).
		WillReturnRows(sqlmock.NewRows([]string{"min"}).AddRow(expectedTemperature))

	temperature, err := repo.FindMinTemperatureForRegion(context.Background(), minCoords, maxCoords)
	if err != nil {
		t.Errorf("error was not expected while finding min temperature: %s", err)
	}

	if temperature != expectedTemperature {
		t.Errorf("expected temperature %f, but got %f", expectedTemperature, temperature)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorRepository_FindAvgTemperatureForSensor(t *testing.T) {
	mockFilters := sensor.SensorFilters{
		CodeName: sensor.Codename{
			GroupName: "alpha",
			Index:     1,
		},
		FromDate: time.Date(2023, time.July, 1, 0, 0, 0, 0, time.UTC),
		TillDate: time.Date(2023, time.July, 31, 23, 59, 59, 0, time.UTC),
	}

	expectedTemperature := float32(25.5)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensor.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectQuery(`SELECT AVG\(sd\.temperature\) FROM sensors AS sens`).
		WithArgs(mockFilters.CodeName.GroupName, mockFilters.CodeName.Index,
			mockFilters.FromDate, mockFilters.TillDate).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(expectedTemperature))

	temperature, err := repo.FindAvgTemperatureForSensor(context.Background(), mockFilters)
	if err != nil {
		t.Errorf("error was not expected while finding average temperature: %s", err)
	}

	if temperature != expectedTemperature {
		t.Errorf("expected temperature %f, but got %f", expectedTemperature, temperature)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
