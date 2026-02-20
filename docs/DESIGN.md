# DESIGN.md — System Design & Architecture

## Overview

This is a backend REST API for shared expense tracking and bill splitting, built in Go using Gin and GORM with a SQLite database.

---

## Database Schema

### `users`
| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER (PK) | Auto-increment |
| name | TEXT | Required |
| email | TEXT | Unique, required |
| password | TEXT | bcrypt hash, never plain text |
| created_at | DATETIME | Auto |
| updated_at | DATETIME | Auto |
| deleted_at | DATETIME | Soft delete (GORM) |

### `groups`
| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER (PK) | Auto-increment |
| name | TEXT | Required |
| created_by | INTEGER (FK → users.id) | Creator user |
| created_at | DATETIME | Auto |

### `group_members`
| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER (PK) | Auto-increment |
| group_id | INTEGER (FK → groups.id) | Required |
| user_id | INTEGER (FK → users.id) | Required |
| created_at | DATETIME | Joined timestamp |

### `expenses`
| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER (PK) | Auto-increment |
| group_id | INTEGER (FK → groups.id) | Required |
| paid_by | INTEGER (FK → users.id) | Who paid |
| amount | INTEGER (int64) | **In paise**, not rupees |
| description | TEXT | Optional |
| created_at | DATETIME | Auto |
| deleted_at | DATETIME | Soft delete |

### `expense_splits`
| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER (PK) | Auto-increment |
| expense_id | INTEGER (FK → expenses.id) | Required |
| user_id | INTEGER (FK → users.id) | Who owes |
| amount_owed | INTEGER (int64) | **In paise** |

---

## Why `int64` for Money?

Floating point (float64) cannot represent all decimal values exactly in binary.

```
0.1 + 0.2 = 0.30000000000000004  // float64 failure
```

By storing all amounts as **paise (integer)**, we guarantee:
- **Exact arithmetic** — no rounding surprises
- **Correct split calculations** — integer division with remainder correction
- **Auditability** — every paisa is accounted for

---

## Why Greedy Settlement Algorithm?

The debt minimization problem is equivalent to finding the minimum number of edges to zero out all balances in a flow graph.

The greedy approach:
- Is **O(n log n)** — fast enough for any real group size
- Is **proven optimal** for the case where we just minimize transaction count (not considering who pays whom)
- Is simple to understand and audit
- Handles any combination of creditors and debtors

Alternative approaches (e.g., graph-based flow algorithms) are more complex with no practical benefit at this scale.

---

## Architecture Decisions

### Why Gin?
- Lightweight, fast HTTP router
- Idiomatic Go
- Battle-tested in production

### Why GORM + SQLite?
- GORM provides clean model-to-table mapping and auto-migration
- SQLite is zero-config, file-based — ideal for a capstone/demo project
- Easy to swap to PostgreSQL/MySQL by changing the driver

### Why modular folder structure?
```
models/    — pure data structs
handlers/  — HTTP request handling
algorithms/— pure business logic (no HTTP dependency)
config/    — database setup only
```

### Why dynamic user summary?
The `/users/:id/summary` endpoint iterates through all of a user's group memberships and calls `computeNetBalances` for each. This ensures the summary is always **fresh** without needing to sync redundant totals in the database, matching the "single source of truth" philosophy. This keeps concerns separated and makes each layer independently testable.

### Why bcrypt?
- Industry standard for password hashing
- One-way (hashes cannot be reversed)
- Slow by design (resists brute-force attacks)
- Salt is built-in (no duplicate hashes for same password)

---

## Security Decisions

- Passwords are **never stored in plain text**
- Passwords are **never returned** in API responses
- Input validation on all endpoints
- Duplicate membership checks before adding group members
- Soft deletes on Expenses (GORM's `deleted_at`) — data preserved for audit
