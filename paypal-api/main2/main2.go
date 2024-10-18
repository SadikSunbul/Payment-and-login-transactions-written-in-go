package main2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
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
	})

	/*
		http://localhost:3000/success?paymentId=<PAYMENT_ID>&PayerID=<PAYER_ID>
	*/

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessToken()
		if err != nil {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		paymentID := r.URL.Query().Get("paymentId")
		payerID := r.URL.Query().Get("PayerID")

		err = executePayment(token, paymentID, payerID)
		if err != nil {
			http.Error(w, "Failed to execute payment", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Payment completed successfully!")
	})

	http.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Payment canceled!")
	})

	log.Println("Server starting at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

const (
	clientID     = ""
	clientSecret = ""
)

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
					"total":    "100.00", // Ödeme tutarı
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

func executePayment(token, paymentID, payerID string) error {
	url := fmt.Sprintf("https://api.sandbox.paypal.com/v1/payments/payment/%s/execute", paymentID)

	executeBody := map[string]interface{}{
		"payer_id": payerID,
	}

	reqBody, err := json.Marshal(executeBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment execution failed with status: %d", resp.StatusCode)
	}

	fmt.Println("Payment executed successfully!")
	return nil
}

// Admin nasıl isterse ona göre değiştirilebilir. fee olayalrını ısterse kulalnıcıya odetırebırlı
