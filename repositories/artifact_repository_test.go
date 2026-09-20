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

func newArtifactRepositoryTest(t *testing.T) (*ArtifactRepositoryImpl, sqlmock.Sqlmock, func()) {
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
	return &ArtifactRepositoryImpl{db: &sql.Client{DB: gormDB}}, mock, func() { _ = db.Close() }
}

func TestArtifactRepositoryRejectsNilInputs(t *testing.T) {
	repository, _, cleanup := newArtifactRepositoryTest(t)
	defer cleanup()

	if err := repository.CreateArtifact(context.Background(), nil); err == nil {
		t.Fatal("CreateArtifact(nil) expected error")
	}
	if err := repository.UpdateArtifact(context.Background(), nil); err == nil {
		t.Fatal("UpdateArtifact(nil) expected error")
	}
}

func TestArtifactRepositoryGetNotFound(t *testing.T) {
	repository, mock, cleanup := newArtifactRepositoryTest(t)
	defer cleanup()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artifacts" WHERE id = $1 ORDER BY "artifacts"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := repository.GetArtifact(context.Background(), id)
	if got != nil || err != ErrArtifactNotFound {
		t.Fatalf("GetArtifact() = %#v, %v; want nil, ErrArtifactNotFound", got, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestArtifactRepositoryExists(t *testing.T) {
	repository, mock, cleanup := newArtifactRepositoryTest(t)
	defer cleanup()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "artifacts" WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repository.ArtifactExists(context.Background(), id)
	if err != nil || !exists {
		t.Fatalf("ArtifactExists() = %v, %v; want true, nil", exists, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestArtifactRepositoryCreateAndUpdate(t *testing.T) {
	repository, mock, cleanup := newArtifactRepositoryTest(t)
	defer cleanup()
	artifact := &models.Artifact{ID: uuid.New(), Name: "artifact"}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "artifacts" WHERE id = \$1`).
		WithArgs(artifact.ID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(`INSERT INTO "artifacts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err := repository.CreateArtifact(context.Background(), artifact); err != nil {
		t.Fatalf("CreateArtifact() error = %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "artifacts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err := repository.UpdateArtifact(context.Background(), artifact); err != nil {
		t.Fatalf("UpdateArtifact() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
