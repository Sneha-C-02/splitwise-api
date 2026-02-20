package handlers

import (
	"net/http"
	"splitwise-api/algorithms"
	"splitwise-api/config"
	"splitwise-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBalances — GET /groups/:id/balances
// Net balance per user = total paid − total owed
// Positive = creditor, Negative = debtor
func GetBalances(c *gin.Context) {
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

	netBalances := computeNetBalances(uint(groupID))

	// Fetch user names for readability
	var result []gin.H
	for uid, bal := range netBalances {
		var user models.User
		config.DB.First(&user, uid)
		status := "settled"
		if bal > 0 {
			status = "creditor"
		} else if bal < 0 {
			status = "debtor"
		}
		result = append(result, gin.H{
			"user_id": uid,
			"name":    user.Name,
			"balance": bal, // in paise
			"status":  status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id": groupID,
		"balances": result,
		"note":     "Amounts are in paise. Divide by 100 for INR.",
	})
}

// GetSettlements — GET /groups/:id/settlements
// Returns the minimum set of transactions to settle all debts.
func GetSettlements(c *gin.Context) {
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

	netBalances := computeNetBalances(uint(groupID))
	transactions := algorithms.MinimizeTransactions(netBalances)

	// Enrich with user names
	var result []gin.H
	for _, tx := range transactions {
		var from, to models.User
		config.DB.First(&from, tx.From)
		config.DB.First(&to, tx.To)
		result = append(result, gin.H{
			"from":         tx.From,
			"from_name":    from.Name,
			"to":           tx.To,
			"to_name":      to.Name,
			"amount_paise": tx.Amount,
			"amount_inr":   formatINR(tx.Amount),
		})
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"group_id":     groupID,
			"transactions": []gin.H{},
			"message":      "All debts are settled! 🎉",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":                groupID,
		"transactions":            result,
		"total_transaction_count": len(result),
		"algorithm":               "Greedy minimization — O(n log n)",
	})
}

// computeNetBalances calculates net balance per user for a group.
// Net = total paid − total owed (in paise)
func computeNetBalances(groupID uint) map[uint]int64 {
	netBalances := make(map[uint]int64)

	// Credit: what each user paid
	var expenses []models.Expense
	config.DB.Where("group_id = ?", groupID).Find(&expenses)
	for _, e := range expenses {
		netBalances[e.PaidBy] += e.Amount
	}

	// Debit: what each user owes across all expense splits
	var splits []models.ExpenseSplit
	config.DB.
		Joins("JOIN expenses ON expenses.id = expense_splits.expense_id").
		Where("expenses.group_id = ? AND expenses.deleted_at IS NULL AND expense_splits.deleted_at IS NULL", groupID).
		Find(&splits)
	for _, s := range splits {
		netBalances[s.UserID] -= s.AmountOwed
	}

	return netBalances
}

// formatINR converts paise (int64) to a readable INR string like "₹100.50"
func formatINR(paise int64) string {
	rupees := paise / 100
	paiseRemainder := paise % 100
	if paiseRemainder == 0 {
		return "₹" + strconv.FormatInt(rupees, 10) + ".00"
	}
	paiseStr := strconv.FormatInt(paiseRemainder, 10)
	if paiseRemainder < 10 {
		paiseStr = "0" + paiseStr
	}
	return "₹" + strconv.FormatInt(rupees, 10) + "." + paiseStr
}
