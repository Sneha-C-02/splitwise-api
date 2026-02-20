package main

import (
	"splitwise-api/config"
	"splitwise-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running 🚀"})
	})

	// ── Phase 1: Auth ──────────────────────────────────────────
	r.POST("/register", handlers.Register)
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id/summary", handlers.GetUserSummary)

	// ── Phase 2: Groups ────────────────────────────────────────
	r.POST("/groups", handlers.CreateGroup)
	r.POST("/groups/:id/members", handlers.AddMember)
	r.GET("/groups/:id", handlers.GetGroup)

	// ── Phase 3: Expenses ──────────────────────────────────────
	r.POST("/groups/:id/expenses", handlers.AddExpense)
	r.GET("/groups/:id/expenses", handlers.GetExpenses)
	r.DELETE("/expenses/:id", handlers.DeleteExpense)

	// ── Phase 4 & 5: Balances & Settlements ────────────────────
	r.GET("/groups/:id/balances", handlers.GetBalances)
	r.GET("/groups/:id/settlements", handlers.GetSettlements)

	r.Run(":8080")
}
