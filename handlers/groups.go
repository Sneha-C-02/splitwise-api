package handlers

import (
	"net/http"
	"splitwise-api/config"
	"splitwise-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateGroup — POST /groups
func CreateGroup(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		CreatedBy uint   `json:"created_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify creator exists
	var creator models.User
	if err := config.DB.First(&creator, input.CreatedBy).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Creator user not found"})
		return
	}

	group := models.Group{
		Name:      input.Name,
		CreatedBy: input.CreatedBy,
	}
	config.DB.Create(&group)

	// Auto-add creator as a member
	member := models.GroupMember{GroupID: group.ID, UserID: input.CreatedBy}
	config.DB.Create(&member)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group created successfully",
		"group":   gin.H{"id": group.ID, "name": group.Name, "created_by": group.CreatedBy},
	})
}

// AddMember — POST /groups/:id/members
func AddMember(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify group exists
	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Verify user exists
	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Check if already a member
	var existing models.GroupMember
	result := config.DB.Where("group_id = ? AND user_id = ?", groupID, input.UserID).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member of this group"})
		return
	}

	member := models.GroupMember{GroupID: uint(groupID), UserID: input.UserID}
	config.DB.Create(&member)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Member added successfully",
		"group_id": groupID,
		"user_id":  input.UserID,
	})
}

// GetGroup — GET /groups/:id
func GetGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Fetch members with user details
	var members []models.GroupMember
	config.DB.Where("group_id = ?", groupID).Find(&members)

	var memberDetails []gin.H
	for _, m := range members {
		var user models.User
		config.DB.First(&user, m.UserID)
		memberDetails = append(memberDetails, gin.H{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"group": gin.H{
			"id":         group.ID,
			"name":       group.Name,
			"created_by": group.CreatedBy,
			"created_at": group.CreatedAt,
			"members":    memberDetails,
		},
	})
}
