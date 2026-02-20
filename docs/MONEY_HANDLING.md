# MONEY_HANDLING.md — Money Handling Strategy

## The Problem With Floats

Floating-point numbers (float32, float64) cannot represent all decimal values exactly in binary.

```go
// Go example — float64 failure
fmt.Println(0.1 + 0.2)
// Output: 0.30000000000000004
```

In financial software, this causes:
- Wrong balances (e.g., ₹99.99 becoming ₹100.00000000001)
- Split miscalculations (e.g., splitting ₹100 among 3 gives ₹33.333...)
- Cumulative drift over many transactions
- Audit failures (amounts don't add up exactly)

**This is unacceptable in a money-handling system.**

---

## Our Approach: Store Paise as `int64`

The Indian Rupee has 2 decimal places (paise). We multiply all amounts by 100 and store as integers.

| Human Input | Stored Value | Type |
|-------------|--------------|------|
| ₹100.00 | 10000 | int64 |
| ₹100.50 | 10050 | int64 |
| ₹33.33 | 3333 | int64 |
| ₹0.01 | 1 | int64 |

All arithmetic is then **exact integer arithmetic**. No approximations.

---

## Rounding Strategy for Uneven Splits

When ₹100 (10000 paise) is split equally among 3 people:

```
10000 ÷ 3 = 3333 remainder 1
```

Naive approach: give everyone ₹33.33 → total = ₹99.99 ≠ ₹100 ❌

**Our approach (remainder distribution):**

```go
baseShare := amount / memberCount   // 3333
remainder := amount % memberCount   // 1

// First `remainder` members get 1 extra paise
// Member 1: 3334 paise = ₹33.34
// Member 2: 3333 paise = ₹33.33
// Member 3: 3333 paise = ₹33.33
// Total:   10000 paise = ₹100.00 ✅
```

The remainder (at most `memberCount - 1` paise, i.e., a few paise) is distributed 1 paise at a time to the first N members. This is the standard approach used by financial systems.

---

## Rounding Strategy for Percentage Splits

When using `split_type: "percentage"`, we use pure integer math:
`amount * percentage / 100`

To ensure **total conservation of money**, the **last user** in the split list always absorbs any rounding remainder:
`share = total_amount - allocated_so_far`

Example: ₹1000 with 33%, 33%, 34% split
1. User 1 (33%): 100000 * 33 / 100 = 33000
2. User 2 (33%): 100000 * 33 / 100 = 33000
3. User 3 (last): 100000 - (33000 + 33000) = 34000
Total: 100000 paise ✅

## Exact Split Validation

For `split_type: "exact"`, the API strictly validates:
`sum(individual split amounts) == total expense amount`

If the sum of paise provided for each user does not exactly equal the total expense amount, the request is rejected with a **400 Bad Request** error.

---

## Custom Split Validation

When using custom splits, the API validates:

```
sum(split amounts) == total expense amount
```

If they don't match, the expense is rejected with an error showing the expected vs actual total. This enforces **conservation of money** — every rupee paid is accounted for.

---

## API Contract

All amounts in API requests and responses are in **paise (int64)**:

```json
// Request: Add ₹600 expense
{ "amount": 60000 }

// Response: Balance of ₹375
{ "balance": 37500 }

// Response: Settlement of ₹225
{
  "amount_paise": 22500,
  "amount_inr": "₹225.00"
}
```

The response also includes a human-readable `amount_inr` field for convenience, but the source of truth is always `amount_paise`.

---

## Why Not Use a Decimal Library?

Libraries like `shopspring/decimal` are excellent for production financial systems that need arbitrary precision or currency conversion. For this project scope (INR only, 2 decimal places), `int64` paise is simpler, dependency-free, and equally correct.
