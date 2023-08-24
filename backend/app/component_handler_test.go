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
	"gorm.io/gorm"
)

func TestCreateComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(dbPath)
	assert.Nil(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.POST("/component", WrapFunc(app.createComponent))

	createProject := models.Project{
		Name:    "Project",
		OwnerID: "1",
	}

	createProject, err = app.DB.CreateProject(createProject)
	assert.Nil(t, err)

	member := models.Member{ProjectID: int(createProject.ID), Admin: true, UserID: createProject.OwnerID}
	err = app.DB.CreateNewMember(member)
	assert.Nil(t, err)

	t.Run("Success", func(t *testing.T) {
		componentInput := createComponentInput{
			Name:      "Test Component",
			ProjectID: "1",
		}
		payload, err := json.Marshal(componentInput)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}
		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		c, err := app.DB.GetComponent("1")
		assert.Nil(t, err)
		assert.Equal(t, c.Name, componentInput.Name)

	})

	t.Run("Bad Request, name is already exist", func(t *testing.T) {

		componentInput := createComponentInput{
			ProjectID: "1",
			Name:      "Test Component",
		}
		payload, err := json.Marshal(componentInput)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Bad Request, empty body", func(t *testing.T) {

		requestBody := []byte(`{}`)
		payload, err := json.Marshal(requestBody)
		assert.Nil(t, err)

		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(dbPath)
	assert.Nil(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.PUT("/component/:id", WrapFunc(app.updateComponent))

	createComponent := models.Component{
		ProjectID: "1",
		Name:      "New Component",
		UserID:    "1",
	}

	_, err = app.DB.CreateComponent(createComponent)
	assert.Nil(t, err)

	createProject := models.Project{
		Name:    "New Project",
		OwnerID: "1",
	}

	_, err = app.DB.CreateProject(createProject)
	assert.Nil(t, err)

	updatedComponent := models.Component{
		Name: "Updated Component",
	}
	requestBody, err := json.Marshal(updatedComponent)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("PUT", "/component/1", bytes.NewBuffer(requestBody))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		c, err := app.DB.GetComponent("1")
		assert.Nil(t, err)
		assert.Equal(t, c.Name, updatedComponent.Name)
	})

	t.Run("Bad Request", func(t *testing.T) {

		req, err := http.NewRequest("PUT", "/component/1", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/component/999", bytes.NewBuffer(requestBody))
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(dbPath)
	assert.Nil(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.GET("/component/:id", WrapFunc(app.getComponent))

	createComponent := models.Component{
		ProjectID: "1",
		Name:      "New Component",
		UserID:    "1",
	}

	_, err = app.DB.CreateComponent(createComponent)
	assert.Nil(t, err)

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/1", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/123", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetComponents(t *testing.T) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(dbPath)
	assert.Nil(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.GET("/component/filters", WrapFunc(app.getComponents))

	createComponent := models.Component{
		ProjectID: "1",
		Name:      "test",
		UserID:    "1",
	}

	_, err = app.DB.CreateComponent(createComponent)
	assert.Nil(t, err)

	t.Run("Success, retrieve components by project id", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/filters?project_id=1", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success, retrieve component by name", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/filters?name=test", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success, retrieve component by name and project id", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/filters?name=test&project_id=1", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestDeleteComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(dbPath)
	assert.Nil(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.GET("/component/:id", WrapFunc(app.getComponent))

	createComponent := models.Component{
		ProjectID: "1",
		Name:      "New Component",
		UserID:    "1",
	}

	_, err = app.DB.CreateComponent(createComponent)
	assert.Nil(t, err)

	createProject := models.Project{
		Name:    "New Project",
		OwnerID: "1",
	}

	_, err = app.DB.CreateProject(createProject)
	assert.Nil(t, err)

	app.router.DELETE("/component/:id", WrapFunc(app.deleteComponent))

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", "/component/1", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		_, err = app.DB.GetComponent("1")

		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Not Found", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", "/component/123", nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
