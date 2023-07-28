package spiece

import (
	"context"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/logging"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_SpieceRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logging.Init("trace", true)

	repo := spiece.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	groupName := "alpha"
	mockRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Species1", time.Now(), time.Now()).
		AddRow(2, "Species2", time.Now(), time.Now())

	mock.ExpectQuery("SELECT s.id, s.name, s.created_at, s.updated_at FROM sensors as sens JOIN sensor_data sd ON sens.id=sd.sensor_id JOIN detected_spieces ds ON sd.id=ds.sensor_data_id JOIN spieces s ON s.id=ds.spiece_id WHERE sens.group_name=?").
		WithArgs(groupName).
		WillReturnRows(mockRows)

	spieces, err := repo.FindAll(context.Background(), spiece.SpieceFilters{GroupName: groupName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(spieces) != 2 {
		t.Errorf("unexpected number of spieces, got: %d, want: %d", len(spieces), 2)
	}

	expectedSpecies := []string{"Species1", "Species2"}
	for i, sp := range spieces {
		if sp.Name != expectedSpecies[i] {
			t.Errorf("unexpected species name, got: %s, want: %s", sp.Name, expectedSpecies[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SpieceRepository_FindOneByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := spiece.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	spieceID := 1
	mockRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(spieceID, "Species1", time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM spieces WHERE id=?").
		WithArgs(spieceID).
		WillReturnRows(mockRows)

	spiece, err := repo.FindOneByID(context.Background(), spieceID, spiece.SpieceFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if spiece.ID != spieceID {
		t.Errorf("unexpected spiece ID, got: %d, want: %d", spiece.ID, spieceID)
	}
	if spiece.Name != "Species1" {
		t.Errorf("unexpected species name, got: %s, want: %s", spiece.Name, "Species1")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_SpieceRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := spiece.NewPostgresqlRepository(db, logging.GetLogger(), nil)

	spieceName := "Species1"

	mock.ExpectExec("INSERT INTO spieces").
		WithArgs(spieceName, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), spiece.CreateSpieceDTO{Name: spieceName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
