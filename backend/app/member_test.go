package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	t.Run("assert new app returns no error", func(t *testing.T) {
		_, err := NewApp("../flyspray.db")
		assert.NoError(t, err)
	})
}

func TestCreateMemberTable(t *testing.T) {
	t.Run("member table is created successfully", func(t *testing.T) {
		_, err := NewApp("../flyspray.db")
		assert.NoError(t, err)
	})
}

func TestCreateNewMember(t *testing.T) {
	app, err := NewApp("./test.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	t.Run("creating valid member returns status 201", func(t *testing.T) {
		member := models.Member{UserID: 1, ProjectID: 2}
		jsonData, err := json.Marshal(member)
		assert.NoError(t, err, "failed to marshal json data")
		req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, resp.Code)
		}
	})
}
