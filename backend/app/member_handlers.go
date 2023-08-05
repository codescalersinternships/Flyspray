package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

func (a *App) CreateNewMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Printf("error binding json request body to struct %q", err)
		c.Status(http.StatusBadRequest)
		return
	}
	res := a.client.CreateNewMember(member)
	if res != nil {
		log.Printf("error creating new member %q", res)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (a *App) GetAllMembers(c *gin.Context) {
	members, err := a.client.GetAllMembers()
	if err != nil {
		log.Printf("error retrieving all members %q", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, members)
}

func (a *App) UpdateMemberOwnership(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Printf("error binding json request body to struct %q", err)
		c.Status(http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := a.client.UpdateMemberOwnership(member, id)
	if err != nil {
		if err==models.ErrMemberNotFound{
			log.Printf("%q:cannot update member",err)
			c.Status(http.StatusNotFound)
			return
		}
		log.Printf("error updating member ownership %q", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, member)
}
