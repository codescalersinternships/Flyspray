package app

import (
	"strconv"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (a *App) CreateNewMember(c *gin.Context) (interface{}, Response) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
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

func (a *App) GetAllMembers(c *gin.Context) (interface{}, Response) {
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

func (a *App) UpdateMemberOwnership(c *gin.Context) (interface{}, Response) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("error binding json data"))
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := a.client.UpdateMemberOwnership(member, id)
	if err != nil {
		if err == models.ErrMemberNotFound {
			log.Error().Err(err).Send()
			return nil, NotFound(err)
		}
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("error updating member ownership"))
	}
	return ResponseMsg{
		Message: "member ownership updateds successfully",
		Data:    member,
	}, Ok()
}
