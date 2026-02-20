# prompts_used.md — AI Assistance Transparency Log

## Project: Splitwise-style Expense Tracking REST API  
Capstone Project 7 — Go (Golang)

---

## Purpose of This Document

This document records the responsible use of AI-assisted development tools during the implementation of this project, in accordance with the project guidelines permitting AI usage with full disclosure.

AI tools were used as implementation assistants — not as autonomous project builders.

All architectural design, algorithm selection, financial modeling decisions, and system structure were defined and validated by the developer.

---

## Development Philosophy

AI was used in the same way modern software engineers use tools such as:
- StackOverflow
- Official documentation
- Code templates
- Linting and debugging assistants

Every AI-generated snippet was:
- Reviewed
- Understood
- Tested manually using Postman
- Modified where necessary
- Validated against project requirements

No code was included without understanding its functionality.

---

## Areas Where AI Assistance Was Used

### 1. Project Bootstrapping

AI assistance was used to:
- Generate initial Gin + GORM project scaffold
- Set up database connection
- Create basic User model
- Implement initial endpoints (`/register`, `/users`)

Human decisions:
- Folder structure
- Dependency selection
- SQLite instead of PostgreSQL (for zero-config demo)

---

### 2. Authentication Hardening

AI assistance was used to:
- Integrate bcrypt hashing
- Remove password from API responses
- Add basic input validation

Human decisions:
- Never return hashed password
- Enforce safe response models
- Maintain separation of input vs output structs

---

### 3. Group and Membership Design

AI helped generate:
- Group model
- GroupMember join model
- Basic CRUD handlers

Human design decisions:
- Auto-add creator as group member
- Prevent duplicate membership
- Dynamic membership loading for GET endpoints

---

### 4. Expense Modeling and Money Handling

AI assistance supported:
- Struct definitions for Expense and ExpenseSplit
- Equal split implementation
- Sum validation for split integrity

Human design decisions (Critical):
- Use `int64` paise instead of float or decimal
- Never use floating-point math in financial logic
- Remainder distribution logic for equal splits
- Strict sum validation for exact splits
- Integer-only math for percentage splits

Financial integrity decisions were made independently and intentionally.

---

### 5. Settlement Algorithm

AI assistance was used to:
- Structure the greedy matching function
- Implement sorting logic

Human algorithm decisions:
- Use greedy minimization strategy
- Complexity target: O(n log n)
- Sort creditors and debtors by descending absolute balance
- Match largest debtor to largest creditor iteratively
- Avoid storing settlements in database
- Compute dynamically to ensure consistency

Algorithm logic was reviewed and manually tested with edge-case scenarios:
- Single creditor, multiple debtors
- Multiple creditors and debtors
- Perfectly balanced group
- Exact remainder matching

---

### 6. Enhancements

AI-assisted implementation included:
- Percentage-based splits
- Exact-amount splits
- Global user summary endpoint

Human decisions:
- Percentage must sum to 100
- Exact splits must equal total expense
- Early validation guards
- Reuse existing balance logic (no redundancy)
- No global totals stored in database

---

### 7. Validation Hardening

AI assistance was used to implement:
- Guard: `Amount > 0`
- Guard: `split_type required`
- Early return validation pattern

Human validation reasoning:
- Prevent negative expense creation
- Ensure API contract consistency
- Avoid invalid ledger state

---

## What AI Was NOT Used For

AI was not used to:
- Choose financial modeling strategy
- Decide on paise vs float
- Design overall architecture
- Decide on greedy minimization approach
- Determine API structure
- Write conceptual explanations in documentation
- Design testing scenarios

All high-level system reasoning was developer-driven.

---

## Final Verification

Before submission:

- All endpoints tested manually via Postman
- Edge-case scenarios validated
- Negative input cases validated
- Settlement logic verified with multi-creditor case
- `go build` completed with zero errors
- Database auto-migration verified
- No passwords returned in responses
- All money stored as `int64` paise

---

## Statement of Responsibility

This project reflects the developer’s understanding of:

- REST API design
- Relational modeling
- Financial precision handling
- Greedy algorithm optimization
- Validation and input hardening
- Secure password storage
- Dynamic balance computation

AI tools were used as coding accelerators — not as substitutes for understanding.

All implemented features can be explained and defended during evaluation.