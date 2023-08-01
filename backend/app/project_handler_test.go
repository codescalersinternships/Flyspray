package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	dbFilePath := "./test.db"
	app, err := NewApp(dbFilePath)
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.POST("project", func(ctx *gin.Context) {
		createProject(ctx, app.client)
	})

	t.Run("valid", func(t *testing.T) {
		input := createProjectInput{Name: "new project"}
		payload, err := json.Marshal(input)
		assert.Nil(t, err)

		req, err := http.NewRequest("POST", "/project", bytes.NewBuffer(payload))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid", func(t *testing.T) {
		input := createProjectInput{Name: ""}
		payload, err := json.Marshal(input)
		assert.Nil(t, err)

		req, err := http.NewRequest("POST", "/project", bytes.NewBuffer(payload))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetProject(t *testing.T) {
	dbFilePath := "./test.db"
	app, err := NewApp(dbFilePath)
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.GET("project/:id", func(ctx *gin.Context) {
		getProject(ctx, app.client)
	})

	t.Run("valid", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/project/1", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/project/-1", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetProjects(t *testing.T) {
	dbFilePath := "./test.db"
	app, err := NewApp(dbFilePath)
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.GET("project/filters", func(ctx *gin.Context) {
		getProjects(ctx, app.client)
	})

	t.Run("valid", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/project/filters?userid=10007", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteProject(t *testing.T) {
	dbFilePath := "./test.db"
	app, err := NewApp(dbFilePath)
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.DELETE("project/:id", func(ctx *gin.Context) {
		deleteProject(ctx, app.client)
	})

	t.Run("valid", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/project/1", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
