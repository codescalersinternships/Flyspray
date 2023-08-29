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
	Admin     bool `json:"admin"`
	ProjectID int  `json:"project_id"`
}

// createNewMember creates a new member
// @Summary Create a new member
// @Description Create a new member for a project
// @Tags Members
// @Accept json
// @Produce json
// @Param input body createMemberInput true "Member data"
// @Security ApiKeyAuth
// @Success 201 {object} ResponseMsg "member created successfully (Member details in the 'Data' field)"
// @Failure 400 {object} Response "Failed to read member data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Access denied to create member"
// @Failure 500 {object} Response "Internal server error"
// @Router /member [post]
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
	err := a.DB.CheckUserAccess(member.ProjectID, userID.(string))
	if err == models.ErrAccessDenied {
		return nil, Forbidden(errors.New("access denied to create member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	err = a.DB.CreateNewMember(member)
	if err == gorm.ErrDuplicatedKey {
		log.Error().Err(err).Send()
		return nil, Forbidden(errors.New("member already exists"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "member created successfully",
		Data:    member,
	}, Created()
}

// getMembersInProject retrieves members in a project
// @Summary Get members in a project
// @Description Get a list of members in a project
// @Tags Members
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} ResponseMsg  "members in project retrieved successfully (Members details in the 'Data' field)"
// @Failure 500 {object} Response "Internal server error"
// @Router /member/{project_id} [get]
func (a *App) getMembersInProject(c *gin.Context) (interface{}, Response) {

	project_id, err := strconv.Atoi(c.Param("project_id"))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	members, err := a.DB.GetMembersInProject(project_id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	return ResponseMsg{
		Message: "members in project retrieved successfully",
		Data:    members,
	}, Ok()
}

// updateMemberOwnership updates the ownership of a member
// @Summary Update member ownership
// @Description Update the ownership of a member
// @Tags Members
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param input body updateMemberInput true "Member data"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "Member ownership updated successfully"
// @Failure 400 {object} Response "Failed to read member data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Access denied to update member"
// @Failure 404 {object} Response "Member is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /member/{id} [put]
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
	err = a.DB.CheckUserAccess(memberInput.ProjectID, userID.(string))
	if err == models.ErrAccessDenied {
		return nil, Forbidden(errors.New("access denied to update member"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	err = a.DB.UpdateMemberOwnership(id, memberInput.Admin, memberInput.ProjectID)
	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(err)
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	return ResponseMsg{
		Message: "member ownership updated successfully",
	}, Ok()
}
