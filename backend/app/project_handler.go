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

type responseErr struct {
	Message string `json:"message"`
}

type responseOk struct {
	Message string           `json:"message"`
	Data    []models.Project `json:"data"`
}

func (a *App) createProject(ctx *gin.Context) {
	var input createProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Println("project name must be specified")
		ctx.IndentedJSON(http.StatusBadRequest,
			responseErr{Message: "name must be specified"},
		)
		return
	}

	// TODO: get user id from authorization middleware and assign it to OwnerId
	newProject := models.Project{Name: input.Name, OwnerId: 10007} // 10007 is just a random number
	newProject, err := models.CreateProject(newProject, a.client)

	if err != nil {
		log.Fatal(err)
		ctx.IndentedJSON(http.StatusInternalServerError,
			responseErr{Message: "could not create new project"},
		)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, responseOk{
		Message: "project created successfully",
		Data:    []models.Project{newProject},
	})
}

func (a *App) updateProject(ctx *gin.Context) {
	// TODO: get user id from authorization middleware and check if user has access to update the project
	id := ctx.Param("id")
	var input createProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Println("project name must be specified")
		ctx.IndentedJSON(http.StatusBadRequest,
			responseErr{Message: "name must be specified"},
		)
		return
	}

	updatedProject := models.Project{Name: input.Name}
	updatedProject, err := models.UpdateProject(id, updatedProject, a.client)

	if err == gorm.ErrRecordNotFound {
		log.Println("project not found")
		ctx.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "project not found"},
		)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, responseOk{
		Message: "project updated successfully",
		Data:    []models.Project{updatedProject},
	})
}

func (a *App) getProject(ctx *gin.Context) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	project, err := models.GetProject(id, a.client)
	if err == gorm.ErrRecordNotFound {
		log.Println("project not found")
		ctx.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "project not found"},
		)
		return
	}

	ctx.IndentedJSON(http.StatusOK, responseOk{
		Message: "project retrieved successfully",
		Data:    []models.Project{project},
	})
}

func (a *App) getProjects(ctx *gin.Context) {
	// TODO: add middleware to check if user is signed in
	userId := ctx.Query("userid")
	projectName := ctx.Query("name")
	creationDate := ctx.Query("after")

	projects := models.FilterProjects(userId, projectName, creationDate, a.client)

	ctx.IndentedJSON(http.StatusOK, responseOk{
		Message: "projects retrieved successfully",
		Data:    projects,
	})
}

func (a *App) deleteProject(ctx *gin.Context) {
	// TODO: get user id from authorization middleware and check if user has access to delete the project
	id := ctx.Param("id")

	models.DeleteProject(id, a.client)

	ctx.IndentedJSON(http.StatusOK, responseOk{
		Message: "project deleted successfully",
	})
}
