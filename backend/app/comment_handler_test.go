package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {

	dbPath := "./testing.db"

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		app.CreateComment(c)
	})

	t.Run("create comment successfully", func(t *testing.T) {

		wantedComment := models.Comment{
			OwnerID: 1,
			BugID:   10,
			Summary: "bug to be solved",
		}

		payload, err := json.Marshal(wantedComment)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		var response Response
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		jsonData, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println("response body:", string(jsonData))

		fmt.Println("response", response.Payload[0].Summary)

		assert.Equal(t, response.Payload[0].OwnerID, wantedComment.OwnerID, "got %d ownerID but wanted %d", response.Payload[0].OwnerID, wantedComment.OwnerID)
		assert.Equal(t, response.Payload[0].BugID, wantedComment.BugID, "got %d bugID but wanted %d", response.Payload[0].BugID, wantedComment.BugID)
		assert.Equal(t, response.Payload[0].Summary, wantedComment.Summary, "got %s summary but wanted %s", response.Payload[0].Summary, wantedComment.Summary)
	})

	t.Run("failed to create comment due to incomplete input data", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 1}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)

	})

	t.Run("failed to create comment as it is invalid", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 1, "bug_id":-2,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)

	})
}

func TestUpdateComment(t *testing.T) {
	dbPath := "./testing.db"

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		app.CreateComment(c)
	})

	router.PUT("/comment/:id", func(c *gin.Context) {
		app.UpdateComment(c)
	})

	t.Run("update comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id":100,"owner_id": 3, "bug_id": 10, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":100,"owner_id": 3, "bug_id": 10, "summary": "updated bug"}`)

		request, err := http.NewRequest("PUT", "/comment/100", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

		var response Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		jsonData, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println("response body:", string(jsonData))

		wantedComment := models.Comment{
			ID:      100,
			OwnerID: 3,
			BugID:   10,
			Summary: "updated bug",
		}

		assert.Equal(t, response.Payload[0].ID, wantedComment.ID, "got %d ID but wanted %d", response.Payload[0].ID, wantedComment.ID)
		assert.Equal(t, response.Payload[0].OwnerID, wantedComment.OwnerID, "got %d ownerID but wanted %d", response.Payload[0].OwnerID, wantedComment.OwnerID)
		assert.Equal(t, response.Payload[0].BugID, wantedComment.BugID, "got %d bugID but wanted %d", response.Payload[0].BugID, wantedComment.BugID)
		assert.Equal(t, response.Payload[0].Summary, wantedComment.Summary, "got %s summary but wanted %s", response.Payload[0].Summary, wantedComment.Summary)
	})

	t.Run("comment is not found", func(t *testing.T) {

		updatedRequestBody := []byte(`{"id":100,"owner_id": 3, "bug_id": 10, "summary": "updated bug"}`)

		request, err := http.NewRequest("PUT", "/comment/55", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})

	t.Run("failed to update comment as it is invalid", func(t *testing.T) {

		requestBody := []byte(`{"id":200,"owner_id": 12, "bug_id":15,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":200,"owner_id": 12, "bug_id":15,"summary": 12 }`)

		request, err := http.NewRequest("PUT", "/comment/200", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)

	})

	t.Run("failed to update comment as the request is incomplete (bad request)", func(t *testing.T) {

		requestBody := []byte(`{"id":300,"owner_id": 13, "bug_id":15,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":300,"owner_id": 13, "bug_id":,"summary": "this is a bug" }`)
		request, err := http.NewRequest("PUT", "/comment/200", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)

	})

}

func TestGetComment(t *testing.T) {

	dbPath := "./testing.db"

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		app.CreateComment(c)
	})

	router.GET("/comment/:id", func(c *gin.Context) {
		app.GetComment(c)
	})

	t.Run("get comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 50,"owner_id": 2,"bug_id": 10, "summary": "bug to be solved"}`)
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

	t.Run("comment ID is not given", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/comment/ ", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)

	})

}

func TestDeleteComment(t *testing.T) {

	dbPath := "./testing.db"

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		app.CreateComment(c)
	})

	router.DELETE("/comment/:id", func(c *gin.Context) {
		app.DeleteComment(c)
	})

	t.Run("delete comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 60,"owner_id": 5,"bug_id": 10, "summary": "bug to be solved"}`)
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

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)

	})
	t.Run("comment ID is not given", func(t *testing.T) {

		request, err := http.NewRequest("DELETE", "/comment/ ", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)

	})

}

func TestListComments(t *testing.T) {

	dbPath := "./testing.db"

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", func(c *gin.Context) {
		app.CreateComment(c)
	})

	router.GET("/comment/filters", func(c *gin.Context) {
		app.ListComments(c)
	})

	t.Run("list comments for a specific bug", func(t *testing.T) {

		requestBody := []byte(`{"owner_id": 2,"bug_id": 12, "summary": "bug to be solved"}`)
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

	t.Run("missing the bug_id in the get request (bad request) ", func(t *testing.T) {

		request, err := http.NewRequest("GET", "/comment/filters?bug_id=", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)

	})

}
