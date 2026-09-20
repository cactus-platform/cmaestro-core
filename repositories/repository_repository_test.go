package repositories

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cactus-platform/cmaestro-core/models"
	"github.com/cactus-platform/cmaestro-core/storage/sql"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newRepositoryRepositoryTest(t *testing.T) (*RepositoryRepositoryImpl, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		db.Close()
		t.Fatalf("gorm.Open() error = %v", err)
	}
	return &RepositoryRepositoryImpl{db: &sql.Client{DB: gormDB}}, mock, func() { _ = db.Close() }
}

func TestRepositoryRepositoryRejectsNilInputs(t *testing.T) {
	repository, _, cleanup := newRepositoryRepositoryTest(t)
	defer cleanup()

	if err := repository.Create(context.Background(), nil); err == nil {
		t.Fatal("Create(nil) expected error")
	}
	if err := repository.CreateRevision(context.Background(), nil); err == nil {
		t.Fatal("CreateRevision(nil) expected error")
	}
	if err := repository.Update(context.Background(), nil); err == nil {
		t.Fatal("Update(nil) expected error")
	}
}

func TestRepositoryRepositoryGetNotFound(t *testing.T) {
	repository, mock, cleanup := newRepositoryRepositoryTest(t)
	defer cleanup()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "repositories" WHERE id = $1 ORDER BY "repositories"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := repository.Get(context.Background(), id)
	if got != nil || err != ErrRepositoryNotFound {
		t.Fatalf("Get() = %#v, %v; want nil, ErrRepositoryNotFound", got, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestRepositoryRepositoryExists(t *testing.T) {
	repository, mock, cleanup := newRepositoryRepositoryTest(t)
	defer cleanup()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "repositories" WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := repository.Exists(context.Background(), id)
	if err != nil || exists {
		t.Fatalf("Exists() = %v, %v; want false, nil", exists, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestRepositoryRepositoryCreateAndUpdate(t *testing.T) {
	repository, mock, cleanup := newRepositoryRepositoryTest(t)
	defer cleanup()
	repositoryValue := &models.Repository{
		ID:   uuid.New(),
		Name: "repository",
		Artifacts: []*models.Artifact{{
			ID:   uuid.New(),
			Name: "artifact",
		}},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "repositories"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "artifacts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err := repository.Create(context.Background(), repositoryValue); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repositoryValue.Artifacts[0].RepositoryID != repositoryValue.ID {
		t.Fatal("Create() did not assign repository ID to artifact")
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "repositories"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err := repository.Update(context.Background(), repositoryValue); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
