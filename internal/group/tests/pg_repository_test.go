package group

import (
	"context"
	"sensors-generator/internal/group"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// func Test_SensorGroupRepository_FindSpiecesInGroup(t *testing.T) {
// 	mockGroupName := "alpha"

// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	logging.Init("trace", true)

// 	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

// 	expectedQuery := `SELECT s.id, s.name, s.created_at, s.updated_at, COUNT(s.id)
// 	FROM sensor_groups as sg
// 	JOIN sensors sens ON sg.id = sens.group_id
// 	JOIN sensor_data sd ON sens.id = sd.sensor_id
// 	JOIN detected_spieces ds ON sd.id = ds.sensor_data_id
// 	JOIN spieces s ON s.id = ds.spiece_id
// 	WHERE sg.name = $1
// 	GROUP BY s.id, s.name, s.created_at, s.updated_at`
// 	expectedRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "count"}).
// 		AddRow(1, "Species1", time.Now(), time.Now(), 5).
// 		AddRow(2, "Species2", time.Now(), time.Now(), 10)

// 	mock.ExpectQuery(expectedQuery).
// 		WithArgs(mockGroupName).
// 		WillReturnRows(expectedRows)

// 	speciesCounts, err := repo.FindSpiecesInGroup(context.Background(), mockGroupName, group.SensorGroupFilters{})

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if len(speciesCounts) != 2 {
// 		t.Errorf("unexpected number of species found, got: %d, want: %d", len(speciesCounts), 2)
// 	}

// 	expectedSpecies := map[string]int{"Species1": 10, "Species2": 5}
// 	for sp, count := range speciesCounts {
// 		if expectedCount, ok := expectedSpecies[sp.Name]; ok {
// 			if count != expectedCount {
// 				t.Errorf("unexpected count for species %s, got: %d, want: %d", sp.Name, count, expectedCount)
// 			}
// 		} else {
// 			t.Errorf("unexpected species %s", sp.Name)
// 		}
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

func Test_SensorGroupRepository_FindOneByID(t *testing.T) {
	logging.Init("trace", true)
	mockSensorGroupID := 1
	mockSensorGroup := group.SensorGroup{
		ID:        mockSensorGroupID,
		Name:      "alpha",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM sensor_groups WHERE id=\\$1").
		WithArgs(mockSensorGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(mockSensorGroup.ID, mockSensorGroup.Name, mockSensorGroup.CreatedAt, mockSensorGroup.UpdatedAt))

	sensorGroup, err := repo.FindOneByID(context.Background(), mockSensorGroupID, group.SensorGroupFilters{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sensorGroup.ID != mockSensorGroup.ID {
		t.Errorf("unexpected sensor group ID, got: %d, want: %d", sensorGroup.ID, mockSensorGroup.ID)
	}
	if sensorGroup.Name != mockSensorGroup.Name {
		t.Errorf("unexpected sensor group name, got: %s, want: %s", sensorGroup.Name, mockSensorGroup.Name)
	}

	if !sensorGroup.CreatedAt.Equal(mockSensorGroup.CreatedAt) {
		t.Errorf("unexpected created at timestamp, got: %v, want: %v", sensorGroup.CreatedAt, mockSensorGroup.CreatedAt)
	}
	if !sensorGroup.UpdatedAt.Equal(mockSensorGroup.UpdatedAt) {
		t.Errorf("unexpected updated at timestamp, got: %v, want: %v", sensorGroup.UpdatedAt, mockSensorGroup.UpdatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorGroupRepository_FindAll(t *testing.T) {
	rows := []group.SensorGroup{
		{
			ID:        1,
			Name:      "alpha",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		{
			ID:        2,
			Name:      "beta",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	expectedRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})

	for _, row := range rows {
		expectedRows.AddRow(row.ID, row.Name, row.CreatedAt, row.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM sensor_groups").
		WillReturnRows(expectedRows)

	sensorGroups, err := repo.FindAll(context.Background(), group.SensorGroupFilters{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCount := 2
	if len(sensorGroups) != expectedCount {
		t.Errorf("unexpected number of sensor groups, got: %d, want: %d", len(sensorGroups), expectedCount)
	}

	for i := 0; i < expectedCount; i++ {
		if err != nil {
			t.Fatalf("failed to scan row: %v", err)
		}

		expectedSensorGroup := sensorGroups[i]

		if expectedSensorGroup.ID != rows[i].ID {
			t.Errorf("unexpected sensor group ID, got: %d, want: %d", expectedSensorGroup.ID, rows[i].ID)
		}
		if expectedSensorGroup.Name != rows[i].Name {
			t.Errorf("unexpected sensor group name, got: %s, want: %s", expectedSensorGroup.Name, rows[i].Name)
		}

		if !expectedSensorGroup.CreatedAt.Equal(rows[i].CreatedAt) {
			t.Errorf("unexpected created at timestamp, got: %v, want: %v", expectedSensorGroup.CreatedAt, rows[i].CreatedAt)
		}
		if !expectedSensorGroup.UpdatedAt.Equal(rows[i].UpdatedAt) {
			t.Errorf("unexpected updated at timestamp, got: %v, want: %v", expectedSensorGroup.UpdatedAt, rows[i].UpdatedAt)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorGroupRepository_FindAvgTransparencyInGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	groupName := "alpha"

	expectedTransparency := 50.0

	mock.ExpectQuery("SELECT AVG\\(sd\\.transparency\\) FROM sensor_groups as sg").
		WithArgs(groupName).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(expectedTransparency))

	transparency, err := repo.FindAvgTransparencyInGroup(context.Background(), groupName, group.SensorGroupFilters{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if float64(transparency) != expectedTransparency {
		t.Errorf("unexpected transparency value, got: %f, want: %f", float64(transparency), expectedTransparency)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorGroupRepository_FindAvgTemperatureInGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	groupName := "alpha"

	expectedTemperature := 25.0

	mock.ExpectQuery("SELECT AVG\\(sd\\.temperature\\) FROM sensor_groups as sg").
		WithArgs(groupName).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(expectedTemperature))

	temperature, err := repo.FindAvgTemperatureInGroup(context.Background(), groupName, group.SensorGroupFilters{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if float64(temperature) != expectedTemperature {
		t.Errorf("unexpected temperature value, got: %f, want: %f", float64(temperature), expectedTemperature)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SensorGroupRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := group.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	groupName := "alpha"

	mock.ExpectExec("INSERT INTO sensor_groups\\(name, created_at, updated_at\\)").
		WithArgs(groupName, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), group.CreateSensorGroupDTO{Name: groupName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
