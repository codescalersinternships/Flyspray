package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {

	dbPath := "./testing.db"

	client, err := models.NewDBClient(dbPath)
	assert.NoError(t, err)

	err = client.Migrate()
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		CreateComment(c, client)
	})

	t.Run("create comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 1,"bug_id": 10, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		createdComment := models.Comment{}
		err = json.Unmarshal(recorder.Body.Bytes(), &createdComment)
		assert.NoError(t, err)

		wantedComment := models.Comment{
			OwnerID: 1,
			BugID:   10,
			Summary: "bug to be solved",
		}

		assert.Equal(t, createdComment.OwnerID, wantedComment.OwnerID, "got %d ownerID but wanted %d")
		assert.Equal(t, createdComment.BugID, wantedComment.BugID, "got %d bugID but wanted %d")
		assert.Equal(t, createdComment.Summary, wantedComment.Summary, "got %d summary but wanted %d")

	})

	t.Run("failed to create comment due to bad request", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 1"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 404", recorder.Code)

	})

}

func TestGetComment(t *testing.T) {

	dbPath := "./testing.db"

	client, err := models.NewDBClient(dbPath)
	assert.NoError(t, err)

	err = client.Migrate()
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		CreateComment(c, client)
	})

	router.GET("/comment/:id", func(c *gin.Context) {
		GetComment(c, client)
	})

	t.Run("get comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 50,"owner_id": 1,"bug_id": 10, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/comment/50", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

	t.Run("comment is not found", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/comment/55", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)

	})

}

func TestDeleteComment(t *testing.T) {

	dbPath := "./testing.db"

	client, err := models.NewDBClient(dbPath)
	assert.NoError(t, err)

	err = client.Migrate()
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		CreateComment(c, client)
	})

	router.DELETE("/comment/:id", func(c *gin.Context) {
		DeleteComment(c, client)
	})

	t.Run("delete comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 60,"owner_id": 1,"bug_id": 10, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("DELETE", "/comment/60", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

	t.Run("comment is not found", func(t *testing.T) {

		request, err := http.NewRequest("DELETE", "/comment/55", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

}

func TestListComments(t *testing.T) {

	dbPath := "./testing.db"

	client, err := models.NewDBClient(dbPath)
	assert.NoError(t, err)

	err = client.Migrate()
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		CreateComment(c, client)
	})

	router.POST("/comment/filters", func(c *gin.Context) {
		ListComments(c, client)
	})

	t.Run("list comments for a specific bug", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 1,"bug_id": 12, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/comment/filters?bug_id=12", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

	t.Run("no comments found for the bug", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/comment/filters?bug_id=6", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)

	})

}

