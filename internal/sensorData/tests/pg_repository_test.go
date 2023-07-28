package sensordata

import (
	"context"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_SensorDataRepository_FindOneByID(t *testing.T) {
	mockSensorID := 1
	mockSensorDataID := 100

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logging.Init("trace", true)

	repo := sensordata.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT sd.id, sens.id, sd.temperature, sd.transparency, sd.created_at, sd.updated_at FROM sensor_data AS sd JOIN sensors sens ON sd.sensor_id=sens.id WHERE sd.id=?").
		WithArgs(mockSensorDataID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sensor_id", "temperature", "transparency", "created_at", "updated_at"}).
			AddRow(mockSensorDataID, mockSensorID, 25.5, 0.8, time.Now(), time.Now()))

	mock.ExpectQuery("SELECT s.id, s.name, s.created_at, s.updated_at FROM sensor_data AS sd JOIN detected_spieces ds ON sd.id=ds.sensor_data_id JOIN spieces s ON ds.spiece_id=s.id WHERE sd.id=?").
		WithArgs(mockSensorDataID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(1, "Species1", time.Now(), time.Now()).
			AddRow(2, "Species2", time.Now(), time.Now()))
	mock.ExpectCommit()

	sensorData, err := repo.FindOneByID(context.Background(), mockSensorDataID, sensordata.SensorDataFilters{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sensorData.ID != mockSensorDataID {
		t.Errorf("unexpected sensor data ID, got: %d, want: %d", sensorData.ID, mockSensorDataID)
	}
	if sensorData.SensorID != mockSensorID {
		t.Errorf("unexpected sensor ID, got: %d, want: %d", sensorData.SensorID, mockSensorID)
	}
	if len(sensorData.DetectedSpieces) != 2 {
		t.Errorf("unexpected number of detected species, got: %d, want: %d", len(sensorData.DetectedSpieces), 2)
	}

	expectedSpecies := []string{"Species1", "Species2"}
	for i, sp := range sensorData.DetectedSpieces {
		if sp.Name != expectedSpecies[i] {
			t.Errorf("unexpected detected species name, got: %s, want: %s", sp.Name, expectedSpecies[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorDataRepository_Create(t *testing.T) {
	mockSensorData := sensordata.CreateSensorDataDTO{
		SensorID:     1,
		Temperature:  25.5,
		Transparency: 8,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensordata.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectQuery("INSERT INTO sensor_data\\(sensor_id, temperature, transparency, created_at, updated_at\\) VALUES\\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
		WithArgs(mockSensorData.SensorID, mockSensorData.Temperature, mockSensorData.Transparency, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.Create(context.Background(), mockSensorData)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 1 {
		t.Errorf("unexpected sensor data ID, got: %d, want: %d", id, 1)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorDataRepository_AddDetectedSpiece(t *testing.T) {
	mockSensorDataID := 1
	mockSpieceID := 100

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := sensordata.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectExec("INSERT INTO detected_spieces\\(spiece_id, sensor_data_id\\) VALUES\\(\\$1, \\$2\\)").
		WithArgs(mockSpieceID, mockSensorDataID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.AddDetectedSpiece(context.Background(), mockSensorDataID, spiece.Spiece{ID: mockSpieceID})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
