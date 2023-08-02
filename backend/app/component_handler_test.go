package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

func TestCreateComponent(t *testing.T) {

	databasePath := "./test.db"
	app, err := NewApp(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	app.router = gin.Default()
	app.router.POST("/component", app.CreateComponent)

	t.Run("Success", func(t *testing.T) {

		requestBody := []byte(`{"project_id": 1, "name": "Test Component"}`)
		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		if http.StatusCreated != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
		}

		var responseComponent models.Component
		err = json.Unmarshal(w.Body.Bytes(), &responseComponent)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		expectedComponent := models.Component{
			ProjectID: 1,
			Name:      "Test Component",
		}
		if expectedComponent.ProjectID != responseComponent.ProjectID {
			t.Errorf("Expected status code %d, but got %d", expectedComponent.ProjectID, responseComponent.ProjectID)
		}
		if expectedComponent.Name != responseComponent.Name {
			t.Errorf("Expected status code %s, but got %s", expectedComponent.Name, responseComponent.Name)
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

func TestGetComponentByID(t *testing.T) {

	databasePath := "./test.db"
	app, err := NewApp(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	app.router = gin.Default()
	app.router.GET("/component/:id", app.GetComponentByID)

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/1", nil)
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

func TestListComponentsForProject(t *testing.T) {
	databasePath := "./test.db"
	app, err := NewApp(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	app.router = gin.Default()
	app.router.GET("/component/filters", app.ListComponentsForProject)

	t.Run("Bad request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/component/filters?project_id=", nil)
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

	t.Run("Not Found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/component/filters?project_id=112233", nil)
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

	databasePath := "./test.db"
	app, err := NewApp(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer os.Remove(databasePath)

	app.router = gin.Default()
	app.router.DELETE("/component/:id", app.DeleteComponent)

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
