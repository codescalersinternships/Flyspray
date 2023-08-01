package app

import (
	"log"
	"net/http"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createProjectInput struct {
	Name string `json:"name" binding:"required"`
}

func createProject(ctx *gin.Context, db models.DBClient) {
	var input createProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "name must be specified"})
		return
	}

	// TODO: get user id from authorization middleware and assign it to OwnerId
	newProject := models.Project{Name: input.Name, OwnerId: 10007} // 10007 is just a random number

	if result := db.Client.Create(&newProject); result.Error != nil {
		log.Fatal(result.Error)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newProject)
}

func getProject(ctx *gin.Context, db models.DBClient) {
	id := ctx.Param("id")

	project := models.Project{}
	result := db.Client.First(&project, id)

	if result.Error == gorm.ErrRecordNotFound {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, project)
}

func getProjects(ctx *gin.Context, db models.DBClient) {
	userId := ctx.Query("userid")

	projects := []models.Project{}
	db.Client.Where("owner_id = ?", userId).Find(&projects)

	ctx.IndentedJSON(http.StatusOK, projects)
}

func deleteProject(ctx *gin.Context, db models.DBClient) {
	// TODO: get user id from authorization middleware and check if user has access to delete the project
	id := ctx.Param("id")

	project := models.Project{}
	db.Client.Delete(&project, id)

	ctx.Status(http.StatusOK)
}
