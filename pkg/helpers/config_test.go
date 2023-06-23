package helpers

import (
	"errors"
	"github.com/jackc/pgx"
	"os"
	"reflect"
	"testing"
)

func TestParsePsqlConfig(t *testing.T) {
	host, port, db, user, password := "localhost", uint16(5432), "links_db", "postgres", "pgpwdlinks"

	tests := []struct {
		name    string
		want    *pgx.ConnPoolConfig
		wantErr bool
	}{
		{
			name: "Normal config",
			want: &pgx.ConnPoolConfig{
				ConnConfig: pgx.ConnConfig{
					Host:     host,
					Port:     port,
					Database: db,
					User:     user,
					Password: password,
				},
			},
			wantErr: false,
		},
		{
			name:    "No .env",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "Normal config" {
			normal := []byte("POSTGRES_HOST=localhost\nPOSTGRES_PORT=5432\nPOSTGRES_DB=links_db\nPOSTGRES_USER=postgres\nPOSTGRES_PASSWORD=pgpwdlinks")
			err := os.WriteFile(".env", normal, 0644)
			if err != nil {
				return
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePsqlConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePsqlConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePsqlConfig() got = %v, want %v", got, tt.want)
			}
		})

		err := os.Remove(".env")
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return
			}
		}
	}
}
