package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
	}{
		{
			name: "valid body",
			requestBody: map[string]string{
				"name":     "diaa",
				"email":    "diaabadr82@gmail.com",
				"password": "diaabadr",
			},
			expectedStatusCode: http.StatusCreated,
		}, {
			name: "invalid email",
			requestBody: map[string]string{
				"name":     "diaa",
				"email":    "diaabadr",
				"password": "diaabadr",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "missing field in body",
			requestBody: map[string]string{
				"name":  "diaa",
				"email": "diaabadr82@gmail.com",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "email already exists",
			requestBody: map[string]string{
				"name":     "diaa",
				"email":    "diaabadr82@gmail.com",
				"password": "diaabadr",
			},
			expectedStatusCode: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.requestBody)

			assert.Nil(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(body))

			assert.Nil(t, err)
			res := httptest.NewRecorder()

			app.router.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}
}

func TestVerify(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
	}{
		{
			name: "verify with wrong code",
			requestBody: map[string]string{
				"verification_code": "12345",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.requestBody)

			assert.Nil(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user/signup/verify", bytes.NewReader(body))

			assert.Nil(t, err)
			res := httptest.NewRecorder()

			app.router.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}
}

func TestSignin(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	requestBody := map[string]string{
		"name":     "diaa",
		"email":    "diaabadr82@gmail.com",
		"password": "diaabadr",
	}

	body, err := json.Marshal(requestBody)

	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(body))

	assert.Nil(t, err)
	res := httptest.NewRecorder()

	app.router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
	}{
		{
			name: "account not verified",
			requestBody: map[string]string{
				"email":    "diaabadr82@gmail.com",
				"password": "diaabadr",
			},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name: "wrong email",
			requestBody: map[string]string{
				"email":"diaa@gmail.com",
				"password":"diaabadr",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "wrong password",
			requestBody: map[string]string{
				"email":"diaabadr82@gmail.com",
				"password":"diaa",
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.requestBody)

			assert.Nil(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(body))

			assert.Nil(t, err)
			res := httptest.NewRecorder()

			app.router.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}
}
