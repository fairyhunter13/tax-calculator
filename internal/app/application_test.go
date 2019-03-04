// +build unit

package app

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/fairyhunter13/tax-calculator/internal/bill"
	mocksBill "github.com/fairyhunter13/tax-calculator/internal/bill/mocks"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	mocksTax "github.com/fairyhunter13/tax-calculator/internal/taxobj/mocks"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	errMigrate = errors.New("Database is not connected")
)

func TestApp_Migrate(t *testing.T) {
	t.Parallel()
	type fields struct {
		config    *Config
		pool      *sql.DB
		billRepo  bill.Repository
		billUcase bill.Usecase
		taxRepo   taxobj.Repository
		taxUcase  taxobj.Usecase
		echoMux   *echo.Echo
	}
	tests := []struct {
		name    string
		fields  func() fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			fields: func() fields {
				allFields := fields{}
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("Migrate").Return(nil)
				billUcase := &mocksBill.Usecase{}
				billUcase.On("LoadData").Return(nil)
				allFields.billUcase = billUcase
				allFields.taxRepo = taxRepo
				return allFields
			},
			wantErr: false,
		},
		{
			name: "Tax Repository Migrate Error",
			fields: func() fields {
				allFields := fields{}
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("Migrate").Return(errMigrate)
				allFields.taxRepo = taxRepo
				return allFields
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			app := &App{
				config:    fields.config,
				pool:      fields.pool,
				billRepo:  fields.billRepo,
				billUcase: fields.billUcase,
				taxRepo:   fields.taxRepo,
				taxUcase:  fields.taxUcase,
				echoMux:   fields.echoMux,
			}
			if err := app.Migrate(); (err != nil) != tt.wantErr {
				t.Errorf("App.Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_ParseConfig(t *testing.T) {
	t.Parallel()
	type fields struct {
		config    *Config
		pool      *sql.DB
		billRepo  bill.Repository
		billUcase bill.Usecase
		taxRepo   taxobj.Repository
		taxUcase  taxobj.Usecase
		echoMux   *echo.Echo
	}
	type args struct {
		configPath string
		appConfig  *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "Parse Config Succeed",
			fields: fields{},
			args: args{
				configPath: "../../configs/config.ini",
				appConfig:  &Config{},
			},
			wantErr: false,
		},
		{
			name:   "Parse Config Error",
			fields: fields{},
			args: args{
				configPath: "",
				appConfig:  &Config{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				config:    tt.fields.config,
				pool:      tt.fields.pool,
				billRepo:  tt.fields.billRepo,
				billUcase: tt.fields.billUcase,
				taxRepo:   tt.fields.taxRepo,
				taxUcase:  tt.fields.taxUcase,
				echoMux:   tt.fields.echoMux,
			}
			if err := app.ParseConfig(tt.args.configPath, tt.args.appConfig); (err != nil) != tt.wantErr {
				t.Errorf("App.ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_SetConfig(t *testing.T) {
	t.Parallel()
	type fields struct {
		config    *Config
		pool      *sql.DB
		billRepo  bill.Repository
		billUcase bill.Usecase
		taxRepo   taxobj.Repository
		taxUcase  taxobj.Usecase
		echoMux   *echo.Echo
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name:   "Set Application's Config",
			fields: fields{},
			args: args{
				config: new(Config),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				config:    tt.fields.config,
				pool:      tt.fields.pool,
				billRepo:  tt.fields.billRepo,
				billUcase: tt.fields.billUcase,
				taxRepo:   tt.fields.taxRepo,
				taxUcase:  tt.fields.taxUcase,
				echoMux:   tt.fields.echoMux,
			}
			app.SetConfig(tt.args.config)
		})
	}
}

func TestApp_Init(t *testing.T) {
	t.Parallel()
	type fields struct {
		config    *Config
		pool      *sql.DB
		billRepo  bill.Repository
		billUcase bill.Usecase
		taxRepo   taxobj.Repository
		taxUcase  taxobj.Usecase
		echoMux   *echo.Echo
	}
	type args struct {
		pool *sql.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name:   "Initialize the application",
			fields: fields{},
			args: args{
				pool: new(sql.DB),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				config:    tt.fields.config,
				pool:      tt.fields.pool,
				billRepo:  tt.fields.billRepo,
				billUcase: tt.fields.billUcase,
				taxRepo:   tt.fields.taxRepo,
				taxUcase:  tt.fields.taxUcase,
				echoMux:   tt.fields.echoMux,
			}
			app.Init(tt.args.pool)
		})
	}
}

func TestApp_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		config    *Config
		pool      *sql.DB
		billRepo  bill.Repository
		billUcase bill.Usecase
		taxRepo   taxobj.Repository
		taxUcase  taxobj.Usecase
		echoMux   *echo.Echo
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		// TODO: Add test cases.
		{
			name: "Closing the application",
			fields: func() fields {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Error in starting the mocker: %s", err)
				}

				fields := fields{
					pool: db,
				}
				taxRepo := new(mocksTax.Repository)
				taxRepo.On("Close")
				fields.taxRepo = taxRepo
				return fields
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			app := &App{
				config:    fields.config,
				pool:      fields.pool,
				billRepo:  fields.billRepo,
				billUcase: fields.billUcase,
				taxRepo:   fields.taxRepo,
				taxUcase:  fields.taxUcase,
				echoMux:   fields.echoMux,
			}
			app.Close()
		})
	}
}

func TestNewApp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *App
	}{
		// TODO: Add test cases.
		{
			name: "Create the application",
			want: new(App),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, NewApp(), tt.want)
		})
	}
}
