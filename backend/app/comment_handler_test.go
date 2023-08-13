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

func TestCreateComment(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", WrapFunc(app.createComment))

	t.Run("create comment successfully", func(t *testing.T) {

		commentInput := CreateCommentInput{
			BugID:   10,
			Summary: "bug to be solved",
		}

		wantedComment := models.Comment{
			UserID:  "1000",
			BugID:   commentInput.BugID,
			Summary: commentInput.Summary,
		}

		payload, err := json.Marshal(commentInput)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(payload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		var response ResponseMsg
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		comment := response.Data
		commentMap := comment.(map[string]interface{})
		bugID := commentMap["bug_id"].(float64)
		userID := commentMap["user_id"]
		summary := commentMap["summary"]

		assert.Equal(t, userID, wantedComment.UserID, "got %d userID but wanted %d", userID, wantedComment.UserID)
		assert.Equal(t, uint(bugID), wantedComment.BugID, "got %d bugID but wanted %d", uint(bugID), wantedComment.BugID)
		assert.Equal(t, summary, wantedComment.Summary, "got %s summary but wanted %s", summary, wantedComment.Summary)
	})

	t.Run("failed to create comment due to incomplete input data", func(t *testing.T) {

		requestBody := []byte(`{"user_id": "1"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)

	})

	t.Run("failed to create comment as it is invalid", func(t *testing.T) {

		requestBody := []byte(`{"user_id": "1", "bug_id":-2,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)

	})
}

func TestUpdateComment(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", WrapFunc(app.createComment))

	router.PUT("/comment/:id", WrapFunc(app.updateComment))

	t.Run("update comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 5,"user_id": "3", "bug_id": 20, "summary": "bug to be solved"}`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		commentUpdate := updateCommentInput{
			Summary: "updated bug",
		}

		payload, err := json.Marshal(commentUpdate)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		request, err := http.NewRequest("PUT", "/comment/5", bytes.NewBuffer(payload))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

	t.Run("comment is not found", func(t *testing.T) {

		commentUpdate := updateCommentInput{
			Summary: "updated one",
		}

		payload, err := json.Marshal(commentUpdate)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		request, err := http.NewRequest("PUT", "/comment/50", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})

	t.Run("invalid comment", func(t *testing.T) {

		requestBody := []byte(`{"id": 200,"user_id": "12", "bug_id": 15,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"summary": 12 }`)

		request, err := http.NewRequest("PUT", "/comment/200", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)

	})

	t.Run("failed to update comment as the request is incomplete", func(t *testing.T) {

		requestBody := []byte(`{"id":300,"user_id": "13", "bug_id":15,"summary": "this is a bug" }`)
		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":300,"user_id": "13", "bug_id":,"summary": "this is a bug" }`)
		request, err := http.NewRequest("PUT", "/comment/200", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)

	})

}

func TestGetComment(t *testing.T) {

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", WrapFunc(app.createComment))

	router.GET("/:id", WrapFunc(app.getComment))

	t.Run("get comment successfully", func(t *testing.T) {

		requestBody := []byte(`{"id": 50,"user_id": "2","bug_id": 10, "summary": "bug to be solved"}`)
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

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", WrapFunc(app.createComment))

	router.DELETE("/:id", WrapFunc(app.deleteComment))

	t.Run("delete comment successfully", func(t *testing.T) {

		commentInput := CreateCommentInput{
			BugID:   10,
			Summary: "bug to be solved",
		}

		wantedComment := models.Comment{
			ID:      60,
			UserID:  "1000",
			BugID:   commentInput.BugID,
			Summary: commentInput.Summary,
		}

		payload, err := json.Marshal(wantedComment)
		if err != nil {
			t.Fatal("failed to marshal comment payload")
		}

		req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer(payload))
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

	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "testing.db")

	app, err := NewApp(dbPath)
	assert.NoError(t, err)

	router := gin.Default()

	router.POST("/comment", WrapFunc(app.createComment))

	router.GET("/filters", WrapFunc(app.listComments))

	t.Run("list comments for a specific bug", func(t *testing.T) {

		requestBody := []byte(`{"user_id": "2","bug_id": 12, "summary": "bug to be solved"}`)
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
