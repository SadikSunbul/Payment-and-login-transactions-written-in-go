package main

import (
	"fmt"
	"os"
	"path/filepath"
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

type Courier struct {
	Name  string
	Phone string
	Mail  string
}

type User struct {
	Name    string
	Email   string
	Phone   string
	Address string
}

type OrderItem struct {
	Name     string
	Quantity int
	Price    float64
}

type Order struct {
	Items             []OrderItem
	Total             float64
	Discount          float64
	FinalAmount       float64
	OrderTime         time.Time
	DeliveryTime      time.Time
	Restaurant        string
	RestaurantAddress string
	PaymentMethod     string
	TaxPrice          float64
	OrderCode         string
}

func generateReceipt(company Company, courier Courier, user User, order Order, path string) error {
	baseHeight := 145.0
	itemHeight := 6.0
	pageHeight := baseHeight + itemHeight*float64(len(order.Items))

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: 74,
			Ht: pageHeight,
		},
	})
	pdf.SetMargins(1.5, 1.5, 1.5)
	pdf.SetFont("Helvetica", "", 8)
	pdf.AddPage()

	// Şirket bilgilerini ekle
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(74, 10, company.Name, "", 0, "C", false, 0, "")
	pdf.Ln(12)

	// Kurye bilgilerini ekle
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(74, 5, fmt.Sprintf("Courier: %s, Phone: %s", courier.Name, courier.Phone), "", 0, "L", false, 0, "")
	pdf.Ln(5)

	// Kullanıcı bilgilerini ekle
	pdf.CellFormat(74, 5, fmt.Sprintf("Customer: %s, Phone: %s", user.Name, user.Phone), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(74, 5, fmt.Sprintf("Address: %s", user.Address), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Çizgi çiz
	pdf.Line(1.5, pdf.GetY(), 72.5, pdf.GetY())
	pdf.Ln(3)

	// Sipariş detayları
	pdf.CellFormat(74, 5, fmt.Sprintf("Ordered on %s", order.OrderTime.Format("Jan 02, 15:04 PM")), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(74, 5, fmt.Sprintf("Due at %s", order.DeliveryTime.Format("Jan 02, 15:04 PM")), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Çizgi çiz
	pdf.Line(1.5, pdf.GetY(), 72.5, pdf.GetY())
	pdf.Ln(3)

	// Teslimat bölümü
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(74, 5, "DELIVERY", "", 0, "L", false, 0, "")
	pdf.Ln(7)

	// Sipariş kodunu ekle
	pdf.SetFont("Helvetica", "U", 8)
	pdf.CellFormat(54, 2, fmt.Sprintf("Order Code: %s", order.OrderCode), "", 0, "L", false, 0, "")
	pdf.Ln(3)

	// Öğeler ve vergi miktarları
	pdf.SetFont("Helvetica", "", 8)
	for _, item := range order.Items {
		itemTotal := item.Price * float64(item.Quantity)
		pdf.CellFormat(40, 5, fmt.Sprintf("%dx %s", item.Quantity, item.Name), "", 0, "L", false, 0, "")
		pdf.CellFormat(24, 5, fmt.Sprintf("%.2f EUR", itemTotal), "", 0, "R", false, 0, "")
		pdf.Ln(5)
	}

	// Ara toplam ve indirim
	pdf.Ln(5)
	pdf.CellFormat(40, 5, "Subtotal", "", 0, "L", false, 0, "")
	pdf.CellFormat(24, 5, fmt.Sprintf("%.2f EUR", order.Total), "", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(40, 5, "Discount", "", 0, "L", false, 0, "")
	pdf.CellFormat(24, 5, fmt.Sprintf("(%.2f EUR)", order.Discount), "", 0, "R", false, 0, "")
	pdf.Ln(5)

	// Toplam vergi
	pdf.CellFormat(40, 5, "Total Tax", "", 0, "L", false, 0, "")
	pdf.CellFormat(24, 5, fmt.Sprintf("%.2f EUR", order.TaxPrice), "", 0, "R", false, 0, "")
	pdf.Ln(5)

	// Ödenen miktar
	pdf.SetFont("Helvetica", "B", 8)
	pdf.CellFormat(40, 5, "Amount Paid", "", 0, "L", false, 0, "")
	pdf.CellFormat(24, 5, fmt.Sprintf("%.2f EUR", order.FinalAmount), "", 0, "R", false, 0, "")
	pdf.Ln(10)

	// Çizgi çiz
	pdf.Line(1.5, pdf.GetY(), 72.5, pdf.GetY())
	pdf.Ln(3)

	// Ödeme yöntemi
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(74, 5, fmt.Sprintf("Payment Method: %s", order.PaymentMethod), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Restoran bilgileri
	pdf.CellFormat(74, 5, fmt.Sprintf("Restaurant: %s", order.Restaurant), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(74, 5, fmt.Sprintf("Address: %s", order.RestaurantAddress), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Dosya oluşturma ve yazma
	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return fmt.Errorf("dosya oluşturma ve yazma hatası: %w", err)
	}

	return nil
}

func main() {
	company := Company{
		Name:    "Karadeniz",
		Email:   "support@karadeniz.com",
		Phone:   "+123456789",
		Website: "www.karadeniz.com",
		Logo:    "logo.png",
	}

	courier := Courier{
		Name:  "John Doe",
		Phone: "+987654321",
	}

	user := User{
		Name:    "Jane Smith",
		Email:   "jane.smith@example.com",
		Phone:   "+1122334455",
		Address: "123 Main St, City, Country",
	}

	order := Order{
		Items: []OrderItem{
			{Name: "Doner", Quantity: 2, Price: 19.90},
			{Name: "Ayran", Quantity: 1, Price: 2.00},
			{Name: "Kola", Quantity: 3, Price: 3.00},
			{Name: "ssagsdg", Quantity: 1, Price: 5.00},
			{Name: "sdgsdg", Quantity: 5, Price: 1.00},
		},
		Total:             59.80,
		Discount:          0.00,
		FinalAmount:       39.90,
		OrderTime:         time.Now(),
		DeliveryTime:      time.Now().Add(5 * time.Minute),
		Restaurant:        "Karadeniz",
		RestaurantAddress: "456 Another St, City, Country",
		PaymentMethod:     "Paypal",
		TaxPrice:          10.00,
		OrderCode:         "ORD123456",
	}

	outputPath := "deneme/logo/receipt.pdf"

	// Dizin oluşturma
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Dizin oluşturma hatası: %v\n", err)
		return
	}

	// Dosya oluşturma veya üzerine yazma
	err := generateReceipt(company, courier, user, order, outputPath)
	if err != nil {
		fmt.Printf("Makbuz oluşturma hatası: %v\n", err)
		return
	}

	fmt.Printf("Makbuz başarıyla oluşturuldu: %s\n", outputPath)
}
