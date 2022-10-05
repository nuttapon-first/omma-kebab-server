package stock

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nuttapon-first/omma-kebab-server/router"
)

func NewTestContext(c *gin.Context) *router.Context {
	return &router.Context{Context: c}
}

func NewTestHandler(handler func(router.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(*NewTestContext(c))
	}
}

type response struct {
	ErrMessage string `json:"error"`
}

type testCases struct {
	Name            string
	Body            map[string]interface{}
	CodeExpected    int
	MessageExpected string
}

func TestCreateStockFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []testCases{
		{
			Name: "Test cost is less than zero",
			Body: map[string]interface{}{
				"ingredient":  "test",
				"stockCost":   -1,
				"stockAmount": 1,
			},
			CodeExpected:    http.StatusBadRequest,
			MessageExpected: "cost or amount should greater than zero",
		},
		{
			Name: "Test amount is less than zero",
			Body: map[string]interface{}{
				"ingredient":  "box",
				"stockCost":   1,
				"stockAmount": -1,
			},
			CodeExpected:    http.StatusBadRequest,
			MessageExpected: "cost or amount should greater than zero",
		},
	}

	for _, val := range testCases {
		t.Run(val.Name, func(t *testing.T) {
			res := httptest.NewRecorder()
			c, r := gin.CreateTestContext(res)

			jsonStr, _ := json.Marshal(val.Body)

			c.Request = httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewBuffer(jsonStr))
			c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

			stockHandler := &StockHandler{}
			r.POST("/stocks", NewTestHandler(stockHandler.New)) // Call to a handler method
			r.ServeHTTP(res, c.Request)

			if status := res.Code; status != val.CodeExpected {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			response := &response{}
			json.Unmarshal(res.Body.Bytes(), response)

			if got := response.ErrMessage; got != val.MessageExpected {
				t.Errorf("Message: got %v want %v", got, val.MessageExpected)
			}
		})
	}
}
