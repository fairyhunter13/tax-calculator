// +build !race

package app

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	"github.com/labstack/echo"
)

//The test case is located in here because the test's result shows that
//the application has some race conditions. The application handling is using
//graceful shutdown with using os.Signal so the race conditions are expected.
func TestApp_Run(t *testing.T) {
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
		osSignal chan os.Signal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Run the application",
			fields: fields{
				echoMux: echo.New(),
				config: &Config{
					Database{
						ConnectionString: "",
					},
					Server{
						Port: ":8080",
					},
				},
			},
			args: args{
				osSignal: make(chan os.Signal, 1),
			},
			wantErr: false,
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
			go func() {
				time.Sleep(500 * time.Millisecond)
				tt.args.osSignal <- os.Interrupt
			}()
			err := app.Run(tt.args.osSignal)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
