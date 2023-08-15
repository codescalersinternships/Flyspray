package app

import (
	"strconv"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type createMemberInput struct {
	ProjectID int  `json:"project_id" validate:"required"`
	Admin     bool `json:"admin_bool"`
}

func (a *App) createNewMember(c *gin.Context) (interface{}, Response) {
	memberInput := createMemberInput{}

	if err := c.ShouldBindJSON(&memberInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
	// UserID will be replaced by that of middleware
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: 100}
	err := a.client.CreateNewMember(member)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error creating new member"))
	}
	return ResponseMsg{
		Message: "member created successfully",
		Data:    member,
	}, Created()
}

func (a *App) getAllMembers(c *gin.Context) (interface{}, Response) {
	members, err := a.client.GetAllMembers()
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error getting all members"))
	}
	return ResponseMsg{
		Message: "all members retrieved successfully",
		Data:    members,
	}, Ok()
}

func (a *App) updateMemberOwnership(c *gin.Context) (interface{}, Response) {

	memberInput := createMemberInput{}

	if err := c.ShouldBindJSON(&memberInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error parsing id"))
	}
	// UserID will be replaced by that of middleware
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: 100}
	err = a.client.UpdateMemberOwnership(member, id)
	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(err)
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error updating member ownership"))
	}
	return ResponseMsg{
		Message: "member ownership updated successfully",
		Data:    member,
	}, Ok()
}
