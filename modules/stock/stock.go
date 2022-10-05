package stock

import (
	"net/http"
	"strconv"

	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/router"
	"github.com/nuttapon-first/omma-kebab-server/store"
	"gorm.io/gorm"
)

// ////////////////////////////////////////////////////////////////////
// SPI

type StockHandler struct {
	store store.Storer
}

func NewStockHandler(store store.Storer) *StockHandler {
	return &StockHandler{store: store}
}

// ////////////////////////////////////////////////////////////////////

func (h *StockHandler) New(c router.Context) {
	stock := &model.Stock{}

	if err := c.ShouldBindJSON(stock); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if stock.StockCost <= 0 || stock.StockAmount <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "cost or amount should greater than zero",
		})
		return
	}

	err := h.store.New(stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": 0,
		"ID":      stock.ID,
	})
}

func (h *StockHandler) GetList(c router.Context) {
	stocks := &[]model.Stock{}
	err := h.store.Find(stocks, &model.Stock{}, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
		"stocks":  stocks,
	})
}

func (h *StockHandler) GetById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	stock := &model.Stock{}
	result := h.store.First(stock, id, "")
	if err := result; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
		"stock":   stock,
	})
}

func (h *StockHandler) AddById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	where := map[string]interface{}{"ID": id}
	stock := &model.AddStock{}

	if err := c.ShouldBindJSON(stock); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	r := h.store.Updates(where, &model.Stock{}, map[string]interface{}{"stock_amount": gorm.Expr("stock_amount + ?", stock.StockAmount)})
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

func (h *StockHandler) EditById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	idPayload := map[string]interface{}{"ID": id}
	stock := &model.Stock{}

	if err := c.ShouldBindJSON(stock); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	r := h.store.Updates(idPayload, stock, stock)
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

func (h *StockHandler) RemoveById(c router.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	r := h.store.Delete(&model.Stock{}, id)
	if err := r; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err = h.store.Table("Recipes").Where("stock_id = ?", id).Delete(&model.Recipe{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
	})
}
