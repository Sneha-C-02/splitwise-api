# prompts_used.md — AI Assistance Transparency Log

Project: Splitwise-style Expense Tracking REST API  
Capstone Project 7 — Go (Golang)

---

## Purpose of This Document

This document records the responsible use of AI-assisted development tools during the implementation of this project.

Per the project guidelines, AI assistance is permitted provided all prompts are disclosed transparently.

AI tools were used as implementation accelerators and reference assistants — not as autonomous project builders.

All architectural decisions, financial modeling strategy, algorithm selection, validation rules, and system structure were defined, verified, and validated by the developer.

---

# AI Tool Used

- Conversational AI coding assistant (LLM-based code generation tool)

Total Implementation Phases: 7  
All outputs were reviewed, tested, and understood before integration.

---

# Prompts Used During Development

Below are the structured prompts that guided the AI-generated portions of the implementation.

---

## Prompt 1 — Project Bootstrapping

**Purpose:** Initialize project scaffold.

**Prompt Used:**

> Initialize a Go REST API project named `splitwise-api` using Gin and GORM with a pure Go SQLite driver.  
> Create a User model with name, email, and password.  
> Implement `/register` and `/users` endpoints.  
> Enable database auto-migration.

**Output Integrated:**
- `main.go`
- `config/database.go`
- `models/user.go`
- Basic handler for registration and listing users

**Developer Decisions:**
- SQLite chosen for zero-config demonstration
- Folder structure finalized manually
- Naming conventions defined manually

---

## Prompt 2 — Authentication Hardening

**Purpose:** Improve password security and response safety.

**Prompt Used:**

> Modify the Register endpoint to hash passwords using bcrypt before saving.  
> Ensure password hashes are never returned in API responses.  
> Add basic validation for email and minimum password length.

**Output Integrated:**
- bcrypt hashing in `handlers/auth.go`
- Input struct validation
- Safe JSON response model excluding password

**Developer Validation:**
- Confirmed hash storage
- Confirmed password never exposed
- Tested via Postman

---

## Prompt 3 — Groups & Membership

**Purpose:** Implement group functionality.

**Prompt Used:**

> Implement Group and GroupMember models.  
> Auto-add group creator as a member when creating a group.  
> Implement:
> - POST `/groups`
> - POST `/groups/:id/members`
> - GET `/groups/:id`

**Output Integrated:**
- `models/group.go`
- `handlers/groups.go`

**Developer Decisions:**
- Prevent duplicate membership
- Dynamically load members in GET responses
- Enforce membership validation

---

## Prompt 4 — Expenses with Financial Precision

**Purpose:** Add expense handling with correct money modeling.

**Prompt Used:**

> Implement Expense and ExpenseSplit models using int64 for monetary values.  
> Support equal split logic without using floats.  
> Ensure remainders are distributed correctly.  
> Implement:
> - POST `/groups/:id/expenses`
> - GET `/groups/:id/expenses`
> - DELETE `/expenses/:id`

**Output Integrated:**
- `models/expense.go`
- `handlers/expenses.go`

**Critical Developer Decisions:**
- Use `int64` paise (never float64)
- Integer-only math
- Strict sum validation for exact splits
- Remainder paise distribution logic

---

## Prompt 5 — Balances & Settlement Algorithm

**Purpose:** Compute balances and minimize settlement transactions.

**Prompt Used:**

> Implement:
> - GET `/groups/:id/balances`
> - Greedy settlement algorithm to minimize transactions
> - GET `/groups/:id/settlements`
> Sort creditors and debtors by descending absolute balance and match greedily.

**Output Integrated:**
- `handlers/settlements.go`
- `algorithms/settlement.go`

**Developer Decisions:**
- Use Greedy Minimization (O(n log n))
- Do not store settlements in database
- Compute dynamically for integrity

**Edge Cases Tested:**
- Single creditor, multiple debtors
- Multiple creditors
- Perfectly balanced group
- Large remainder distributions

---

## Prompt 6 — Enhancements (Percentage, Exact, Summary)

**Purpose:** Extend real-world split support.

**Prompt Used:**

> Enhance AddExpense to support:
> - Percentage split (must sum to 100)
> - Exact split (sum must equal expense)
> Add `/users/:id/summary` endpoint to aggregate global balances.  
> Do not modify existing equal split behavior.

**Output Integrated:**
- Percentage logic with integer math
- Exact split validation
- `handlers/summary.go`

**Developer Decisions:**
- Last-user remainder handling for percentage splits
- No totals stored in DB
- Summary reuses dynamic balance computation

---

## Prompt 7 — Validation Hardening

**Purpose:** Prevent invalid financial states.

**Prompt Used:**

> Add validation guards in AddExpense:
> - Amount must be > 0
> - split_type is required  
> Ensure early exit without affecting existing logic.

**Output Integrated:**
- Guard clauses in `handlers/expenses.go`

**Developer Reasoning:**
- Prevent negative ledger states
- Maintain API contract integrity
- Early return defensive pattern

---

# What AI Was NOT Used For

AI was NOT used to:

- Choose paise over float strategy
- Design greedy minimization concept
- Design system architecture
- Determine API endpoint structure
- Define validation philosophy
- Create financial integrity rules
- Design testing scenarios
- Write conceptual explanations in docs

All high-level reasoning and modeling decisions were developer-driven.

---

# Final Verification Checklist

Before submission:

- All endpoints tested via Postman
- Negative validations tested
- Multi-creditor edge case validated
- `go build` passes with zero errors
- Database auto-migration verified
- No passwords returned in responses
- All monetary values stored as int64 paise
- Settlement algorithm manually validated

---

# Statement of Responsibility

This project reflects the developer’s understanding of:

- REST API architecture  
- Relational data modeling  
- Financial precision handling  
- Greedy algorithm optimization  
- Validation hardening  
- Secure password storage  
- Dynamic balance computation  

AI tools were used as implementation assistants, and all integrated code has been reviewed, tested, and fully understood.

All implemented logic can be confidently explained during evaluation.