package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewMember(t *testing.T) {
	app, err := NewApp("./test.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	testCases := []struct {
		name           string
		member         models.Member
		expectedStatus int
	}{
		{
			name:           "creating valid member returns status 201",
			member:         models.Member{UserID: 1, ProjectID: 2},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "creating invalid member returns status 500",
			member:         models.Member{},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tC.member)
			assert.NoError(t, err, "failed to marshal json data")
			req, err := http.NewRequest("POST", "/member", bytes.NewBuffer(jsonData))
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			app.router.ServeHTTP(resp, req)
			if resp.Code != tC.expectedStatus {
				t.Errorf("expected status code %d but got %d", tC.expectedStatus, resp.Code)
			}
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	app, err := NewApp("./test.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	t.Run("getallmembers returns status 200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/member", nil)
		assert.NoError(t, err, "failed to create http request")
		resp := httptest.NewRecorder()
		app.router.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, resp.Code)
		}
	})
}

func TestUpdateMemberOwnership(t *testing.T) {
	app, err := NewApp("./test.db")
	assert.NoError(t, err, "failed to connect to database")
	app.setRoutes()
	testCases := []struct {
		name      string
		id        int
		errStatus int
	}{
		{
			name:      "updatememberownership returns status 200 when id is valid",
			id:        1,
			errStatus: http.StatusOK,
		},
		{
			name:      "updatememberownership returns status 404 when id is invalid",
			id:        -1,
			errStatus: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			jsonData, err := json.Marshal(models.Member{Admin: false})
			assert.NoError(t, err, "failed to marshal json data")
			req, err := http.NewRequest("PUT", fmt.Sprintf("/member/%d", tC.id), bytes.NewBuffer(jsonData))
			assert.NoError(t, err, "failed to create http request")
			resp := httptest.NewRecorder()
			app.router.ServeHTTP(resp, req)
			if resp.Code != tC.errStatus {
				t.Errorf("expected status code %d but got %d", tC.errStatus, resp.Code)
			}
		})
	}
}
