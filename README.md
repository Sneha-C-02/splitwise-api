# Splitwise API — Expense Tracker & Bill Splitter

A production-quality REST API built in Go for tracking shared expenses and splitting bills among friends — similar to Splitwise.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (Golang) |
| Web Framework | Gin |
| ORM | GORM |
| Database | SQLite (pure Go driver) |
| Password Hashing | bcrypt |
| Money Handling | `int64` paise (no floats) |

---

## Prerequisites

- Go 1.21 or higher
- Git
- Any OS (Windows, macOS, Linux)
- No external database required (SQLite is embedded)

---

## Setup & Run

```bash
# 1. Clone / enter the project directory
cd splitwise-api

# 2. Install dependencies
go mod tidy
# All dependencies are managed via Go modules and declared in go.mod.

# 3. Run the server
go run main.go
# Server starts at http://localhost:8080
```

The SQLite database (`splitwise.db`) is auto-created and all tables are auto-migrated on startup.

---

## API Endpoints

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ping` | Server health check |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register a new user |
| GET | `/users` | Get all users |

### Groups
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/groups` | Create a group |
| POST | `/groups/:id/members` | Add a member to a group |
| GET | `/groups/:id` | Get group details + members |

### Expenses
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/groups/:id/expenses` | Add an expense (equal, percentage, or exact split) |
| GET | `/groups/:id/expenses` | List all expenses in a group |
| DELETE | `/expenses/:id` | Delete an expense |

### Balances & Settlements
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/groups/:id/balances` | Net balance per user in a specific group |
| GET | `/groups/:id/settlements` | Optimized settlement transactions for a group |
| GET | `/users/:id/summary` | User's global financial position across ALL groups |

---

## Example Requests (curl)

### Register a user
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Priya","email":"priya@example.com","password":"secret123"}'
```

### Create a group
```bash
curl -X POST http://localhost:8080/groups \
  -H "Content-Type: application/json" \
  -d '{"name":"Goa Trip","created_by":1}'
```

### Add member to group
```bash
curl -X POST http://localhost:8080/groups/1/members \
  -H "Content-Type: application/json" \
  -d '{"user_id":2}'
```

### Add an expense — Equal split
```bash
curl -X POST http://localhost:8080/groups/1/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "paid_by": 1,
    "amount": 60000,
    "description": "Hotel dinner",
    "split_type": "equal"
  }'
```
> `amount` is in **paise**: `60000` = ₹600

### Add an expense — Percentage split
```bash
curl -X POST http://localhost:8080/groups/1/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "paid_by": 1,
    "amount": 100000,
    "description": "Dinner %",
    "split_type": "percentage",
    "splits": [
        {"user_id": 1, "percentage": 50},
        {"user_id": 2, "percentage": 25},
        {"user_id": 3, "percentage": 25}
    ]
  }'
```

### Add an expense — Exact split
```bash
curl -X POST http://localhost:8080/groups/1/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "paid_by": 2,
    "amount": 100000,
    "description": "Rent Exact",
    "split_type": "exact",
    "splits": [
        {"user_id": 1, "amount": 40000},
        {"user_id": 3, "amount": 30000},
        {"user_id": 4, "amount": 30000}
    ]
  }'
```

### Get summary across groups
```bash
curl http://localhost:8080/users/1/summary
```

### Get balances
```bash
curl http://localhost:8080/groups/1/balances
```

### Get optimized settlements
```bash
curl http://localhost:8080/groups/1/settlements
```

---

## Example Settlement Scenario

**Group: Goa Trip (4 members — Priya, Sneha, Arun, Raj)**

| Event | Payer | Amount |
|-------|-------|--------|
| Hotel | Priya | ₹600 (split equally → ₹150 each) |
| Food | Sneha | ₹200 (split equally → ₹50 each) |
| Taxi  | Arun  | ₹100 (split equally → ₹25 each) |

**Net Balances:**

| Person | Paid | Owes | Net |
|--------|------|------|-----|
| Priya  | ₹600 | ₹225 | **+₹375** |
| Sneha  | ₹200 | ₹225 | **-₹25** |
| Arun   | ₹100 | ₹225 | **-₹125** |
| Raj    | ₹0   | ₹225 | **-₹225** |

**Optimized Settlements (3 transactions, minimized):**
```
Raj   → Priya  ₹225
Arun  → Priya  ₹125
Sneha → Priya  ₹25
```

Without optimization, this group might require up to 6 transactions. The greedy algorithm reduces it to the minimum.

---

## Settlement Algorithm

**Time Complexity: O(n log n)**

1. Compute net balance per user (paid − owed)
2. Split into **creditors** (positive) and **debtors** (negative)
3. Sort both lists descending by absolute amount
4. Greedily match the largest debtor to the largest creditor
5. Transfer `min(debt, credit)`, update both balances
6. Advance pointer for whichever side reaches zero
7. Repeat until all balances are zero

This greedy approach is proven to yield the minimum number of transactions.

---

## Money Handling

All amounts are stored as **`int64` in paise** (1 INR = 100 paise).

- ₹100.50 → stored as `10050`
- **No `float64` anywhere** in the codebase
- Rounding on equal splits: remainder paise distributed 1 at a time to the first N members
- Example: ₹100 split 3 ways → `3334 + 3333 + 3333 paise`

See [`docs/MONEY_HANDLING.md`](docs/MONEY_HANDLING.md) for full explanation.

---

## Validation Rules

- Amount must be greater than 0  
- `split_type` is required  
- Percentage splits must sum exactly to 100  
- Exact split amounts must equal the total expense amount  
- Only group members can be included in expense splits  
- Payer must belong to the group  

---

## Security Considerations

- Passwords are stored using bcrypt hashing  
- Password hashes are never returned in API responses  
- Financial calculations avoid floating-point arithmetic  
- All balances are computed dynamically (no redundant stored totals)

---

## Project Structure

```
splitwise-api/
├── main.go                   # Entry point + all routes
├── config/
│   └── database.go           # GORM + SQLite setup + AutoMigrate
├── models/
│   ├── user.go               # User model
│   ├── group.go              # Group + GroupMember models
│   └── expense.go            # Expense + ExpenseSplit models
├── handlers/
│   ├── auth.go               # Register, GetUsers
│   ├── groups.go             # CreateGroup, AddMember, GetGroup
│   ├── expenses.go           # AddExpense, GetExpenses, DeleteExpense
│   ├── settlements.go        # GetBalances, GetSettlements
│   └── summary.go            # Global summary endpoint
├── algorithms/
│   └── settlement.go         # Greedy minimization algorithm
├── docs/
│   ├── DESIGN.md             # Architecture & DB schema
│   ├── MONEY_HANDLING.md     # Money handling strategy
│   └── prompts_used.md       # AI transparency log
└── README.md
```
