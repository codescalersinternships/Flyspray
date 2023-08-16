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
	ProjectID int    `json:"project_id" validate:"required"`
	Admin     bool   `json:"admin"`
	UserID    string `json:"user_id" validate:"required"`
}
type updateMemberInput struct {
	Admin bool `json:"admin"`
	ID    int
}

func (a *App) createNewMember(c *gin.Context) (interface{}, Response) {
	memberInput := createMemberInput{}
	if err := c.ShouldBindJSON(&memberInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error reading member data"))
	}
	userID, _ := c.Get("user_id")
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: memberInput.UserID}
	if err := member.Validate(); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error invalid input data"))
	}
	err := a.DB.CreateNewMember(member, userID.(string))
	if err == gorm.ErrDuplicatedKey {
		log.Error().Err(err).Send()
		return nil, Forbidden(errors.New("error member already exists"))
	}
	if err == models.ErrAccessDenied {
		return nil, Forbidden(errors.New("error access denied to create member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error cannot create new member"))
	}

	return ResponseMsg{
		Message: "member created successfully",
		Data:    member,
	}, Created()
}

func (a *App) getAllMembers(c *gin.Context) (interface{}, Response) {
	members, err := a.DB.GetAllMembers()
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error retrieving all members"))
	}

	return ResponseMsg{
		Message: "all members retrieved successfully",
		Data:    members,
	}, Ok()
}
func (a *App) getMembersInProject(c *gin.Context) (interface{}, Response) {

	project_id, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	members, err := a.DB.GetMembersInProject(project_id)
	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(err)
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	return ResponseMsg{
		Message: "members in project retrieved successfully",
		Data:    members,
	}, Ok()
}
func (a *App) updateMemberOwnership(c *gin.Context) (interface{}, Response) {
	memberInput := updateMemberInput{}
	if err := c.ShouldBindJSON(&memberInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error reading member data"))
	}
	userID, _ := c.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	err = a.DB.UpdateMemberOwnership(id, memberInput.Admin, userID.(string))
	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(err)
	}
	if err == models.ErrAccessDenied {
		return nil, Forbidden(errors.New("error access denied to update member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	return ResponseMsg{
		Message: "member ownership updated successfully",
	}, Ok()
}
