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

func TestCreateComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.POST("/component", WrapFunc(app.createComponent))

	t.Run("Success", func(t *testing.T) {
		componentInput := createComponentInput{
			ProjectID: "1",
			Name:      "Test Component",
		}
		payload, err := json.Marshal(componentInput)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}
		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusCreated != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Bad Request", func(t *testing.T) {

		componentInput := createComponentInput{
			ProjectID: "1",
			Name:      "Test Component",
		}
		payload, err := json.Marshal(componentInput)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusBadRequest != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Bad Request", func(t *testing.T) {

		requestBody := []byte(`{}`)
		payload, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusBadRequest != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUpdateComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

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

	updatedComponent := models.Component{
		Name:   "Updated Component",
		UserID: "1",
	}
	requestBody, err := json.Marshal(updatedComponent)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("PUT", "/component/1", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Bad Request", func(t *testing.T) {

		req, err := http.NewRequest("PUT", "/component/1", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		if http.StatusBadRequest != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/component/999", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		if http.StatusNotFound != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestGetComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		// assert.Equal(t, tc.expectedStatusCode, w.Code)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/123", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusNotFound != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestGetComponents(t *testing.T) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.router.GET("/component/filters", WrapFunc(app.getComponents))

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/filters?project_id=1", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}
	})
}

func TestDeleteComponent(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

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
	app.router.DELETE("/component/:id", WrapFunc(app.deleteComponent))

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", "/component/1", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", "/component/123", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusNotFound != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}
