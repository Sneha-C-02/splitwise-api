package algorithms

import "sort"

// Transaction represents a single optimized settlement payment
type Transaction struct {
	From   uint  `json:"from"`
	To     uint  `json:"to"`
	Amount int64 `json:"amount"` // in paise
}

type balance struct {
	UserID uint
	Amount int64
}

// MinimizeTransactions takes a map of userID -> net balance (in paise)
// and returns the minimum set of transactions to settle all debts.
//
// Algorithm (Greedy, O(n log n)):
//  1. Separate users into creditors (balance > 0) and debtors (balance < 0).
//  2. Sort both by absolute value descending.
//  3. Match the largest debtor to the largest creditor.
//  4. Transfer min(|debt|, credit). Reduce both balances.
//  5. If one side reaches 0, advance its pointer.
//  6. Repeat until all balances are zero.
//
// This produces the minimum number of transactions.
func MinimizeTransactions(balances map[uint]int64) []Transaction {
	var creditors []balance
	var debtors []balance

	for uid, amt := range balances {
		if amt > 0 {
			creditors = append(creditors, balance{uid, amt})
		} else if amt < 0 {
			debtors = append(debtors, balance{uid, -amt}) // store as positive
		}
	}

	// Sort descending by absolute amount
	sort.Slice(creditors, func(i, j int) bool { return creditors[i].Amount > creditors[j].Amount })
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].Amount > debtors[j].Amount })

	var transactions []Transaction
	i, j := 0, 0

	for i < len(debtors) && j < len(creditors) {
		debtor := &debtors[i]
		creditor := &creditors[j]

		// Transfer the smaller of the two amounts
		transfer := debtor.Amount
		if creditor.Amount < transfer {
			transfer = creditor.Amount
		}

		transactions = append(transactions, Transaction{
			From:   debtor.UserID,
			To:     creditor.UserID,
			Amount: transfer,
		})

		debtor.Amount -= transfer
		creditor.Amount -= transfer

		if debtor.Amount == 0 {
			i++
		}
		if creditor.Amount == 0 {
			j++
		}
	}

	return transactions
}
