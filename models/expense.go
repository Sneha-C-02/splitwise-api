package models

import "gorm.io/gorm"

// Expense represents a shared expense paid by one member of a group.
// Amount is stored in paise (int64) to avoid float precision errors.
// Example: ₹100.50 = 10050 paise
type Expense struct {
	gorm.Model
	GroupID     uint           `json:"group_id" gorm:"not null"`
	PaidBy      uint           `json:"paid_by" gorm:"not null"`
	Amount      int64          `json:"amount" gorm:"not null"` // in paise
	Description string         `json:"description"`
	Splits      []ExpenseSplit `json:"splits,omitempty" gorm:"foreignKey:ExpenseID"`
}

// ExpenseSplit records how much each member owes for a given expense.
// AmountOwed is in paise (int64).
type ExpenseSplit struct {
	gorm.Model
	ExpenseID  uint  `json:"expense_id" gorm:"not null"`
	UserID     uint  `json:"user_id" gorm:"not null"`
	AmountOwed int64 `json:"amount_owed" gorm:"not null"` // in paise
}
