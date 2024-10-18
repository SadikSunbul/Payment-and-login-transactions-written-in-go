package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	clientID     = ""
	clientSecret = ""
)

func main() {
	http.HandleFunc("/pay", handlePay)
	http.HandleFunc("/success", handleSuccess)
	http.HandleFunc("/cancel", handleCancel)

	log.Println("Server starting at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handlePay(w http.ResponseWriter, r *http.Request) {
	token, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	approvalURL, err := createPayment(token)
	if err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Kullanıcıyı PayPal onay sayfasına yönlendirme
	http.Redirect(w, r, approvalURL, http.StatusTemporaryRedirect)
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	token, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	paymentID := r.URL.Query().Get("paymentId")
	payerID := r.URL.Query().Get("PayerID")

	saleID, err := executePayment(token, paymentID, payerID)
	if err != nil {
		http.Error(w, "Failed to execute payment", http.StatusInternalServerError)
		return
	}

	// Sipariş oluşturma işlemi burada yapılabilir
	// Örneğin, veritabanına kaydetme

	fmt.Fprintln(w, "Payment completed successfully! Sale ID:", saleID)
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	token, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// URL'den saleID'yi al
	saleID := r.URL.Query().Get("saleID")
	if saleID == "" {
		http.Error(w, "Missing saleID in request", http.StatusBadRequest)
		return
	}

	err = refundPayment(token, saleID)
	if err != nil {
		http.Error(w, "Failed to refund payment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Payment canceled and refunded successfully!")
}

func getAccessToken() (string, error) {
	url := "https://api.sandbox.paypal.com/v1/oauth2/token"
	reqBody := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("could not get access token")
}

func createPayment(token string) (string, error) {
	url := "https://api.sandbox.paypal.com/v1/payments/payment"

	paymentBody := map[string]interface{}{
		"intent": "sale",
		"payer": map[string]interface{}{
			"payment_method": "paypal",
		},
		"transactions": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"total":    "10.00", // Ödeme tutarı
					"currency": "USD",
				},
				"description": "This is the payment description.",
			},
		},
		"redirect_urls": map[string]interface{}{
			"return_url": "http://localhost:3000/success", // Ödeme başarılı olursa
			"cancel_url": "http://localhost:3000/cancel",  // Ödeme iptal olursa
		},
	}

	reqBody, err := json.Marshal(paymentBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Ödeme onayı için kullanıcıyı yönlendireceğimiz PayPal URL'sini alıyoruz
	if links, ok := result["links"].([]interface{}); ok {
		for _, link := range links {
			linkMap := link.(map[string]interface{})
			if rel, exists := linkMap["rel"]; exists && rel == "approval_url" {
				return linkMap["href"].(string), nil
			}
		}
	}

	return "", fmt.Errorf("approval_url not found")
}

func executePayment(token, paymentID, payerID string) (string, error) {
	url := fmt.Sprintf("https://api.sandbox.paypal.com/v1/payments/payment/%s/execute", paymentID)

	executeBody := map[string]interface{}{
		"payer_id": payerID,
	}

	reqBody, err := json.Marshal(executeBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("payment execution failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Sale ID'yi al
	if transactions, ok := result["transactions"].([]interface{}); ok {
		for _, transaction := range transactions {
			transactionMap := transaction.(map[string]interface{})
			if relatedResources, exists := transactionMap["related_resources"].([]interface{}); exists {
				for _, resource := range relatedResources {
					resourceMap := resource.(map[string]interface{})
					if sale, exists := resourceMap["sale"].(map[string]interface{}); exists {
						if saleID, exists := sale["id"].(string); exists {
							return saleID, nil
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("sale_id not found")
}
func refundPayment(token, saleID string) error {
	url := fmt.Sprintf("https://api.sandbox.paypal.com/v1/payments/sale/%s/refund", saleID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 200, 201 veya 204 durum kodlarını kabul et
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("refund failed with status: %d", resp.StatusCode)
	}

	fmt.Println("Payment refunded successfully!")
	return nil
}
