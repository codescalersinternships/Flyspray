package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.setRoutes()

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		input              createProjectInput
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{},
			input:              createProjectInput{Name: "new project"},
			expectedStatusCode: http.StatusCreated,
		}, {
			name:               "repeated project name",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			input:              createProjectInput{Name: "new project"},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "empty project name",
			preCreatedProjects: []models.Project{},
			input:              createProjectInput{Name: ""},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "invalid format",
			preCreatedProjects: []models.Project{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.client.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.client.CreateProject(p)
				assert.Nil(t, err)
			}

			payload, err := json.Marshal(tc.input)
			assert.Nil(t, err)

			req, err := http.NewRequest("POST", "/project", bytes.NewBuffer(payload))
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.setRoutes()

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		input              updateProjectInput
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			input:              updateProjectInput{Name: "updated project", OwnerId: 1},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "no change",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			input:              updateProjectInput{Name: "new project", OwnerId: 1},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "repeated project name",
			preCreatedProjects: []models.Project{{Name: "new project1", OwnerId: 1}, {Name: "new project2", OwnerId: 1}},
			input:              updateProjectInput{Name: "new project2", OwnerId: 1},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "empty name",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			input:              updateProjectInput{Name: "", OwnerId: 1},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "invalid format",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "not found",
			input:              updateProjectInput{Name: "updated project", OwnerId: 1},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.client.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.client.CreateProject(p)
				assert.Nil(t, err)
			}

			payload, err := json.Marshal(tc.input)
			assert.Nil(t, err)

			req, err := http.NewRequest("PUT", "/project/1", bytes.NewBuffer(payload))
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestGetProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.setRoutes()

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "not found",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.client.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.client.CreateProject(p)
				assert.Nil(t, err)
			}

			req, err := http.NewRequest("GET", "/project/1", nil)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestGetProjects(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.setRoutes()

	t.Run("valid", func(t *testing.T) {
		defer app.client.Client.Exec("DELETE FROM projects")

		p := models.Project{Name: "new project", OwnerId: 1}
		_, err := app.client.CreateProject(p)
		assert.Nil(t, err)

		req, err := http.NewRequest("GET", "/project/filters?userid=10007", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.setRoutes()

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerId: 1}},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "valid",
			preCreatedProjects: []models.Project{},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.client.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.client.CreateProject(p)
				assert.Nil(t, err)
			}

			req, err := http.NewRequest("DELETE", "/project/1", nil)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
