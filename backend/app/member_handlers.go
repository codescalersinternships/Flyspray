package app

import (
	"fmt"
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
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: userID.(string)}
	if err := member.Validate(); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	if fmt.Sprint(userID) != fmt.Sprint(member.UserID) && !member.Admin {
		return nil, Forbidden(errors.New("you have no access to create member"))
	}

	addedMember, err := a.DB.CreateNewMember(member)
	if err == gorm.ErrDuplicatedKey{
		log.Error().Err(err).Send()
		return nil, Forbidden(errors.New("member already exists"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error creating new member"))
	}

	return ResponseMsg{
		Message: "member created successfully",
		Data:    addedMember,
	}, Created()
}

func (a *App) getAllMembers(c *gin.Context) (interface{}, Response) {
	members, err := a.DB.GetAllMembers()
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
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error parsing id"))
	}
	member := models.Member{ProjectID: memberInput.ProjectID, Admin: memberInput.Admin, UserID: userID.(string)}
	if fmt.Sprint(userID) != fmt.Sprint(member.UserID) && !member.Admin {
		return nil, Forbidden(errors.New("you have no access to update member"))
	}
	updatedMember, err := a.DB.UpdateMemberOwnership(member, id)
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
		Data:    updatedMember,
	}, Ok()
}
