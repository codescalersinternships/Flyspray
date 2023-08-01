package models

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed sql/01createMemberTable.sql
var createMemberTableQuery string

//go:embed sql/01getAllMembers.sql
var getAllMembersQuery string

//go:embed sql/01createNewMember.sql
var createNewMemberQuery string

//go:embed sql/01updateMemberTable.sql
var updateMemberQuery string

// member struct has Member table's columns
type member struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	ProjectID int  `json:"project_id"`
	Admin     bool `json:"admin_bool"`
}

// CreateMemberTable creates member table if it does not exist
func (db *DBClient) CreateMemberTable() error {
	res := db.Client.Exec(createMemberTableQuery)
	return res.Error
}

// Crete adds a new member to member table
func (db *DBClient) Create(c *gin.Context) {
	var member member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	res := db.Client.Exec(createNewMemberQuery, member.UserID, member.ProjectID, member.Admin)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, member)
}

// GetAllMembers returns all members in member table
func (db *DBClient) GetAllMembers(c *gin.Context) {
	var members []member
	rows, err := db.Client.Raw(getAllMembersQuery).Rows()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var member member
		err := rows.Scan(&member.ID, &member.Admin, &member.UserID, &member.ProjectID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		members = append(members, member)
	}
	c.JSON(http.StatusOK, members)
}

// UpdateMemberOwnership updates the admin bool in member table
func (db *DBClient) UpdateMemberOwnership(c *gin.Context) {
	var member member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	res := db.Client.Exec(updateMemberQuery, boolToInt(member.Admin), c.Param("id"))
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, member)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
