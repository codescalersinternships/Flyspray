package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

	// adding verified user
	user := signupBody{
		Name:            "diaa",
		Email:           "diaabadr@gmail.com",
		Password:        "diaabadr",
		ConfirmPassword: "diaabadr",
		Verified:        true,
	}

	AddUserToDB(t, user, &app)

	testCases := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
	}{
		{
			name: "valid body",
			requestBody: map[string]string{
				"name":             "diaa",
				"email":            "diaabadr82@gmail.com",
				"password":         "diaabadr",
				"confirm_password": "diaabadr",
			},
			expectedStatusCode: http.StatusCreated,
		}, {
			name: "invalid email",
			requestBody: map[string]string{
				"name":             "diaa",
				"email":            "diaabadr",
				"password":         "diaabadr",
				"confirm_password": "diaabadr",
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
			name: "email already exists and not verified",
			requestBody: map[string]string{
				"name":             "diaa",
				"email":            "diaabadr82@gmail.com",
				"password":         "diaabadr",
				"confirm_password": "diaabadr",
			},
			expectedStatusCode: http.StatusOK,
		}, {
			name: "email already exists and verified",
			requestBody: map[string]string{
				"name":             "diaa",
				"email":            "diaabadr@gmail.com",
				"password":         "diaabadr",
				"confirm_password": "diaabadr",
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
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

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
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

	unverifiedUser := signupBody{
		Name:            "diaa",
		Email:           "diaabadr82@gmail.com",
		Password:        "diaabadr",
		ConfirmPassword: "diaabadr",
	}
	AddUserToDB(t, unverifiedUser, &app)

	verifiedUser := signupBody{
		Name:            "Omar",
		Email:           "omar12345678912@gmail.com",
		Password:        "omar",
		ConfirmPassword: "omar",
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
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

	user := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
		Verified:        true,
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
	reqBody := signinBody{
		Email:    user.Email,
		Password: user.Password,
	}
	token := SigninUser(t, reqBody, &app)
	t.Run("authorized user", func(t *testing.T) {
		requestBody := signupBody{
			Name:     "Omar",
			Email:    "omar@gmail.com",
			Password: "omar ahmed",
		}

		body, err = json.Marshal(requestBody)
		assert.Nil(t, err)
		request, err := http.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))
		assert.Nil(t, err)

		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

	user := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
		Verified:        true,
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

	reqBody := signinBody{
		Email:    user.Email,
		Password: user.Password,
	}
	// login user
	token := SigninUser(t, reqBody, &app)

	t.Run("authorized user", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))
		assert.Nil(t, err)

		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)
	app.setRoutes()

	user := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
		Verified:        true,
	}

	AddUserToDB(t, user, &app)
	reqBody := signinBody{
		Email:    user.Email,
		Password: user.Password,
	}

	token := SigninUser(t, reqBody, &app)

	testCases := []struct {
		name               string
		token              string
		expectedStatusCode int
	}{
		{
			name:               "invalid token",
			token:              "token",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "valid token",
			token:              token,
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodPost, "/user/refresh_token", nil)

			assert.Nil(t, err)

			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			res := httptest.NewRecorder()

			app.router.ServeHTTP(res, request)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}

}

func AddUserToDB(t testing.TB, user signupBody, app *App) {
	t.Helper()

	body, err := json.Marshal(user)

	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(body))

	assert.Nil(t, err)
	res := httptest.NewRecorder()

	app.router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func SigninUser(t testing.TB, user signinBody, app *App) string {
	t.Helper()

	body, err := json.Marshal(user)
	assert.Nil(t, err)
	request, err := http.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(body))
	assert.Nil(t, err)
	res := httptest.NewRecorder()

	app.router.ServeHTTP(res, request)

	assert.Equal(t, http.StatusOK, res.Code)

	var responseBody struct {
		Data struct {
			AccessToken string `json:"access_token"`
		}
	}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	assert.Nil(t, err)
	return responseBody.Data.AccessToken
}
