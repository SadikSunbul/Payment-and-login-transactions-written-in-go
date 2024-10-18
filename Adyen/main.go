package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const apiKey = ""
const adyenURL = "https://checkout-test.adyen.com/v68/payments"
const paymentMethodsURL = "https://checkout-test.adyen.com/v68/paymentMethods"

// Amount yapı
type Amount struct {
	Currency string `json:"currency"`
	Value    int    `json:"value"`
}

// PaymentMethod yapı
type PaymentMethod struct {
	Type                  string `json:"type"`
	EncryptedCardNumber   string `json:"encryptedCardNumber"`
	EncryptedExpiryMonth  string `json:"encryptedExpiryMonth"`
	EncryptedExpiryYear   string `json:"encryptedExpiryYear"`
	EncryptedSecurityCode string `json:"encryptedSecurityCode"`
}

// PaymentRequest ödeme isteği için kullanılan yapı
type PaymentRequest struct {
	Amount          Amount        `json:"amount"`
	Reference       string        `json:"reference"`
	PaymentMethod   PaymentMethod `json:"paymentMethod"`
	ReturnUrl       string        `json:"returnUrl"`
	MerchantAccount string        `json:"merchantAccount"`
}

// PaymentResponse ödeme yanıtı için kullanılan yapı
type PaymentResponse struct {
	PspReference  string `json:"pspReference"`
	ResultCode    string `json:"resultCode"`
	RefusalReason string `json:"refusalReason"`
}

// PaymentMethodsRequest ödeme yöntemleri isteği için kullanılan yapı
type PaymentMethodsRequest struct {
	MerchantAccount string `json:"merchantAccount"`
}

// PaymentMethodsResponse ödeme yöntemleri yanıtı için kullanılan yapı
type PaymentMethodsResponse struct {
	PaymentMethods []interface{} `json:"paymentMethods"`
}

// getPaymentMethods ödeme yöntemlerini alma işlemi
func getPaymentMethods() (*PaymentMethodsResponse, error) {
	request := PaymentMethodsRequest{
		MerchantAccount: "Sadcar_123456_TEST",
	}

	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", paymentMethodsURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentMethodsResponse PaymentMethodsResponse
	err = json.Unmarshal(body, &paymentMethodsResponse)
	if err != nil {
		return nil, err
	}

	return &paymentMethodsResponse, nil
}

// createPayment ödeme oluşturma işlemi
func createPayment(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	var paymentReq PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		// Adyen API'ye istek gönderme
		encryptedCardNumber := "test_4111111111111111" // Bu, şifrelenmiş kart numarasıdır
		encryptedExpiryMonth := "test_03"              // Bu, şifrelenmiş son kullanma ayıdır
		encryptedExpiryYear := "test_2030"             // Bu, şifrelenmiş son kullanma yılıdır
		encryptedSecurityCode := "test_737"
		// PaymentRequest yapısına şifrelenmiş kart bilgilerini ekleyin
		paymentReq.PaymentMethod.EncryptedCardNumber = encryptedCardNumber
		paymentReq.PaymentMethod.EncryptedExpiryMonth = encryptedExpiryMonth
		paymentReq.PaymentMethod.EncryptedExpiryYear = encryptedExpiryYear
		paymentReq.PaymentMethod.EncryptedSecurityCode = encryptedSecurityCode */

	adyenReq, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest("POST", adyenURL, bytes.NewBuffer(adyenReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var paymentResp PaymentResponse
	err = json.Unmarshal(body, &paymentResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Yanıtı loglama
	log.Printf("Adyen API Yanıtı: %s\n", string(body))

	// Adyen API yanıtını kontrol et
	if paymentResp.ResultCode != "Authorised" {
		http.Error(w, paymentResp.RefusalReason, http.StatusPaymentRequired)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentResp)
}

// main fonksiyonu HTTP sunucusunu başlatır ve endpoint'leri dinler
func main() {
	http.HandleFunc("/create-payment", createPayment)
	http.HandleFunc("/payment-methods", func(w http.ResponseWriter, r *http.Request) {
		/*enableCors(&w)
		// OPTIONS isteği için CORS ön uçuş kontrolü
		// Bu, tarayıcının CORS politikasını doğrulamasına izin verir
		if r.Method == "OPTIONS" {
			return
		} */
		paymentMethods, err := getPaymentMethods()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paymentMethods)
	})
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
}
