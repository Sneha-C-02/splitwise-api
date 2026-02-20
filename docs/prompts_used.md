# prompts_used.md — AI Transparency Log

This document records how AI assistance (Antigravity / Google DeepMind) was used in building this project, in the spirit of academic transparency.

---

## Project: Expense Tracker with Bill Splitting API (Capstone Project 7)

### Overall Approach

The project was built phase-by-phase with AI assistance. The human developer defined the architecture, requirements, and constraints. The AI generated the code, which was reviewed before use.

---

## Prompts Used

### Prompt 1 — Project Bootstrap (Phase 1)
> Initialize a Go REST API project named `splitwise-api` with Gin, GORM, and a pure Go SQLite driver. Create a User model with name, email, and password. Implement `/register` and `/users` endpoints. Enable database auto-migration.

**AI Output Used:** initial main.go, config/database.go, models/user.go, handlers/auth.go

---

### Prompt 2 — Phase 1 Fix: Password Hashing
> The Register handler saves passwords as plain text. Fix it to hash passwords with bcrypt before saving. Never return the hashed password in API responses. Validate email and enforce minimum 6-character passwords.

**AI Output Used:** Updated handlers/auth.go with bcrypt hashing, input struct validation, safe response (no password field).

---

### Prompt 3 — Phase 2: Groups
> Implement a Group model (ID, Name, CreatedBy, CreatedAt) and GroupMember model (ID, GroupID, UserID). Auto-add the creator as a member when a group is created. Implement POST /groups, POST /groups/:id/members (with duplicate check), GET /groups/:id (with member list).

**AI Output Used:** models/group.go, handlers/groups.go

---

### Prompt 4 — Phase 3: Expenses with Paise
> Implement Expense model (GroupID, PaidBy, Amount as int64 paise, Description) and ExpenseSplit model (ExpenseID, UserID, AmountOwed as int64 paise). Support equal split (with proper remainder distribution — no floats) and custom split (with sum validation). Implement POST /groups/:id/expenses, GET /groups/:id/expenses, DELETE /expenses/:id.

**AI Output Used:** models/expense.go, handlers/expenses.go

---

### Prompt 5 — Phase 4 & 5: Balances and Settlement Algorithm
> Implement GET /groups/:id/balances that returns net balance per user (paid − owed). Implement a greedy settlement algorithm in algorithms/settlement.go that minimizes the number of transactions by sorting creditors and debtors by descending absolute balance and greedily matching them. Implement GET /groups/:id/settlements that returns the optimized transaction list with user names.

**AI Output Used:** handlers/settlements.go, algorithms/settlement.go

---

### Prompt 7 — Enhancements: Percentage, Exact, and Summary
> Surgically enhance the working system to support `percentage` and `exact` split types without touching existing `equal` logic. Percentages must sum to 100 and use integer math. Exact amounts must sum to total. Add a `/users/:id/summary` endpoint that aggregates a user's net position across all groups by reusing existing balance logic. Do not store totals in DB.

---

### Prompt 8 — Hardening: Validation Guards in AddExpense
> Surgically harden `AddExpense` by adding guard clauses for `Amount <= 0` and `SplitType == ""`. Guards must execute before split-type processing to ensure early exit. Do not modify any other logic.

**AI Output Used:** Updated handlers/expenses.go with minimal guard clauses.

---

## Human Decisions (Not AI-Generated)

- Choice of SQLite over PostgreSQL for zero-config demo
- Decision to use paise (int64) instead of decimal library
- Phased implementation order
- Decision to auto-add creator as group member
- Endpoint naming conventions
- Project folder structure
