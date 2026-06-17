package config_test

import (
	"flag"
	"os"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFlag(t *testing.T) {

	tests := []struct {
		name string
		args []string
		env  map[string]string
		want config.Config
	}{
		{
			name: "Add_ENV",
			args: []string{"test"},
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8888",
				"DATABASE_URI":           "postgres://test:test@localhost:5432/test",
				"JWT_SECRET_KEY":         "test_secret_key",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:8881",
			},
			want: config.Config{
				RunAddr:           "localhost:8888",
				DatabaseURI:       "postgres://test:test@localhost:5432/test",
				JWTSecretKey:      "test_secret_key",
				AccrualSystemAddr: "http://localhost:8881",
			},
		},
		{
			name: "Add_flag",
			args: []string{
				"test",
				"-a=localhost:9090",
				"-d=postgres://test:test@localhost:5432/test",
				"-r=http://localhost:9000",
			},
			env: map[string]string{},
			want: config.Config{
				RunAddr:           "localhost:9090",
				DatabaseURI:       "postgres://test:test@localhost:5432/test",
				JWTSecretKey:      "default_secret_key",
				AccrualSystemAddr: "http://localhost:9000",
			},
		},
		{
			name: "Default",
			args: []string{"test"},
			env:  map[string]string{},
			want: config.Config{
				RunAddr:           "localhost:8080",
				DatabaseURI:       "",
				JWTSecretKey:      "default_secret_key",
				AccrualSystemAddr: "http://localhost:8081",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			os.Args = tt.args

			config.ParseFlag()

			assert.Equal(t, tt.want.RunAddr, config.Envs.RunAddr)
			assert.Equal(t, tt.want.DatabaseURI, config.Envs.DatabaseURI)
		})
	}
}
