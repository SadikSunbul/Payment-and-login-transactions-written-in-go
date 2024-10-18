package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Company struct {
	Name    string
	Email   string
	Phone   string
	Website string
	Logo    string
}

type Customer struct {
	ID    string
	Email string
	Phone string
}

type Product struct {
	Description string
	Rate        float64
	Quantity    float64
}

func generateInvoice(company Company, customer Customer, products []Product, invoiceNo string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set colors
	pdf.SetTextColor(70, 70, 70)
	pdf.SetFillColor(65, 105, 225) // Royal Blue

	// Add logo placeholder
	pdf.Image(company.Logo, 10, 10, 40, 0, false, "", 0, "")

	// Add "Invoice" text
	pdf.SetFont("Helvetica", "B", 30)
	pdf.SetXY(150, 10)
	pdf.Cell(50, 10, "Invoice")

	// Add invoice details
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(150, 25)
	pdf.Cell(30, 10, "Invoice no.:")
	pdf.SetXY(180, 25)
	pdf.Cell(20, 10, invoiceNo)

	pdf.SetXY(150, 31)
	pdf.Cell(30, 10, "Invoice date:")
	pdf.SetXY(180, 31)
	pdf.Cell(20, 10, time.Now().Format("Jan 02, 2006"))

	pdf.SetXY(150, 37)
	pdf.Cell(30, 10, "Due:")
	pdf.SetXY(180, 37)
	pdf.Cell(20, 10, "on receipt")

	// Add company details
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(10, 50)
	pdf.Cell(50, 10, "From")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(10, 60)
	pdf.Cell(50, 10, company.Name)
	pdf.SetXY(10, 65)
	pdf.Cell(50, 10, company.Email)
	pdf.SetXY(10, 70)
	pdf.Cell(50, 10, company.Phone)
	pdf.SetXY(10, 75)
	pdf.Cell(50, 10, company.Website)

	// Add customer details
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetXY(150, 50)
	pdf.Cell(50, 10, "Bill to")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(150, 60)
	pdf.Cell(50, 10, customer.ID)
	pdf.SetXY(150, 65)
	pdf.Cell(50, 10, customer.Email)
	pdf.SetXY(150, 70)
	pdf.Cell(50, 10, customer.Phone)

	// Add product table
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(65, 105, 225)
	pdf.SetTextColor(255, 255, 255)
	pdf.Rect(10, 90, 190, 10, "F")
	pdf.SetXY(10, 90)
	pdf.Cell(80, 10, "DESCRIPTION")
	pdf.SetXY(90, 90)
	pdf.Cell(30, 10, "RATE, USD")
	pdf.SetXY(120, 90)
	pdf.Cell(30, 10, "QTY")
	pdf.SetXY(150, 90)
	pdf.Cell(50, 10, "AMOUNT, USD")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(70, 70, 70)
	var subtotal float64
	for i, product := range products {
		y := 100 + float64(i*10)
		amount := product.Rate * product.Quantity
		subtotal += amount

		if i%2 == 1 {
			pdf.SetFillColor(240, 248, 255) // Light blue
			pdf.Rect(10, y, 190, 10, "F")
		}

		pdf.SetXY(10, y)
		pdf.Cell(80, 10, product.Description)
		pdf.SetXY(90, y)
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", product.Rate))
		pdf.SetXY(120, y)
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", product.Quantity))
		pdf.SetXY(150, y)
		pdf.Cell(50, 10, fmt.Sprintf("%.2f", amount))
	}

	// Add totals
	y := 120 + float64(len(products)*10)
	tax := subtotal * 0.1145 // 11.45% tax
	total := subtotal + tax

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(120, y)
	pdf.Cell(30, 10, "Subtotal:")
	pdf.SetXY(150, y)
	pdf.Cell(50, 10, fmt.Sprintf("$%.2f", subtotal))

	pdf.SetXY(120, y+10)
	pdf.Cell(30, 10, "Tax in items:")
	pdf.SetXY(150, y+10)
	pdf.Cell(50, 10, fmt.Sprintf("$%.2f", tax))

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetXY(120, y+20)
	pdf.Cell(30, 10, "Total:")
	pdf.SetXY(150, y+20)
	pdf.Cell(50, 10, fmt.Sprintf("$%.2f", total))

	pdf.SetFillColor(240, 248, 255) // Light blue
	pdf.Rect(120, y+30, 80, 10, "F")
	pdf.SetXY(120, y+30)
	pdf.Cell(30, 10, "Balance Due:")
	pdf.SetXY(150, y+30)
	pdf.Cell(50, 10, fmt.Sprintf("$%.2f", total))

	// Save the PDF
	return pdf.OutputFileAndClose("invoice.pdf")
}

func main() {
	company := Company{
		Name:    "log-sysdev.de",
		Email:   "test@gmail.com",
		Phone:   "+2342455464567",
		Website: "https://log-sysdev.de/about",
		Logo:    "logo.png",
	}

	customer := Customer{
		ID:    "000000000000001",
		Email: "customer@asfjsaf.com",
		Phone: "+34345475675",
	}

	products := []Product{
		{Description: "Doner", Rate: 5.00, Quantity: 1.00},
		{Description: "Ayran", Rate: 2.00, Quantity: 2.00},
	}

	err := generateInvoice(company, customer, products, "001")
	if err != nil {
		fmt.Printf("Error generating invoice: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Invoice generated successfully!")
}
