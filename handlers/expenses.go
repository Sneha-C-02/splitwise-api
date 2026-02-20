package handlers

import (
	"net/http"
	"splitwise-api/config"
	"splitwise-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// splitEntry is used for percentage and exact split inputs
type splitEntry struct {
	UserID     uint  `json:"user_id"`
	Percentage int64 `json:"percentage"` // for percentage split
	Amount     int64 `json:"amount"`     // for exact split (paise)
}

// AddExpense — POST /groups/:id/expenses
// Supports: "equal", "percentage", "exact" split types.
// All amounts are in PAISE (int64). No floats anywhere.
func AddExpense(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var input struct {
		PaidBy      uint         `json:"paid_by" binding:"required"`
		Amount      int64        `json:"amount" binding:"required"` // in paise
		Description string       `json:"description"`
		SplitType   string       `json:"split_type"` // "equal", "percentage", "exact"
		Splits      []splitEntry `json:"splits"`     // used for percentage and exact
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ── VALIDATION GUARDS ────────────────────────────────────────────────
	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	if input.SplitType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "split_type is required"})
		return
	}
	// ─────────────────────────────────────────────────────────────────────

	// Verify group exists
	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Verify payer is a member
	var payer models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, input.PaidBy).First(&payer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payer is not a member of this group"})
		return
	}

	// Fetch all group members (needed for equal split)
	var members []models.GroupMember
	config.DB.Where("group_id = ?", groupID).Find(&members)
	memberCount := int64(len(members))
	if memberCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group has no members"})
		return
	}

	// Create Expense record
	expense := models.Expense{
		GroupID:     uint(groupID),
		PaidBy:      input.PaidBy,
		Amount:      input.Amount,
		Description: input.Description,
	}
	config.DB.Create(&expense)

	var splits []models.ExpenseSplit

	switch input.SplitType {

	// ── EQUAL SPLIT (original logic — UNTOUCHED) ─────────────────────────
	case "equal", "":
		// Equal split with rounding correction.
		// e.g., ₹100 among 3 → 3334, 3333, 3333 paise  (total = 10000 ✅)
		baseShare := input.Amount / memberCount
		remainder := input.Amount % memberCount
		for idx, m := range members {
			share := baseShare
			if int64(idx) < remainder {
				share++ // distribute 1 extra paise to first `remainder` members
			}
			splits = append(splits, models.ExpenseSplit{
				ExpenseID:  expense.ID,
				UserID:     m.UserID,
				AmountOwed: share,
			})
		}

	// ── PERCENTAGE SPLIT ──────────────────────────────────────────────────
	case "percentage":
		if len(input.Splits) == 0 {
			config.DB.Delete(&expense)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provide splits[] for percentage split"})
			return
		}
		// Validate: percentages must sum to exactly 100
		var totalPct int64
		for _, s := range input.Splits {
			totalPct += s.Percentage
		}
		if totalPct != 100 {
			config.DB.Delete(&expense)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Percentages must sum to 100",
				"got":      totalPct,
				"expected": 100,
			})
			return
		}
		// Integer math only — no float64.
		// Remainder assigned to last user to guarantee total conservation.
		var allocated int64
		for idx, s := range input.Splits {
			var share int64
			if idx == len(input.Splits)-1 {
				share = input.Amount - allocated // absorb any rounding remainder
			} else {
				share = input.Amount * s.Percentage / 100
			}
			allocated += share
			splits = append(splits, models.ExpenseSplit{
				ExpenseID:  expense.ID,
				UserID:     s.UserID,
				AmountOwed: share,
			})
		}

	// ── EXACT SPLIT ───────────────────────────────────────────────────────
	case "exact":
		if len(input.Splits) == 0 {
			config.DB.Delete(&expense)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provide splits[] for exact split"})
			return
		}
		// Validate: exact amounts must sum to total expense amount
		var total int64
		for _, s := range input.Splits {
			total += s.Amount
		}
		if total != input.Amount {
			config.DB.Delete(&expense)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Exact split amounts do not sum to total expense amount",
				"expected": input.Amount,
				"got":      total,
			})
			return
		}
		for _, s := range input.Splits {
			splits = append(splits, models.ExpenseSplit{
				ExpenseID:  expense.ID,
				UserID:     s.UserID,
				AmountOwed: s.Amount,
			})
		}

	default:
		config.DB.Delete(&expense)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid split_type. Must be one of: equal, percentage, exact",
		})
		return
	}

	// Bulk insert splits
	config.DB.Create(&splits)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Expense added successfully",
		"expense": gin.H{
			"id":          expense.ID,
			"group_id":    expense.GroupID,
			"paid_by":     expense.PaidBy,
			"amount":      expense.Amount,
			"split_type":  input.SplitType,
			"description": expense.Description,
			"splits":      splits,
		},
	})
}

// GetExpenses — GET /groups/:id/expenses
func GetExpenses(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var expenses []models.Expense
	config.DB.Where("group_id = ?", groupID).Preload("Splits").Find(&expenses)

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

// DeleteExpense — DELETE /expenses/:id
func DeleteExpense(c *gin.Context) {
	expenseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var expense models.Expense
	if err := config.DB.First(&expense, expenseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Delete associated splits first, then the expense
	config.DB.Where("expense_id = ?", expenseID).Delete(&models.ExpenseSplit{})
	config.DB.Delete(&expense)

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
