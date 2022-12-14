package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nuttapon-first/omma-kebab-server/modules/dto"
	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/modules/pkg"
	"github.com/nuttapon-first/omma-kebab-server/router"
	"github.com/nuttapon-first/omma-kebab-server/store"
	"gorm.io/gorm"
)

// ////////////////////////////////////////////////////////////////////
// SPI

type TransactionHandler struct {
	store store.Storer
}

func NewTransactionHandler(store store.Storer) *TransactionHandler {
	return &TransactionHandler{store: store}
}

// ////////////////////////////////////////////////////////////////////

func (h *TransactionHandler) New(c router.Context) {
	payload := &dto.TransactionDto{}

	if err := c.ShouldBindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// convert unix millisecond to time
	createTime := time.Unix(0, int64(payload.CreatedAt)*int64(time.Millisecond))

	transaction := &model.Transaction{
		MenuID:           payload.MenuID,
		CreatedAt:        createTime,
		Branch:           payload.Branch,
		TransactionType:  payload.TransactionType,
		Channel:          payload.Channel,
		TransactionPrice: payload.TransactionPrice,
		TransactionUnit:  payload.TransactionUnit,
		Discount:         payload.Discount,
		PaymentChannel:   payload.PaymentChannel,
	}

	if transaction.TransactionUnit <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "unit should greater than zero",
		})
		return
	}

	if transaction.TransactionPrice <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid type and price",
		})
		return
	}

	menu := &model.Menu{}
	addOn := []string{}
	result := h.store.First(menu, transaction.MenuID, "MenuRecipe")
	if err := result; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Invalid menu id",
			})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return
	}

	tx := h.store.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, val := range menu.MenuRecipe {
		// where := map[string]interface{}{"ingredient": []string{"box", "box lid", "kebab wrapper", "tortillas flour"}}
		r := tx.Model(&model.Stock{}).Where(map[string]interface{}{"ID": val.StockID}).Updates(map[string]interface{}{"stock_amount": gorm.Expr("stock_amount - ?", val.IngredientAmount)})
		if err := r.Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	}

	var addOnCost float32
	if len(payload.AddOns) > 0 {
		for _, val := range payload.AddOns {
			stock := &model.Stock{
				ID: val.StockId,
			}
			r := tx.Table("Stocks").Where(stock).First(stock)
			if err := r.Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			stock.StockAmount = stock.StockAmount - val.Amount
			addOnCost += (stock.StockCost * val.Amount)
			r = tx.Save(stock)
			if err := r.Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// r := tx.Model(&model.Stock{}).Where(map[string]interface{}{"ingredient": val.Name}).Updates(map[string]interface{}{"stock_amount": gorm.Expr("stock_amount - ?", val.Amount)})
			// fmt.Printf("%#v\n", r)
			// if err := r.Error; err != nil {
			// 	tx.Rollback()
			// 	c.JSON(http.StatusInternalServerError, map[string]interface{}{
			// 		"error": err.Error(),
			// 	})
			// 	return
			// }
			addOn = append(addOn, stock.Ingredient)
		}
		transaction.AddOn = strings.Join(addOn, ", ")
	}

	transaction.TotalPrice = float32(transaction.TransactionPrice*transaction.TransactionUnit) - transaction.Discount
	transaction.TotalCost = menu.MenuCost*float32(transaction.TransactionUnit) + addOnCost

	if transaction.Channel == "lineman_30" {
		transaction.Fee = transaction.TotalPrice * 0.3
		transaction.Vat = transaction.Fee * 0.07
	} else if transaction.Channel == "lineman_09" {
		transaction.Fee = transaction.TotalPrice * 0.09
		transaction.Vat = transaction.Fee * 0.07
	}

	transaction.TotalPrice = transaction.TotalPrice - transaction.Vat - transaction.Fee
	transaction.TotalProfit = transaction.TotalPrice - transaction.TotalCost
	transaction.TotalProfitPercent = (transaction.TotalProfit / transaction.TotalPrice) * 100

	fmt.Printf("%#v\n", transaction)

	err := tx.Create(transaction)
	if err := err.Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": 0,
		"ID":      transaction.ID,
	})
}

func (h *TransactionHandler) GetList(c router.Context) {
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

	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	sortQuery := c.Query("sort")
	sort := ""
	if sortQuery != "" {
		direction := c.Query("sortDesc")
		if direction != "" {
			if direction == "true" {
				sort = sortQuery + " desc"
			} else if direction == "false" {
				sort = sortQuery + " asc"
			}
		}

	}

	switch {
	case pageSize <= 0:
		pageSize = 30
	}

	pagination := &pkg.Pagination{
		Page:  page,
		Limit: pageSize,
	}

	if sort != "" {
		pagination.Sort = sort
	} else {
		pagination.Sort = "created_at desc"
	}

	where := fmt.Sprintf("Transactions.created_at BETWEEN '%s' AND '%s'", startDate, endDate)
	selectRow := "`Transactions`.`id`,`Transactions`.`created_at`,`Transactions`.`updated_at`,`Transactions`.`menu_id`,`Transactions`.`branch`,`Transactions`.`transaction_type`,`Transactions`.`channel`,`Transactions`.`transaction_price`,`Transactions`.`transaction_unit`,`Transactions`.`fee`,`Transactions`.`vat`,`Transactions`.`discount`,`Transactions`.`total_price`,`Transactions`.`total_cost`,`Transactions`.`total_profit`,`Transactions`.`total_profit_percent`,`Transactions`.`payment_channel`,`Transactions`.`add_on`, `MenuList`.`menu_name`"
	rows, err := h.store.Table("Transactions").Where(where).Scopes(h.store.Paginate(&model.Transaction{}, pagination, h.store.Table("Transactions"), where)).Select(selectRow).Joins("join MenuList on Transactions.menu_id = MenuList.id").Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()
	transactionList := &[]dto.TransactionList{}

	for rows.Next() {
		h.store.ScanRows(rows, transactionList)
	}

	pagination.Rows = *transactionList

	salesReport := &dto.SalesReport{
		Transactions: transactionList,
	}

	for _, transaction := range *transactionList {
		salesReport.SummaryAmount += transaction.TransactionUnit
		salesReport.SummaryFee += transaction.Fee
		salesReport.SummaryVat += transaction.Vat
		salesReport.SummarySales += transaction.TotalPrice
		salesReport.SummaryCost += transaction.TotalCost
		salesReport.SummaryProfit += transaction.TotalProfit
		salesReport.AvgProfitPercent += transaction.TotalProfitPercent
	}

	if salesReport.AvgProfitPercent > 0 {
		salesReport.AvgProfitPercent = salesReport.AvgProfitPercent / float32(len(*transactionList))
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":      0,
		"transactions": pagination,
	})
}

func (h *TransactionHandler) GetById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	transaction := &model.Transaction{}
	result := h.store.First(transaction, id, "")
	if err := result; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":     0,
		"transaction": transaction,
	})
}

func (h *TransactionHandler) RemoveById(c router.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	r := h.store.Delete(&model.Transaction{}, id)
	if err := r; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
	})
}
