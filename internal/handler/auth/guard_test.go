package auth_test

import (
	"context"
	"testing"

	"github.com/Albert-Ti/go-diploma-tpl/internal/handler/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthUserID(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantID  int
		wantErr bool
	}{
		{
			name:    "Success",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, "123"),
			wantID:  123,
			wantErr: false,
		},
		{
			name:    "Empty userID",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, ""),
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "No userID in context",
			ctx:     context.Background(),
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "Invalid userID (not string)",
			ctx:     context.WithValue(context.Background(), auth.UserIDKey, 123),
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := auth.GetAuthUserID(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}
