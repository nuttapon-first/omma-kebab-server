package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nuttapon-first/omma-kebab-server/middleware"
	"github.com/nuttapon-first/omma-kebab-server/migration"
	"github.com/nuttapon-first/omma-kebab-server/modules/expense"
	"github.com/nuttapon-first/omma-kebab-server/modules/login"
	"github.com/nuttapon-first/omma-kebab-server/modules/menu"
	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/modules/report"
	"github.com/nuttapon-first/omma-kebab-server/modules/stock"
	"github.com/nuttapon-first/omma-kebab-server/modules/transaction"
	"github.com/nuttapon-first/omma-kebab-server/router"
	"github.com/nuttapon-first/omma-kebab-server/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}
}

func main() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.Stock{}, &model.Transaction{}, &model.Expense{}, &model.User{}, &model.UserCredential{})

	err = migration.CheckAdminUser(db)
	if err != nil {
		fmt.Printf("Migrate admin error : %s\n", err.Error())
	}

	r := router.NewRouter()
	gormStore := store.NewGormStore(db)

	loginHandler := login.NewLoginHandler(gormStore)
	r.POST("/login", loginHandler.Login)

	v1 := r.Group("/v1")
	v1.Use(middleware.JwtGuard())
	{
		menuHandler := menu.NewMenuHandler(gormStore)
		v1.POST("menus", menuHandler.NewMenu)
		v1.GET("menus", menuHandler.GetMenuList)
		v1.GET("menus/:id", menuHandler.GetMenuById)
		v1.PUT("menus/:id", menuHandler.EditById)
		v1.DELETE("menus/:id", menuHandler.RemoveById)

		transactionHandler := transaction.NewTransactionHandler(gormStore)
		v1.POST("transactions", transactionHandler.New)
		v1.GET("transactions", transactionHandler.GetList)
		v1.GET("transactions/:id", transactionHandler.GetById)
		v1.DELETE("transactions/:id", transactionHandler.RemoveById)

		stockHandler := stock.NewStockHandler(gormStore)
		v1.POST("stocks", stockHandler.New)
		v1.GET("stocks", stockHandler.GetList)
		v1.GET("stocks/:id", stockHandler.GetById)
		v1.PUT("stocks/:id", stockHandler.EditById)
		v1.PUT("stocks/:id/add", stockHandler.AddById)
		v1.DELETE("stocks/:id", stockHandler.RemoveById)

		expenseHandler := expense.NewExpenseHandler(gormStore)
		v1.POST("expenses", expenseHandler.New)
		v1.GET("expenses", expenseHandler.GetList)
		// v1.DELETE("expenses/:id", router.NewGinHandler(expenseHandler.RemoveById))

		reportHandler := report.NewReportHandler(gormStore)
		v1.GET("reports/dashboard", reportHandler.GetDashboard)

	}

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
