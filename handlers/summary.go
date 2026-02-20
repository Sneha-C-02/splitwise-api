package handlers

import (
	"net/http"
	"splitwise-api/config"
	"splitwise-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserSummary — GET /users/:id/summary
// Returns a user's financial position across ALL groups they belong to.
// Reuses computeNetBalances() — no logic duplicated, nothing stored in DB.
func GetUserSummary(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify user exists
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Find all groups this user belongs to
	var memberships []models.GroupMember
	config.DB.Where("user_id = ?", userID).Find(&memberships)

	var totalOwedToUser int64 // user is creditor in these amounts
	var totalUserOwes int64   // user is debtor in these amounts

	for _, m := range memberships {
		// Reuse existing balance function for each group
		balances := computeNetBalances(m.GroupID)
		netInGroup := balances[uint(userID)]

		if netInGroup > 0 {
			totalOwedToUser += netInGroup
		} else if netInGroup < 0 {
			totalUserOwes += -netInGroup // store as positive
		}
	}

	netBalance := totalOwedToUser - totalUserOwes

	status := "settled"
	if netBalance > 0 {
		status = "creditor"
	} else if netBalance < 0 {
		status = "debtor"
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":                  userID,
		"name":                     user.Name,
		"total_owed_to_user_paise": totalOwedToUser,
		"total_user_owes_paise":    totalUserOwes,
		"net_balance_paise":        netBalance,
		"status":                   status,
		"note":                     "Amounts are in paise. Divide by 100 for INR.",
	})
}
