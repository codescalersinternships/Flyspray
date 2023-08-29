package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	//gin.SetMode(gin.TestMode)

	app := App{}
	var err error
	app.DB, err = models.NewDBClient(path)
	assert.NoError(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	router := gin.Default()
	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})
	router.POST("/bug", WrapFunc(app.createBug))

	t.Run("create new bug successfully", func(t *testing.T) {
		bugInput := createBugInput{
			Summary:     "hello world!",
			ComponentID: 13,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		payload, err := json.Marshal(bugInput)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		wantedBug := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput.Summary,
			ComponentID: bugInput.ComponentID,
			Category:    bugInput.Category,
			Severity:    bugInput.Severity,
			Status:      bugInput.Status,
		}

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		var response ResponseMsg
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.NoError(t, err)

		bug := response.Data
		bugMap := bug.(map[string]interface{})
		userID := bugMap["user_id"]
		componentID := int(bugMap["component_id"].(float64))
		summary := bugMap["summary"]
		category := bugMap["category"]
		severity := bugMap["severity"]
		status := bugMap["status"]

		assert.Equal(t, userID, wantedBug.UserID, "got %d user id but wanted %d", userID, wantedBug.UserID)
		assert.Equal(t, componentID, wantedBug.ComponentID, "got %d component but wanted %d", componentID, wantedBug.ComponentID)
		assert.Equal(t, summary, wantedBug.Summary, "got %d as a summary but wanted %d", summary, wantedBug.Summary)
		assert.Equal(t, category, wantedBug.Category, "got %d category but wanted %d", category, wantedBug.Category)
		assert.Equal(t, severity, wantedBug.Severity, "got %d severity but wanted %d", severity, wantedBug.Severity)
		assert.Equal(t, status, wantedBug.Status, "got %d status but wanted %d", status, wantedBug.Status)
	})

	t.Run("failed to create new bug", func(t *testing.T) {
		requestBody := []byte(`{"user_id": "99"}`)

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "got %d status code but want status code 400", recorder.Code)
	})
}

