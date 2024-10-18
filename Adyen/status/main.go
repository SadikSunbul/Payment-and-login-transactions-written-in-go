package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Adyen API URL'si
	url := "https://checkout-test.adyen.com/v68/payments/details"

	// API anahtarınızı buraya ekleyin
	const apiKey = ""

	// İstek gövdesi
	requestBody, err := json.Marshal(map[string]string{
		"pspReference": "L2WW5WPSM3TPBC75",
		"paymentData":  "transactionRiskAnalysis", // Buraya geçerli paymentData'yı ekleyin
	})
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	// HTTP isteği oluşturma
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Başlıkları ayarlama
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// HTTP istemcisi oluşturma ve isteği gönderme
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Yanıtı okuma
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Yanıtı yazdırma
	fmt.Println("Response:", string(body))
}
