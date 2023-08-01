package handlers

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

	client, err := models.NewDBClient(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	err = client.Migrate()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	router := gin.Default()

	router.POST("/component", func(c *gin.Context) {
		CreateComponent(c, client)
	})

	t.Run("Success", func(t *testing.T) {

		requestBody := []byte(`{"project_id": 1, "name": "Test Component"}`)
		req, err := http.NewRequest("POST", "/component", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

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

		router.ServeHTTP(w, req)

		if http.StatusBadRequest != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetComponentByID(t *testing.T) {

	databasePath := "./test.db"

	client, err := models.NewDBClient(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	err = client.Migrate()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	router := gin.Default()

	router.GET("/component/:id", func(c *gin.Context) {
		GetComponentByID(c, client)
	})

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/component/1", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

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

		router.ServeHTTP(w, req)

		if http.StatusNotFound != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestListComponentsForProject(t *testing.T) {

	databasePath := "./test.db"

	client, err := models.NewDBClient(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	err = client.Migrate()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	router := gin.Default()

	router.GET("/component/filters", func(c *gin.Context) {
		ListComponentsForProject(c, client)
	})

	t.Run("Bad request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/component/filters?project_id=", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

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

		router.ServeHTTP(w, req)

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

		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
		}
	})
}

func TestDeleteComponent(t *testing.T) {

	databasePath := "./test.db"

	client, err := models.NewDBClient(databasePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer os.Remove(databasePath)

	err = client.Migrate()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	router := gin.Default()

	router.DELETE("/component/:id", func(c *gin.Context) {
		DeleteComponent(c, client)
	})

	t.Run("Success", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", "/component/1", nil)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

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

		router.ServeHTTP(w, req)

		if http.StatusNotFound != w.Code {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}
