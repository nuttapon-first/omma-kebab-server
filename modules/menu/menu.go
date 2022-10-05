package menu

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/router"
	"github.com/nuttapon-first/omma-kebab-server/store"
)

// ////////////////////////////////////////////////////////////////////
// SPI

type MenuHandler struct {
	store store.Storer
}

func NewMenuHandler(store store.Storer) *MenuHandler {
	return &MenuHandler{store: store}
}

// ////////////////////////////////////////////////////////////////////

func (m *MenuHandler) NewMenu(c router.Context) {
	menu := &model.Menu{}

	if err := c.ShouldBindJSON(menu); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if menu.MenuCost <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "cost should greater than zero",
		})
		return
	}

	if menu.MenuUnit <= 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "unit should greater than zero",
		})
		return
	}

	err := m.store.New(menu)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"ID": menu.ID,
	})
}

func (m *MenuHandler) GetMenuList(c router.Context) {
	menus := &[]model.Menu{}
	err := m.store.Find(menus, &model.Menu{}, "")
	if err != nil {
		fmt.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
		"menus":   menus,
	})
}

func (m *MenuHandler) GetMenuById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	menu := &model.Menu{}
	result := m.store.First(menu, id, "MenuRecipe")
	if err := result; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
		"menu":    menu,
	})
}

func (m *MenuHandler) EditById(c router.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	idPayload := map[string]interface{}{"ID": id}
	menu := &model.EditMenu{}
	if err := c.ShouldBindJSON(menu); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tx := m.store.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	r := tx.Model(&model.Menu{}).Where(idPayload).Updates(menu)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// recipeIds := []string{}
	for _, val := range menu.MenuRecipe {
		if val.ID != 0 {
			r := tx.Model(&model.Recipe{}).Where(map[string]interface{}{"ID": val.ID}).Updates(val)
			if err := r.Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
		} else {
			val.MenuID = id
			r := tx.Model(&model.Recipe{}).Create(&val)
			if err := r.Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
		}
	}

	if len(menu.DeletedRecipeIds) > 0 {
		r := tx.Model(&model.Recipe{}).Where(map[string]interface{}{"ID": menu.DeletedRecipeIds}).Delete(&model.Recipe{})
		if err := r.Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
	})
}

func (m *MenuHandler) RemoveById(c router.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tx := m.store.Begin()

	r := tx.Where("menu_id = ?", id).Delete(&model.Recipe{})
	if err := r.Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	r = tx.Delete(&model.Menu{}, id)
	if err := r.Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
	})
}
