package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
)

type Order struct {
	OrderID     string    `json:"orderId"`
	OrderCode   string    `json:"orderCode"`
	TotalPrice  float64   `json:"totalPrice"`
	TotalTax    float64   `json:"totalTax"`
	Date        time.Time `json:"date"`
	CompanyName string    `json:"companyName"`
	IsCanceled  bool      `json:"isCanceled"`
}

func main() {
	http.HandleFunc("/generate-excel", generateExcelHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func generateExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)
		return
	}

	var orders []Order
	err := json.NewDecoder(r.Body).Decode(&orders)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	f := excelize.NewFile()

	// Normal Orders Table
	createOrderSheet(f, "Normal Orders", filterOrders(orders, false))

	// Canceled Orders Table
	createOrderSheet(f, "Canceled Orders", filterOrders(orders, true))

	// Summary Table
	createSummarySheet(f, orders)

	// Send Excel file as response
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=orders.xlsx")
	err = f.Write(w)
	if err != nil {
		http.Error(w, "Error occurred while creating Excel file", http.StatusInternalServerError)
		return
	}
}

func createOrderSheet(f *excelize.File, sheetName string, orders []Order) {
	f.NewSheet(sheetName)

	headers := []string{"Order ID", "Order Code", "Total Price", "Total Tax", "Date", "Company Name", "Is Canceled"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, order := range orders {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), order.OrderID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), order.OrderCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), order.TotalPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), order.TotalTax)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), order.Date.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), order.CompanyName)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), order.IsCanceled)
	}
}

func createSummarySheet(f *excelize.File, orders []Order) {
	sheetName := "Summary"
	f.NewSheet(sheetName)

	var totalRevenue, totalTax float64
	var normalOrderCount, canceledOrderCount int

	for _, order := range orders {
		if !order.IsCanceled {
			totalRevenue += order.TotalPrice
			totalTax += order.TotalTax
			normalOrderCount++
		} else {
			canceledOrderCount++
		}
	}

	f.SetCellValue(sheetName, "A1", "Total Revenue")
	f.SetCellValue(sheetName, "B1", totalRevenue)
	f.SetCellValue(sheetName, "A2", "Total Tax")
	f.SetCellValue(sheetName, "B2", totalTax)
	f.SetCellValue(sheetName, "A3", "Normal Order Count")
	f.SetCellValue(sheetName, "B3", normalOrderCount)
	f.SetCellValue(sheetName, "A4", "Canceled Order Count")
	f.SetCellValue(sheetName, "B4", canceledOrderCount)
}

func filterOrders(orders []Order, isCanceled bool) []Order {
	var filteredOrders []Order
	for _, order := range orders {
		if order.IsCanceled == isCanceled {
			filteredOrders = append(filteredOrders, order)
		}
	}
	return filteredOrders
}
