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
}

func (a *App) createNewMember(c *gin.Context) (interface{}, Response) {
	memberInput := createMemberInput{}
	if err := c.ShouldBindJSON(&memberInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("cannot read member data"))
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: memberInput.UserID}
	if err := member.Validate(); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid input data"))
	}
	err := a.DB.CheckUserAccess(member, userID.(string))
	if err == models.ErrAccessDenied {
		return nil, Forbidden(errors.New("access denied to create member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("cannot check user access to create new member"))
	}
	err = a.DB.CreateNewMember(member)
	if err == gorm.ErrDuplicatedKey {
		log.Error().Err(err).Send()
		return nil, Forbidden(errors.New("member already exists"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("cannot create new member"))
	}

	return ResponseMsg{
		Message: "member created successfully",
		Data:    member,
	}, Created()
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
		return nil, BadRequest(errors.New("cannot read member data"))
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}
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
		return nil, Forbidden(errors.New("access denied to update member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	return ResponseMsg{
		Message: "member ownership updated successfully",
	}, Ok()
}
