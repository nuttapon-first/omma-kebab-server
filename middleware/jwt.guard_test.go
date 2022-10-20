package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testCases struct {
	Name            string
	Body            string
	CodeExpected    int
	MessageExpected string
}

type response struct {
	ErrMessage string `json:"error"`
}

func FakeHandler(c *gin.Context) {}

func TestAuthorizationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []testCases{
		{
			Name:            "Test no permission role",
			Body:            "manager",
			CodeExpected:    http.StatusForbidden,
			MessageExpected: "Permission denied",
		},
	}

	for _, val := range testCases {
		t.Run(val.Name, func(t *testing.T) {
			res := httptest.NewRecorder()
			c, r := gin.CreateTestContext(res)

			jsonStr, _ := json.Marshal(val.Body)

			c.Request = httptest.NewRequest(http.MethodPost, "/menus", bytes.NewBuffer(jsonStr))
			c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

			r.Use(func(c *gin.Context) {
				c.Set("userRole", val.Body)
			})
			r.Use(Authorization("admin"))
			r.POST("/menus", FakeHandler) // Call to a handler method
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
