package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()
	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})
	router.POST("/bug", WrapFunc(app.createBug))

	t.Run("create new bug successfully", func(t *testing.T) {
		bugInput := createBugInput{
			Summary:     "hello world!",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugInput)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)
	})

	t.Run("failed to create new bug", func(t *testing.T) {
		requestBody := []byte(`{"user_id": "99"}`)

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)
	})
}

func TestGetbug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()
	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})
	router.POST("/bug", WrapFunc(app.createBug))

	router.GET("/bug/filters", WrapFunc(app.getbug))

	t.Run("get all bug successfully", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/bug/filters", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})
}

func TestGetBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.POST("/bug", WrapFunc(app.createBug))

	router.GET("/bug/:id", WrapFunc(app.getSpecificBug))

	t.Run("get bug successfully", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug/1", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("get bug is not found", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug/5", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug ID is not given", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})
}

func TestUpdateBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.POST("/bug", WrapFunc(app.createBug))

	router.PUT("/bug/:id", WrapFunc(app.updateBug))

	t.Run("update bug successfully", func(t *testing.T) {
		bugInput := createBugInput{
			Summary:     "hello from test",
			ComponentID: 90,
		}

		wantedBug := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput.Summary,
			ComponentID: bugInput.ComponentID,
		}

		result := app.client.Client.Create(&wantedBug)
		assert.NoError(t, result.Error)

		bugUpdate := updateBugInput{
			Summary: "update bug",
		}

		payload, err := json.Marshal(bugUpdate)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		request, err := http.NewRequest("PUT", "/bug/10", bytes.NewBuffer(payload))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug is not found", func(t *testing.T) {

		bugUpdate := updateBugInput{
			Summary: "update bug",
		}

		payload, err := json.Marshal(bugUpdate)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		request, err := http.NewRequest("PUT", "/bug/50", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})

	t.Run("failed to update bug as it is invalid", func(t *testing.T) {
		requestBody := []byte(`{"id":50,"user_id": "99", "component_id":15}`)
		request, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":50,"user_id": "99", "component_id" }`)

		req, err := http.NewRequest("PUT", "/bug/33", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)
	})
}

func TestDeleteBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.DELETE("/bug/:id", WrapFunc(app.deleteBug))

	t.Run("delete bug successfully", func(t *testing.T) {
		bugInput := createBugInput{
			Summary:     "hello from test",
			ComponentID: 90,
		}

		wantedBug := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput.Summary,
			ComponentID: bugInput.ComponentID,
		}

		result := app.client.Client.Create(&wantedBug)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("DELETE", "/bug/10", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code, "got %d status code but want status code 201", recorder.Code)
	})

	t.Run("bug is not found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/bug/50", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug ID is not given", func(t *testing.T) {
		request, err := http.NewRequest("DELETE", "/bug/", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})
}
