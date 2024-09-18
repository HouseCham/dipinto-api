package dto

type AdminDashboardDTO struct {
	Cards      *AdminCardsDTO     `json:"cards"`
	SalesChart *[]MonthlySalesDTO `json:"sales_chart"`
}

type AdminCardsDTO struct {
	TotalCustomers     uint64  `json:"total_customers"`
	TotalSales         float64 `json:"total_sales"`
	TotalExpenses      float64 `json:"total_expenses"`
	TotalPendingOrders uint64  `json:"total_pending_orders"`
}

type MonthlySalesDTO struct {
	Month      string  `json:"month"`
	TotalSales float64 `json:"total_sales"`
}