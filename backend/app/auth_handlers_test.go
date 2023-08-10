package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
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

	unverifiedUser := models.User{
		Name:     "diaa",
		Email:    "diaabadr82@gmail.com",
		Password: "diaabadr",
	}
	AddUserToDB(t, unverifiedUser, &app)

	verifiedUser := models.User{
		Name:     "Omar",
		Email:    "omar12345678912@gmail.com",
		Password: "omar",
		Verified: true,
	}

	AddUserToDB(t, verifiedUser, &app)

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
				"email":    "diaa@gmail.com",
				"password": "diaabadr",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "wrong password",
			requestBody: map[string]string{
				"email":    "diaabadr82@gmail.com",
				"password": "diaa",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "verified email",
			requestBody: map[string]string{
				"email":    "omar12345678912@gmail.com",
				"password": "omar",
			},
			expectedStatusCode: http.StatusOK,
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

func TestUpdateUser(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	user := models.User{
		Name:     "diaa",
		Email:    "diaa12345678912@gmail.com",
		Password: "diaa",
		Verified: true,
	}

	AddUserToDB(t, user, &app)

	body, err := json.Marshal(user)
	assert.Nil(t, err)
	t.Run("unauthorized user", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))

		assert.Nil(t, err)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	// login user
	authCookie, _ := SigninUser(t, user, &app)
	t.Run("authorized user", func(t *testing.T) {
		user = models.User{
			Name:     "Omar",
			Email:    "omar@gmail.com",
			Password: "omar ahmed",
		}

		body, err = json.Marshal(user)
		assert.Nil(t, err)
		request, err := http.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
		assert.Nil(t, err)
		request.AddCookie(authCookie)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)
		fmt.Println(res.Body)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
}

func TestGetUser(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	user := models.User{
		Name:     "diaa",
		Email:    "diaa12345678912@gmail.com",
		Password: "diaa",
		Verified: true,
	}

	AddUserToDB(t, user, &app)

	body, err := json.Marshal(user)
	assert.Nil(t, err)
	t.Run("unauthorized user", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))

		assert.Nil(t, err)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	// login user
	authCookie, _ := SigninUser(t, user, &app)

	t.Run("authorized user", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))
		assert.Nil(t, err)

		request.AddCookie(authCookie)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestRefreshToken(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	app.client.Migrate()
	app.setUserRoutes()

	user := models.User{
		Name:     "diaa",
		Email:    "diaa12345678912@gmail.com",
		Password: "diaa",
		Verified: true,
	}

	AddUserToDB(t, user, &app)

	_, body := SigninUser(t, user, &app)

	var responseBody struct {
		Data struct {
			RefreshToken string `json:"refresh_token"`
		}
	}
	err = json.Unmarshal(body.Bytes(), &responseBody)

	assert.Nil(t, err)

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
	}{
		{
			name: "invalid token",
			requestBody: map[string]string{
				"refresh_token": "token",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "valid token",
			requestBody: map[string]string{
				"refresh_token": responseBody.Data.RefreshToken,
			},
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.requestBody)

			assert.Nil(t, err)
			request, err := http.NewRequest(http.MethodPost, "/user/refresh_token", bytes.NewReader(body))

			assert.Nil(t, err)
			res := httptest.NewRecorder()

			app.router.ServeHTTP(res, request)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}

}

func AddUserToDB(t testing.TB, user models.User, app *App) {
	t.Helper()

	body, err := json.Marshal(user)

	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(body))

	assert.Nil(t, err)
	res := httptest.NewRecorder()

	app.router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func SigninUser(t testing.TB, user models.User, app *App) (*http.Cookie, bytes.Buffer) {
	t.Helper()

	body, err := json.Marshal(user)
	assert.Nil(t, err)
	request, err := http.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(body))
	assert.Nil(t, err)
	res := httptest.NewRecorder()

	app.router.ServeHTTP(res, request)

	assert.Equal(t, http.StatusOK, res.Code)

	return res.Result().Cookies()[0], *res.Body
}
