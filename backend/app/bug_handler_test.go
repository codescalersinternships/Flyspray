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
	router.POST("/bugs", func(ctx *gin.Context) {
		_, response := app.createBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	t.Run("create new bug successfully", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)
	})

	t.Run("failed to create new bug", func(t *testing.T) {
		bugData := models.Bug{
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)
	})
}

func TestGetBugs(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app, err := NewApp(path)
	assert.NoError(t, err)

	// Create a new HTTP request
	router := gin.Default()
	router.POST("/bugs", func(ctx *gin.Context) {
		_, response := app.createBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	router.GET("/bugs/filters", func(ctx *gin.Context) {
		_, response := app.getBugs(ctx)
		ctx.JSON(response.Status(), response)
	})

	t.Run("get all bugs successfully", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/bugs/filters", nil)
		assert.NoError(t, err)

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

	router.POST("/bugs", func(ctx *gin.Context) {
		_, response := app.createBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	router.GET("/bugs/:id", func(ctx *gin.Context) {
		_, response := app.getSpecificBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	t.Run("get bug successfully", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bugs/1", nil)
		assert.NoError(t, err)

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

		req, err := http.NewRequest("POST", "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bugs/5", nil)
		assert.NoError(t, err)

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

		req, err := http.NewRequest(http.MethodPost, "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bugs", nil)
		assert.NoError(t, err)

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

	router.POST("/bugs", func(ctx *gin.Context) {
		_, response := app.createBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	router.PUT("/bugs/:id", func(ctx *gin.Context) {
		_, response := app.updateBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	t.Run("update bug successfully", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bugs", bytes.NewReader(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"user_id": "13", "component_id": 13}`)

		request, err := http.NewRequest("PUT", "/bugs/1", bytes.NewReader(updatedRequestBody))
		assert.NoError(t, err)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug is not found", func(t *testing.T) {
		updatedRequestBody := []byte(`{"user_id": "13", "component_id": 13}`)

		request, err := http.NewRequest("PUT", "/bugs/4", bytes.NewReader(updatedRequestBody))
		assert.NoError(t, err)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("failed to update bug as it is invalid", func(t *testing.T) {
		requestBody := []byte(`{"id":50,"user_id": "12", "component_id":15}`)
		request, err := http.NewRequest("POST", "/bugs", bytes.NewReader(requestBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":50,"user_id": 12, "component_id":14 }`)

		req, err := http.NewRequest("PUT", "/bugs/50", bytes.NewReader(updatedRequestBody))
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

	router.POST("/bugs", func(ctx *gin.Context) {
		_, response := app.createBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	router.DELETE("/bugs/:id", func(ctx *gin.Context) {
		_, response := app.deleteBug(ctx)
		ctx.JSON(response.Status(), response)
	})

	t.Run("delete bug successfully", func(t *testing.T) {
		requestBody := []byte(`{"id":51,"user_id": "12", "component_id":15}`)
		request, err := http.NewRequest("POST", "/bugs", bytes.NewReader(requestBody))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		req, err := http.NewRequest("DELETE", "/bugs/51", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug is not found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/bugs/50", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug ID is not given", func(t *testing.T) {
		request, err := http.NewRequest("DELETE", "/bugs/", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})
}