func TestGetbugs(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app := App{}
	var err error
	app.DB, err = models.NewDBClient(path)
	assert.NoError(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	router := gin.Default()
	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.POST("/bug", WrapFunc(app.createBug))

	router.GET("/bug/filters", WrapFunc(app.getbugs))

	t.Run("get all bugs successfully", func(t *testing.T) {

		bugInput1 := createBugInput{
			Summary:     "first",
			ComponentID: 90,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		wantedBug1 := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput1.Summary,
			ComponentID: bugInput1.ComponentID,
			Category:    bugInput1.Category,
			Severity:    bugInput1.Severity,
			Status:      bugInput1.Status,
		}

		result := app.DB.Client.Create(&wantedBug1)
		assert.NoError(t, result.Error)

		bugInput2 := createBugInput{
			Summary:     "second bug",
			ComponentID: 100,
			Category:    "parts",
			Severity:    "high",
			Status:      "nearly solved",
		}

		wantedBug2 := models.Bug{
			ID:          11,
			UserID:      "99",
			Summary:     bugInput2.Summary,
			ComponentID: bugInput2.ComponentID,
			Category:    bugInput2.Category,
			Severity:    bugInput2.Severity,
			Status:      bugInput2.Status,
		}

		result = app.DB.Client.Create(&wantedBug2)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("GET", "/bug/filters", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("get all bugs for a specific component", func(t *testing.T) {

		bugInput1 := createBugInput{
			Summary:     "first",
			ComponentID: 90,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		wantedBug1 := models.Bug{
			ID:          12,
			UserID:      "99",
			Summary:     bugInput1.Summary,
			ComponentID: bugInput1.ComponentID,
			Category:    bugInput1.Category,
			Severity:    bugInput1.Severity,
			Status:      bugInput1.Status,
		}

		result := app.DB.Client.Create(&wantedBug1)
		assert.NoError(t, result.Error)

		bugInput2 := createBugInput{
			Summary:     "second bug",
			ComponentID: 90,
			Category:    "parts",
			Severity:    "high",
			Status:      "In progress",
		}

		wantedBug2 := models.Bug{
			ID:          13,
			UserID:      "99",
			Summary:     bugInput2.Summary,
			ComponentID: bugInput2.ComponentID,
			Category:    bugInput2.Category,
			Severity:    bugInput2.Severity,
			Status:      bugInput2.Status,
		}

		result = app.DB.Client.Create(&wantedBug2)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("GET", "/bug/filters?component_id=90", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("get all bugs for a specific category", func(t *testing.T) {

		bugInput1 := createBugInput{
			Summary:     "first",
			ComponentID: 100,
			Category:    "tests",
			Severity:    "low",
			Status:      "In progress",
		}

		wantedBug1 := models.Bug{
			ID:          14,
			UserID:      "99",
			Summary:     bugInput1.Summary,
			ComponentID: bugInput1.ComponentID,
			Category:    bugInput1.Category,
			Severity:    bugInput1.Severity,
			Status:      bugInput1.Status,
		}

		result := app.DB.Client.Create(&wantedBug1)
		assert.NoError(t, result.Error)

		bugInput2 := createBugInput{
			Summary:     "second bug",
			ComponentID: 90,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		wantedBug2 := models.Bug{
			ID:          15,
			UserID:      "99",
			Summary:     bugInput2.Summary,
			ComponentID: bugInput2.ComponentID,
			Category:    bugInput2.Category,
			Severity:    bugInput2.Severity,
			Status:      bugInput2.Status,
		}

		result = app.DB.Client.Create(&wantedBug2)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("GET", "/bug/filters?category=tests", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("get all bugs for a specific status", func(t *testing.T) {

		bugInput1 := createBugInput{
			Summary:     "first",
			ComponentID: 100,
			Category:    "tests",
			Severity:    "low",
			Status:      "In progress",
		}

		wantedBug1 := models.Bug{
			ID:          16,
			UserID:      "99",
			Summary:     bugInput1.Summary,
			ComponentID: bugInput1.ComponentID,
			Category:    bugInput1.Category,
			Severity:    bugInput1.Severity,
			Status:      bugInput1.Status,
		}

		result := app.DB.Client.Create(&wantedBug1)
		assert.NoError(t, result.Error)

		bugInput2 := createBugInput{
			Summary:     "second bug",
			ComponentID: 90,
			Category:    "parts",
			Severity:    "high",
			Status:      "In progress",
		}

		wantedBug2 := models.Bug{
			ID:          17,
			UserID:      "99",
			Summary:     bugInput2.Summary,
			ComponentID: bugInput2.ComponentID,
			Category:    bugInput2.Category,
			Severity:    bugInput2.Severity,
			Status:      bugInput2.Status,
		}

		result = app.DB.Client.Create(&wantedBug2)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("GET", "/bug/filters?status=In progress", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})
}

func TestGetBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app := App{}
	var err error
	app.DB, err = models.NewDBClient(path)
	assert.NoError(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.POST("/bug", WrapFunc(app.createBug))

	router.GET("/bug/:id", WrapFunc(app.getBug))

	t.Run("get bug successfully", func(t *testing.T) {
		projectInput := createProjectInput{
			Name: "flyspray",
		}

		project := models.Project{
			OwnerID: "99",
			Name:    projectInput.Name,
		}

		result := app.DB.Client.Create(&project)
		assert.NoError(t, result.Error)

		projectId := fmt.Sprintf("%x", project.ID)

		componentInput := createComponentInput{
			Name:      "backend",
			ProjectID: projectId,
		}

		component := models.Component{
			ID:        13,
			UserID:    "99",
			Name:      componentInput.Name,
			ProjectID: componentInput.ProjectID,
		}

		result = app.DB.Client.Create(&component)
		assert.NoError(t, result.Error)
		bugData := models.Bug{
			UserID:      "99",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug/1", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("get bug is not found", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug/5", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug ID is not given", func(t *testing.T) {
		bugData := models.Bug{
			UserID:      "12",
			ComponentID: 13,
		}

		payload, err := json.Marshal(bugData)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		req, err := http.NewRequest(http.MethodPost, "/bug", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		request, err := http.NewRequest("GET", "/bug", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})
}

func TestUpdateBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app := App{}
	var err error
	app.DB, err = models.NewDBClient(path)
	assert.NoError(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.POST("/bug", WrapFunc(app.createBug))

	router.PUT("/bug/:id", WrapFunc(app.updateBug))

	t.Run("update bug successfully", func(t *testing.T) {
		bugInput := createBugInput{
			Summary:     "hello from test",
			ComponentID: 90,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		wantedBug := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput.Summary,
			ComponentID: bugInput.ComponentID,
			Category:    bugInput.Category,
			Severity:    bugInput.Severity,
			Status:      bugInput.Status,
		}

		result := app.DB.Client.Create(&wantedBug)
		assert.NoError(t, result.Error)

		bugUpdate := updateBugInput{
			Summary:  "update bug",
			Category: "tests",
			Severity: "low",
			Status:   "not solved yet",
		}

		payload, err := json.Marshal(bugUpdate)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		request, err := http.NewRequest("PUT", "/bug/10", bytes.NewBuffer(payload))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)
		assert.Equal(t, http.StatusOK, rec.Code, "got %d status code but want status code 200", rec.Code)

	})

	t.Run("bug is not found", func(t *testing.T) {

		bugUpdate := updateBugInput{
			Summary:  "update bug",
			Category: "tests",
			Severity: "low",
			Status:   "not solved yet",
		}

		payload, err := json.Marshal(bugUpdate)
		if err != nil {
			t.Fatal("failed to marshal bug payload")
		}

		request, err := http.NewRequest("PUT", "/bug/50", bytes.NewBuffer(payload))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})

	t.Run("failed to update bug as it is invalid", func(t *testing.T) {
		requestBody := []byte(`{"id":50,"user_id": "99", "component_id":15}`)
		request, err := http.NewRequest("POST", "/bug", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code, "got %d status code but want status code 201", recorder.Code)

		updatedRequestBody := []byte(`{"id":50,"user_id": "99", "component_id" }`)

		req, err := http.NewRequest("PUT", "/bug/33", bytes.NewBuffer(updatedRequestBody))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "got %d status code but want status code 400", rec.Code)
	})
}

func TestDeleteBug(t *testing.T) {
	tempDir := t.TempDir()
	path := path.Join(tempDir, "flyspray.db")
	gin.SetMode(gin.TestMode)

	// Create a new instance of your App
	app := App{}
	var err error
	app.DB, err = models.NewDBClient(path)
	assert.NoError(t, err)
	err = app.DB.Migrate()
	assert.Nil(t, err)

	router := gin.Default()

	// middleware to set the "user_id" in the gin context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "99")
		c.Next()
	})

	router.DELETE("/bug/:id", WrapFunc(app.deleteBug))

	t.Run("delete bug successfully", func(t *testing.T) {

		projectInput := createProjectInput{
			Name: "flyspray",
		}

		project := models.Project{
			OwnerID: "99",
			Name:    projectInput.Name,
		}

		result := app.DB.Client.Create(&project)
		assert.NoError(t, result.Error)

		projectId := fmt.Sprintf("%x", project.ID)

		componentInput := createComponentInput{
			Name:      "backend",
			ProjectID: projectId,
		}

		component := models.Component{
			ID:        90,
			UserID:    "99",
			Name:      componentInput.Name,
			ProjectID: componentInput.ProjectID,
		}

		result = app.DB.Client.Create(&component)
		assert.NoError(t, result.Error)

		bugInput := createBugInput{
			Summary:     "hello from test",
			ComponentID: 90,
			Category:    "tests",
			Severity:    "low",
			Status:      "not solved yet",
		}

		wantedBug := models.Bug{
			ID:          10,
			UserID:      "99",
			Summary:     bugInput.Summary,
			ComponentID: bugInput.ComponentID,
			Category:    bugInput.Category,
			Severity:    bugInput.Severity,
			Status:      bugInput.Status,
		}

		result = app.DB.Client.Create(&wantedBug)
		assert.NoError(t, result.Error)

		request, err := http.NewRequest("DELETE", "/bug/10", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code, "got %d status code but want status code 201", recorder.Code)
	})

	t.Run("bug is not found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/bug/50", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 200", rec.Code)
	})

	t.Run("bug ID is not given", func(t *testing.T) {
		request, err := http.NewRequest("DELETE", "/bug/", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, request)

		assert.Equal(t, http.StatusNotFound, rec.Code, "got %d status code but want status code 404", rec.Code)
	})
}
