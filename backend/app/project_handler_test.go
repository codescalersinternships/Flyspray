package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", "1")
		ctx.Next()
	})
	app.router.POST("/project", WrapFunc(app.createProject))

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
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
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
			defer app.DB.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.DB.CreateProject(p)
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

			// test changes in db
			if tc.expectedStatusCode == http.StatusCreated {
				p, err := app.DB.GetProject("1")
				assert.Nil(t, err)
				assert.Equal(t, p.Name, tc.input.Name)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", "1")
		ctx.Next()
	})
	app.router.PUT("/project/:id", WrapFunc(app.updateProject))

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		input              updateProjectInput
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			input:              updateProjectInput{Name: "updated project", OwnerID: "2"},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "no change",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			input:              updateProjectInput{Name: "new project", OwnerID: "1"},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "repeated project name",
			preCreatedProjects: []models.Project{{Name: "new project1", OwnerID: "1"}, {Name: "new project2", OwnerID: "1"}},
			input:              updateProjectInput{Name: "new project2", OwnerID: "1"},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "empty name",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			input:              updateProjectInput{Name: "", OwnerID: "1"},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "invalid format",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			expectedStatusCode: http.StatusBadRequest,
		}, {
			name:               "not found",
			input:              updateProjectInput{Name: "updated project", OwnerID: "1"},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.DB.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.DB.CreateProject(p)
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

			// test changes in db
			if tc.expectedStatusCode == http.StatusOK {
				p, err := app.DB.GetProject("1")
				assert.Nil(t, err)
				assert.Equal(t, p.Name, tc.input.Name)
				assert.Equal(t, p.OwnerID, tc.input.OwnerID)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.GET("/project/:id", WrapFunc(app.getProject))

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "not found",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.DB.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.DB.CreateProject(p)
				assert.Nil(t, err)
			}

			req, err := http.NewRequest("GET", "/project/1", nil)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			// test response body with db
			if tc.expectedStatusCode == http.StatusOK {
				// remove all whitespace from got
				got := strings.ReplaceAll(w.Body.String(), " ", "")
				got = strings.ReplaceAll(got, "\n", "")

				p, err := app.DB.GetProject("1")
				assert.Nil(t, err)
				resp := ResponseMsg{
					Message: "project is retrieved successfully",
					Data:    p,
				}
				EncodedWant, err := json.Marshal(resp)
				assert.Nil(t, err)
				// remove all whitespace from want
				want := strings.ReplaceAll(string(EncodedWant), " ", "")

				assert.Equal(t, want, got)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.GET("/project/filters", WrapFunc(app.getProjects))

	t.Run("valid", func(t *testing.T) {
		defer app.DB.Client.Exec("DELETE FROM projects")

		p := models.Project{Name: "new project", OwnerID: "10007"}
		_, err := app.DB.CreateProject(p)
		assert.Nil(t, err)

		req, err := http.NewRequest("GET", "/project/filters?userid=10007", nil)
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// test response body with db

		// remove all whitespace from got
		got := strings.ReplaceAll(w.Body.String(), " ", "")
		got = strings.ReplaceAll(got, "\n", "")

		ps, err := app.DB.FilterProjects("10007", "new project", "")
		assert.Nil(t, err)
		resp := ResponseMsg{
			Message: "projects are retrieved successfully",
			Data:    ps,
		}
		EncodedWant, err := json.Marshal(resp)
		assert.Nil(t, err)
		// remove all whitespace from want
		want := strings.ReplaceAll(string(EncodedWant), " ", "")

		assert.Equal(t, want, got)
	})
}

func TestDeleteProject(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", "1")
		ctx.Next()
	})
	app.router.DELETE("/project/:id", WrapFunc(app.deleteProject))

	tests := []struct {
		name               string
		preCreatedProjects []models.Project
		expectedStatusCode int
	}{
		{
			name:               "valid",
			preCreatedProjects: []models.Project{{Name: "new project", OwnerID: "1"}},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "not found",
			preCreatedProjects: []models.Project{},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer app.DB.Client.Exec("DELETE FROM projects")

			for _, p := range tc.preCreatedProjects {
				_, err := app.DB.CreateProject(p)
				assert.Nil(t, err)
			}

			req, err := http.NewRequest("DELETE", "/project/1", nil)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			// test changes in db
			if tc.expectedStatusCode == http.StatusOK {
				_, err := app.DB.GetProject("1")
				assert.Equal(t, gorm.ErrRecordNotFound, err)
			}
		})
	}
}
