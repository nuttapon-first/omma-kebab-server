package expense

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nuttapon-first/omma-kebab-server/modules/dto"
	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/modules/pkg"
	"github.com/nuttapon-first/omma-kebab-server/router"
	"github.com/nuttapon-first/omma-kebab-server/store"
)

// ////////////////////////////////////////////////////////////////////
// SPI

type ExpenseHandler struct {
	store store.Storer
}

func NewExpenseHandler(store store.Storer) *ExpenseHandler {
	return &ExpenseHandler{store: store}
}

// ////////////////////////////////////////////////////////////////////

func (h *ExpenseHandler) New(c router.Context) {
	payload := &dto.ExpenseDto{}

	if err := c.ShouldBindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	createTime := time.Unix(0, int64(payload.CreatedAt)*int64(time.Millisecond))

	expense := &model.Expense{
		Branch:        payload.Branch,
		ExpenseDetail: payload.ExpenseDetail,
		ExpenseCost:   payload.ExpenseCost,
		CreatedAt:     createTime,
	}

	if expense.ExpenseCost <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "expense cost should greater than zero",
		})
		return
	}

	err := h.store.New(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": 0,
		"ID":      expense.ID,
	})
}

func (h *ExpenseHandler) GetList(c router.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	timeFormat := "2006-01-02 15:04:05"
	startDate, err := pkg.FormatDateQuery(timeFormat, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	endDate, err = pkg.FormatDateQuery(timeFormat, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	where := fmt.Sprintf("Expenses.created_at BETWEEN '%s' AND '%s'", startDate, endDate)

	expenses := &[]dto.ExpenseList{}
	err = h.store.Table("Expenses").Order("created_at asc").Where(where).Find(expenses).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	expenseReport := &dto.ExpenseReport{
		Transactions: expenses,
	}

	for _, transaction := range *expenses {
		expenseReport.SummaryCost += transaction.ExpenseCost
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":       0,
		"expenseReport": expenseReport,
	})
}

// func (h *ExpenseHandler) RemoveById(c router.Context) {
// 	idParam := c.Param("id")

// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	r := h.store.Delete(&model.Expense{}, id)
// 	if err := r; err != nil {
// 		c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"success": 0,
// 	})
// }
