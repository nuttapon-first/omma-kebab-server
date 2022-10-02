package dto

type DashboardInfo struct {
	TotalIncome        float32         `json:"totalIncome"`
	TotalCost          float32         `json:"totalCost"`
	TotalProfit        float32         `json:"totalProfit"`
	TotalProfitPercent float32         `json:"totalProfitPercent"`
	TotalExpense       float32         `json:"totalExpense"`
	IncomeDetails      []IncomeDetails `json:"incomeDetails"`
}

type IncomeDetails struct {
	Channel string  `json:"channel"`
	Income  float32 `json:"income"`
	Cost    float32 `json:"cost"`
	Profit  float32 `json:"profit"`
	Expense float32 `json:"expense"`
}
