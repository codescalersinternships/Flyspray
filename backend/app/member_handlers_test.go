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
)

func TestCreateNewMember(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.NoError(t, err, "failed to connect to database")
	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.setRoutes()
	t.Run("creating valid member returns status 201 and correct response", func(t *testing.T) {
		member := models.Member{UserID: "1", ProjectID: 2}
		expectedResponse := ResponseMsg{
			Message: "member created successfully",
			Data:    models.Member{ID: 1, UserID: "1", ProjectID: 2, Admin: false, Project: nil},
		}
		jsonData, err := json.Marshal(member)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		got := strings.ReplaceAll(resp.Body.String(), " ", "")
		got = strings.ReplaceAll(got, "\n", "")
		wantJson, err := json.Marshal(expectedResponse)
		want := string(wantJson)
		want = strings.ReplaceAll(want, " ", "")
		assert.NoError(t, err, "failed to marshal json data")
		assert.Equal(t, want, got)
		if resp.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, resp.Code)
		}
	})
	t.Run("creating invalid member returns status 400", func(t *testing.T) {
		member := models.Member{}
		jsonData, err := json.Marshal(member)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, resp.Code)
		}
	})
	t.Run("creating a member that already exists returns status 403", func(t *testing.T) {
		member := models.Member{UserID: "1", ProjectID: 2}
		jsonData, err := json.Marshal(member)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusForbidden {
			t.Errorf("expected status code %d but got %d", http.StatusForbidden, resp.Code)
		}
	})

}

func TestGetAllMembers(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	assert.NoError(t, err, "failed to connect to database")
	member := models.Member{UserID: "1", ProjectID: 2}
	jsonData, err := json.Marshal(member)
	assert.NoError(t, err, "failed to marshal json data")
	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.setRoutes()
	req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
	resp := httptest.NewRecorder()
	app.router.ServeHTTP(resp, req)
	assert.NoError(t, err, "failed to create http request")
	t.Run("getallmembers returns status 200", func(t *testing.T) {
		expectedResponse := ResponseMsg{
			Message: "all members retrieved successfully",
			Data:    []models.Member{{ID: 1, UserID: "1", ProjectID: 2, Admin: false, Project: nil}},
		}
		req, err := http.NewRequest("GET", "/member", nil)
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		got := strings.ReplaceAll(resp.Body.String(), " ", "")
		got = strings.ReplaceAll(got, "\n", "")
		wantJson, err := json.Marshal(expectedResponse)
		want := string(wantJson)
		want = strings.ReplaceAll(want, " ", "")
		assert.NoError(t, err, "failed to marshal json data")
		assert.Equal(t, want, got)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, resp.Code)
		}
	})
}

func TestUpdateMemberOwnership(t *testing.T) {
	dir := t.TempDir()
	app, err := NewApp(filepath.Join(dir, "test.db"))
	app.router = gin.Default()
	app.router.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	app.setRoutes()
	assert.NoError(t, err, "failed to connect to database")
	member := models.Member{UserID: "1", ProjectID: 2}
	createdData, err := json.Marshal(member)
	assert.NoError(t, err, "failed to marshal json data")
	req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(createdData))
	assert.NoError(t, err, "failed to create http request")
	resp := httptest.NewRecorder()
	app.router.ServeHTTP(resp, req)
	jsonData, err := json.Marshal(models.Member{Admin: true})
	assert.NoError(t, err, "failed to marshal json data")
	t.Run("updating member successfully returns status 200 and correct response", func(t *testing.T) {
		req, err = http.NewRequest("PUT", "/member/1", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp = httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		got := strings.ReplaceAll(resp.Body.String(), " ", "")
		got = strings.ReplaceAll(got, "\n", "")
		expectedResponse := ResponseMsg{
			Message: "member ownership updated successfully",
			Data:    models.Member{ID: 1, UserID: "1", ProjectID: 2, Admin: true, Project: nil},
		}
		wantJson, err := json.Marshal(expectedResponse)
		want := string(wantJson)
		want = strings.ReplaceAll(want, " ", "")
		assert.NoError(t, err, "failed to marshal json data")
		assert.Equal(t, want, got)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, resp.Code)
		}
	})
	t.Run("updating member fails with status 404 with invalid id", func(t *testing.T) {
		req, err = http.NewRequest("PUT", "/member/2", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp = httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusNotFound {
			t.Errorf("expected status code %d but got %d", http.StatusNotFound, resp.Code)
		}
	})

}
