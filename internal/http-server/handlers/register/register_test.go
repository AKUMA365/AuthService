package register

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"AuthServis/internal/storage"

	"github.com/stretchr/testify/require"
)

type MockUserSaver struct {
	SaveUserFunc func(ctx context.Context, email string, passHash []byte) (int64, error)
}

func (m *MockUserSaver) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	return m.SaveUserFunc(ctx, email, passHash)
}

func TestRegisterHandler(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		mockError error
		mockUID   int64
		wantCode  int
		wantError string
	}{
		{
			name:      "Success",
			body:      `{"email": "test@example.com", "password": "passwordsuper"}`,
			mockError: nil,
			mockUID:   1,
			wantCode:  http.StatusOK,
		},
		{
			name:      "Empty Password",
			body:      `{"email": "test@example.com", "password": ""}`,
			mockError: nil,
			wantCode:  http.StatusBadRequest,
			wantError: "invalid request fields",
		},
		{
			name:      "Invalid Email",
			body:      `{"email": "not-email", "password": "password123"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "invalid request fields",
		},
		{
			name:      "User Already Exists",
			body:      `{"email": "exists@example.com", "password": "password123"}`,
			mockError: storage.ErrUserExists,
			wantCode:  http.StatusConflict,
			wantError: "user already exists",
		},
		{
			name:      "Database Error",
			body:      `{"email": "fail@example.com", "password": "password123"}`,
			mockError: errors.New("something went wrong"),
			wantCode:  http.StatusInternalServerError,
			wantError: "failed to save user",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userSaver := &MockUserSaver{}
			userSaver.SaveUserFunc = func(ctx context.Context, email string, passHash []byte) (int64, error) {
				return tc.mockUID, tc.mockError
			}

			logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))

			handler := New(logger, userSaver)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.wantCode, rr.Code)

			if tc.wantError != "" {
				require.Contains(t, rr.Body.String(), tc.wantError)
			}
		})
	}
}
