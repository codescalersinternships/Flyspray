package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	// adding verified user
	newUser := signupBody{
		Name:            "diaa",
		Email:           "diaabadr@gmail.com",
		Password:        "diaabadr",
		ConfirmPassword: "diaabadr",
	}

	AddUserToDB(t, newUser, &app)

	user, err := app.DB.GetUserByEmail(newUser.Email)

	assert.Nil(t, err)
	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)

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
			name: "passwords do not match",
			requestBody: map[string]string{
				"name":             "diaa",
				"email":            "diaabadr82@gmail.com",
				"password":         "diaabadr",
				"confirm_password": "diaa",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
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

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	// adding verified user
	newUser := signupBody{
		Name:            "diaa",
		Email:           "diaabadr@gmail.com",
		Password:        "diaabadr",
		ConfirmPassword: "diaabadr",
	}

	AddUserToDB(t, newUser, &app)

	user, err := app.DB.GetUserByEmail(newUser.Email)

	assert.Nil(t, err)
	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)

	newUser = signupBody{
		Name:            "diaa",
		Email:           "diaabadr82@gmail.com",
		Password:        "diaabadr",
		ConfirmPassword: "diaabadr",
	}

	AddUserToDB(t, newUser, &app)

	testCases := []struct {
		name               string
		requestBody        verifyBody
		expectedStatusCode int
	}{
		{
			name: "verify with wrong code",
			requestBody: verifyBody{
				VerificationCode: 12345,
				Email:            "diaabadr82@gmail.com",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "user already verified",
			requestBody: verifyBody{
				VerificationCode: 12345,
				Email:            "diaabadr@gmail.com",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "user not found",
			requestBody: verifyBody{
				VerificationCode: 1234,
				Email:            "diaabadr8@gmail.com",
			},
			expectedStatusCode: http.StatusNotFound,
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

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

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

	user, err := app.DB.GetUserByEmail(verifiedUser.Email)
	assert.Nil(t, err)

	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)

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

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	app.config = internal.Configuration{}
	app.config.JWT.Timeout = 15

	newUser := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
	}

	AddUserToDB(t, newUser, &app)

	user, err := app.DB.GetUserByEmail(newUser.Email)
	assert.Nil(t, err)

	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)

	t.Run("unauthorized user", func(t *testing.T) {
		body, err := json.Marshal(struct {
			Name string `json:"name"`
		}{
			Name: "omar",
		})
		assert.Nil(t, err)
		request, err := http.NewRequest(http.MethodPut, "/user", bytes.NewReader(body))

		assert.Nil(t, err)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	user.Password = "diaa"
	// login user
	reqBody := signinBody{
		Email:    user.Email,
		Password: user.Password,
	}

	token := SigninUser(t, reqBody, &app)
	t.Run("authorized user", func(t *testing.T) {
		requestBody := signupBody{
			Name: "Omar",
		}

		body, err := json.Marshal(requestBody)
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

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	app.config = internal.Configuration{}
	app.config.JWT.Timeout = 15

	newUser := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
	}

	AddUserToDB(t, newUser, &app)

	user, err := app.DB.GetUserByEmail(newUser.Email)

	assert.Nil(t, err)

	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)
	body, err := json.Marshal(user)
	assert.Nil(t, err)
	t.Run("unauthorized user", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))

		assert.Nil(t, err)

		res := httptest.NewRecorder()

		app.router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	user.Password = "diaa"
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

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	app.config = internal.Configuration{}
	app.config.JWT.Timeout = 15

	newUser := signupBody{
		Name:            "diaa",
		Email:           "diaa12345678912@gmail.com",
		Password:        "diaa",
		ConfirmPassword: "diaa",
	}

	AddUserToDB(t, newUser, &app)

	user, err := app.DB.GetUserByEmail(newUser.Email)
	assert.Nil(t, err)

	err = app.DB.VerifyUser(user.ID)
	assert.Nil(t, err)
	user.Password = "diaa"
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

func TestForgetPassword(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	testCases := []struct {
		name               string
		preCreatedUser     models.User
		input              emailInput
		expectedStatusCode int
	}{
		{
			name: "invalid input",
			preCreatedUser: models.User{
				Name:     "omar",
				Email:    "omar@gmail.com",
				Password: "123456!Abc",
			},
			input:              emailInput{},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "not exist",
			preCreatedUser: models.User{
				Name:     "omar",
				Email:    "omar@gmail.com",
				Password: "123456!Abc",
			},
			input:              emailInput{Email: "another@gmail.com"},
			expectedStatusCode: http.StatusNotFound,
		}, {
			name: "not verified",
			preCreatedUser: models.User{
				Name:     "omar",
				Email:    "omar@gmail.com",
				Password: "123456!Abc",
			},
			input:              emailInput{Email: "omar@gmail.com"},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "valid",
			preCreatedUser: models.User{
				Name:     "omar",
				Email:    "omar@gmail.com",
				Password: "123456!Abc",
				Verified: true,
			},
			input:              emailInput{Email: "omar@gmail.com"},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer app.DB.Client.Exec("DELETE FROM users")

			_, err := app.DB.CreateUser(tc.preCreatedUser)
			assert.Nil(t, err)

			body, err := json.Marshal(tc.input)
			assert.Nil(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user/forget_password", bytes.NewReader(body))
			assert.Nil(t, err)

			res := httptest.NewRecorder()
			app.router.ServeHTTP(res, req)

			assert.Equal(t, tc.expectedStatusCode, res.Code)
		})
	}
}

func TestVerifyForgetPassword(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(filepath.Join(dir, "flyspray.db"))
	assert.Nil(t, err)

	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.registerRoutes()

	timeout := 15

	testCases := []struct {
		name               string
		preCreatedUser     models.User
		input              verifyForgetPasswordBody
		expectedStatusCode int
	}{
		{
			name: "invalid input",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input:              verifyForgetPasswordBody{},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "not exist",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input: verifyForgetPasswordBody{
				Email:            "another@gmail.com",
				VerificationCode: 100001,
				Password:         "changed-123456!Abc",
				ConfirmPassword:  "changed-123456!Abc",
			},
			expectedStatusCode: http.StatusNotFound,
		}, {
			name: "not verified",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                false,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input: verifyForgetPasswordBody{
				Email:            "omar@gmail.com",
				VerificationCode: 100001,
				Password:         "changed-123456!Abc",
				ConfirmPassword:  "changed-123456!Abc",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "passwords do not match",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input: verifyForgetPasswordBody{
				Email:            "omar@gmail.com",
				VerificationCode: 100001,
				Password:         "changed-a-123456!Abc",
				ConfirmPassword:  "changed-b-123456!Abc",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "wrong verification code",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input: verifyForgetPasswordBody{
				Email:            "omar@gmail.com",
				VerificationCode: 100003,
				Password:         "changed-123456!Abc",
				ConfirmPassword:  "changed-123456!Abc",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "verification code expired",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now(),
			},
			input: verifyForgetPasswordBody{
				Email:            "omar@gmail.com",
				VerificationCode: 100001,
				Password:         "changed-123456!Abc",
				ConfirmPassword:  "changed-123456!Abc",
			},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name: "valid",
			preCreatedUser: models.User{
				Name:                    "omar",
				Email:                   "omar@gmail.com",
				Password:                "123456!Abc",
				Verified:                true,
				VerificationCode:        100001,
				VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
			},
			input: verifyForgetPasswordBody{
				Email:            "omar@gmail.com",
				VerificationCode: 100001,
				Password:         "changed-123456!Abc",
				ConfirmPassword:  "changed-123456!Abc",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer app.DB.Client.Exec("DELETE FROM users")

			_, err := app.DB.CreateUser(tc.preCreatedUser)
			assert.Nil(t, err)

			body, err := json.Marshal(tc.input)
			assert.Nil(t, err)

			req, err := http.NewRequest(http.MethodPut, "/user/forget_password/verify", bytes.NewReader(body))
			assert.Nil(t, err)

			res := httptest.NewRecorder()
			app.router.ServeHTTP(res, req)

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
