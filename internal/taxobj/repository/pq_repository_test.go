// +build unit

package repository

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

const (
	regexQueryInsert = `
		INSERT INTO tax_object
			(.+)
		VALUES
			(.+)
		RETURNING id
	`
	regexQuerySelectAll = `
		SELECT 
			\*
		FROM
			tax_object
	`
	regexQuerySelectOne = `
		SELECT
			id
		FROM
			tax_object
		LIMIT 1
	`
	regexQueryCreateTable = `
		CREATE TABLE tax_object (.+)
	`
)

var (
	errPreparingStatement = errors.New("Error preparing the statement")
	errQuerying           = errors.New("Error in querying rows")
	errRelationNotExist   = errors.New("Relation still doesn't exist")
)

func TestPqRepository_GetAll(t *testing.T) {
	t.Parallel()
	const logFail = `[TestPqRepository_GetAll] %s: %s`
	tests := []struct {
		name           string
		customFunc     func() (*PqRepository, sqlmock.Sqlmock, *sql.DB)
		wantTaxObjects []taxobj.TaxObject
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				resultRow := sqlmock.NewRows([]string{"id", "name", "tax_code", "price"})
				resultRow.AddRow(1, "MACD", 1, 20000)
				//Init the mock!
				mock.ExpectPrepare(regexQuerySelectAll)
				mock.ExpectQuery(regexQuerySelectAll).
					WillReturnRows(resultRow)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantTaxObjects: []taxobj.TaxObject{
				taxobj.TaxObject{
					ID:      1,
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			wantErr: false,
		},
		{
			name: "Error preparing the statement",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					log.Fatalf(logFail, "Error starting the mocker", err)
				}
				//Init the mock!
				mock.ExpectPrepare(regexQuerySelectAll).WillReturnError(errPreparingStatement)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantTaxObjects: []taxobj.TaxObject{},
			wantErr:        true,
		},
		{
			name: "Error querying rows",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					log.Fatalf(logFail, "Error starting the mocker", err)
				}
				//Init the mock!
				mock.ExpectPrepare(regexQuerySelectAll)
				mock.ExpectQuery(regexQuerySelectAll).
					WillReturnError(errQuerying)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantTaxObjects: []taxobj.TaxObject{},
			wantErr:        true,
		},
		{
			name: "Error scanning the data",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					log.Fatalf(logFail, "Error starting the mocker", err)
				}
				resultRow := sqlmock.NewRows([]string{"hello"})
				resultRow.AddRow(1)
				//Init the mock!
				mock.ExpectPrepare(regexQuerySelectAll)
				mock.ExpectQuery(regexQuerySelectAll).
					WillReturnRows(resultRow)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantTaxObjects: []taxobj.TaxObject{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, db := tt.customFunc()
			defer db.Close()
			gotTaxObjects, err := repo.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("PqRepository.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTaxObjects, tt.wantTaxObjects) {
				t.Errorf("PqRepository.GetAll() = %v, want %v", gotTaxObjects, tt.wantTaxObjects)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("PqRepository.GetAll() mock expectation were not met: %s", err)
			}
		})
	}
}

func TestPqRepository_Create(t *testing.T) {
	t.Parallel()
	const logFail = `[TestPqRepository_Create] %s: %s`
	type args struct {
		taxObj *taxobj.TaxObject
	}
	tests := []struct {
		name       string
		customFunc func() (*PqRepository, sqlmock.Sqlmock, *sql.DB)
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				resultRow := sqlmock.NewRows([]string{"id"})
				resultRow.AddRow(1)
				//Init the mock!
				mock.ExpectPrepare(regexQueryInsert)
				mock.ExpectQuery(regexQueryInsert).
					WithArgs("MACD", 1, float64(20000)).
					WillReturnRows(resultRow)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			args: args{
				taxObj: &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			wantErr: false,
		},
		{
			name: "Error in preparing statement",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				//Init the mock!
				mock.ExpectPrepare(regexQueryInsert).
					WillReturnError(errPreparingStatement)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			args: args{
				taxObj: &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, db := tt.customFunc()
			defer db.Close()
			err := repo.Create(tt.args.taxObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("PqRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("PqRepository.Create() mock expectation were not met: %s", err)
			}
		})
	}
}

func TestPqRepository_Migrate(t *testing.T) {
	t.Parallel()
	const logFail = `[TestPqRepository_Migrate] %s: %s`
	tests := []struct {
		name       string
		customFunc func() (*PqRepository, sqlmock.Sqlmock, *sql.DB)
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Table has already exist",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				resultRow := sqlmock.NewRows([]string{"id"})
				resultRow.AddRow(1)
				//Init the mock!
				mock.ExpectQuery(regexQuerySelectOne).
					WillReturnRows(resultRow)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantErr: false,
		},
		{
			name: "Create table for the first time",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				//Init the mock!
				mock.ExpectQuery(regexQuerySelectOne).
					WillReturnError(errRelationNotExist)
				mock.ExpectExec(regexQueryCreateTable)

				repo := NewPqRepository(db)
				return repo.(*PqRepository), mock, db
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, db := tt.customFunc()
			defer db.Close()
			err := repo.Migrate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PqRepository.Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("PqRepository.Migrate() mock expectation were not met: %s", err)
			}
		})
	}
}

func TestPqRepository_Close(t *testing.T) {
	t.Parallel()
	const logFail = `[TestPqRepository_Migrate] %s: %s`
	tests := []struct {
		name       string
		customFunc func() (*PqRepository, sqlmock.Sqlmock, *sql.DB)
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			customFunc: func() (*PqRepository, sqlmock.Sqlmock, *sql.DB) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Errorf(logFail, "Error starting the mocker", err)
				}
				//Init the mock!
				mock.ExpectPrepare(regexQueryInsert)
				mock.ExpectPrepare(regexQuerySelectAll)

				repo := NewPqRepository(db).(*PqRepository)
				stmt, err := db.Prepare(queryInsert)
				if err != nil {
					t.Errorf(logFail, "Error preparing the statement", err)
				}
				repo.statement.insert = stmt
				stmt, err = db.Prepare(querySelectAll)
				if err != nil {
					t.Errorf(logFail, "Error preparing the statement", err)
				}
				repo.statement.selectAll = stmt
				return repo, mock, db
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, db := tt.customFunc()
			defer db.Close()
			repo.Close()
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("PqRepository.Close() mock expectation were not met: %s", err)
			}
		})
	}
}

func TestNewPqRepository(t *testing.T) {
	t.Parallel()
	type args struct {
		pool *sql.DB
	}
	db, _, err := sqlmock.New()
	if err != nil {
		assert.Errorf(t, err, "Error starting the mock: %s", err)
	}
	defer db.Close()
	tests := []struct {
		name string
		args args
		want taxobj.Repository
	}{
		// TODO: Add test cases.
		{
			name: "Init Tax Object Pq Repository",
			args: args{
				pool: db,
			},
			want: &PqRepository{
				pool: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, NewPqRepository(tt.args.pool), tt.want)
		})
	}
}
