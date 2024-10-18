package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func main() {
	// Adyen'in API anahtarını ve URL'sini ayarla
	const apiKey = ""
	checkoutURL := "https://checkout-test.adyen.com/v68/payments"

	// HTTP istemciyi oluştur
	client := resty.New()

	// Ödeme talebini hazırlıyoruz
	request := map[string]interface{}{
		"merchantAccount": "Sadcar_123456_TEST",
		"amount": map[string]interface{}{
			"currency": "EUR",
			"value":    1000, // 10.00 EUR anlamına gelir
		},
		"reference": "YOUR_ORDER_REFERENCE",
		"paymentMethod": map[string]string{
			"type":                  "scheme",
			"encryptedCardNumber":   "test_4111111111111111", // Test kart numarası
			"encryptedExpiryMonth":  "03",
			"encryptedExpiryYear":   "2030",
			"encryptedSecurityCode": "737",
		},
		"returnUrl": "https://your-return-url.com/",
	}

	// Adyen API'ye ödeme talebi gönder
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-API-Key", apiKey).
		SetBody(request).
		Post(checkoutURL)

	// Hata kontrolü yap
	if err != nil {
		log.Fatalf("API çağrısı sırasında hata: %v", err)
	}

	// Yanıtı yazdır
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Body:", resp.String())
}
