package report

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuttapon-first/omma-kebab-server/modules/dto"
	"github.com/nuttapon-first/omma-kebab-server/store"
)

type ReportHandler struct {
	store store.Storer
}

func NewReportHandler(store store.Storer) *ReportHandler {
	return &ReportHandler{store: store}
}

func (h *ReportHandler) GetDashboard(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	t := time.Now()
	timeFormat := "2006-01-02 15:04:05"
	if startDate == "" {
		year, month, day := t.Date()
		startDate = time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Format(timeFormat)
	} else {
		start, err := time.Parse("20060102150405", startDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		startDate = start.Format(timeFormat)
	}

	if endDate == "" {
		year, month, day := t.Date()
		endDate = time.Date(year, month, day, 23, 59, 59, 59, t.Location()).Format(timeFormat)
	} else {
		end, err := time.Parse("20060102150405", endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		endDate = end.Format(timeFormat)
	}

	transactionList := &[]dto.TransactionList{}
	where := fmt.Sprintf("Transactions.created_at BETWEEN '%s' AND '%s'", startDate, endDate)
	err := h.store.Table("Transactions").Order("created_at asc").Where(where).Find(transactionList).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	where = fmt.Sprintf("Expenses.created_at BETWEEN '%s' AND '%s'", startDate, endDate)
	expenseList := &[]dto.ExpenseList{}
	err = h.store.Table("Expenses").Order("created_at asc").Where(where).Find(expenseList).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	dashboard := &dto.DashboardInfo{}
	linemanChannel := &dto.IncomeDetails{
		Channel: "lineman",
	}
	frontendChannel := &dto.IncomeDetails{
		Channel: "frontend",
	}

	for _, transaction := range *transactionList {
		dashboard.TotalIncome += float32(transaction.TotalPrice)
		dashboard.TotalCost += transaction.TotalCost
		dashboard.TotalProfit += transaction.TotalProfit

		if transaction.Channel == "frontend" {
			frontendChannel.Income += float32(transaction.TotalPrice)
			frontendChannel.Cost += transaction.TotalCost
			frontendChannel.Profit += transaction.TotalProfit
		} else if strings.Contains(transaction.Channel, "lineman") {
			linemanChannel.Income += float32(transaction.TotalPrice)
			linemanChannel.Cost += transaction.TotalCost
			linemanChannel.Profit += transaction.TotalProfit
		}
	}

	if dashboard.TotalIncome > 0 {
		dashboard.TotalProfitPercent = (dashboard.TotalProfit / dashboard.TotalIncome) * 100
	} else {
		dashboard.TotalProfitPercent = 0
	}

	dashboard.IncomeDetails = append(dashboard.IncomeDetails, *linemanChannel, *frontendChannel)

	for _, expense := range *expenseList {
		dashboard.TotalExpense += float32(expense.ExpenseCost)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":   0,
		"dashboard": dashboard,
	})
}
